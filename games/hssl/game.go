package lostcities

import (
	"errors"
	"games-backend/games/host"
	"games-backend/games/host/multi_player"
	"games-backend/util"
)

type (
	Game struct {
		board   *Board
		seatMap multi_player.RoomGameSeatMap
	}
)

func (g *Game) stub() {
	// goland always report error if remove this
}

func (g *Game) NewGame(ctx multi_player.Context) error {
	if err := g.checkPlaying(false); err != nil {
		return err
	}
	g.board = NewStandardBoard(ctx.GetPlayerCount())
	g.seatMap = ctx.GetRoomGameSeatMap()
	ctx.SendAllEach(func(id string) host.TopicEvent {
		return g.toInfoEvent(ctx, id)
	})
	return nil
}

func (g *Game) Handle(ctx multi_player.Context) (err error) {
	switch ctx.Topic {
	case topicInfo:
		if g.board != nil {
			ctx.SendBack(g.toInfoEvent(ctx, ctx.ClientId))
		}
	case topicSet:
		event := SettingRequest{}
		pos := 0
		if g.board.status == StatusSet2 {
			pos = 1
		}
		var seat int

		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return ctx.GetPayload(&event) },
			func() error { return g.board.Swap(pos, event.Card, -1, -1) },
			func() error {
				ctx.SendAll(host.TopicEvent{
					Topic: topicSet,
					Payload: SettingResponse{
						SeatEvent:      SeatEvent{Seat: g.seatMap.GameToRoom[seat]},
						SettingRequest: event,
						Round:          pos,
					},
				})
				g.sendStatus(ctx)
				return nil
			})
	case topicBuy:
		event := BuyRequest{}
		var seat int
		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return ctx.GetPayload(&event) },
			func() error { return g.board.BuyItem(event.Item, event.Card) },
			func() error {
				ctx.SendAll(host.TopicEvent{
					Topic: topicBuy,
					Payload: BuyResponse{
						SeatEvent:  SeatEvent{Seat: g.seatMap.GameToRoom[seat]},
						BuyRequest: event,
					},
				})
				g.sendStatus(ctx)
				return nil
			})
	case topicSwap:
		event := SwapRequest{}
		var seat int
		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return ctx.GetPayload(&event) },
			func() error { return g.board.Swap(event.Index1, event.Card1, event.Index2, event.Card2) },
			func() error {
				ctx.SendAll(host.TopicEvent{
					Topic: topicSwap,
					Payload: SwapResponse{
						SeatEvent:   SeatEvent{Seat: g.seatMap.GameToRoom[seat]},
						SwapRequest: event,
					},
				})
				g.sendStatus(ctx)
				return nil
			})
	case topicBanYun:
		event := BanYunRequest{}
		var seat int
		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return ctx.GetPayload(&event) },
			func() error { return g.board.Swap(event.Index, event.Card, -1, -1) },
			func() error {
				ctx.SendAll(host.TopicEvent{
					Topic: topicBanYun,
					Payload: BanYunResponse{
						SeatEvent:     SeatEvent{Seat: g.seatMap.GameToRoom[seat]},
						BanYunRequest: event,
					},
				})
				g.sendStatus(ctx)
				return nil
			})
	case topicSkip:
		var seat int
		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return g.board.SkipSwap() },
			func() error {
				ctx.SendAll(host.TopicEvent{
					Topic: topicSkip,
					Payload: SeatEvent{
						Seat: g.seatMap.GameToRoom[seat],
					},
				})
				g.sendStatus(ctx)
				return nil
			})
	case topicPlay:
		event := PlayRequest{}
		var seat int
		biyueCard := Card(-1)
		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return ctx.GetPayload(&event) },
			func() error {
				c, err := g.board.PlayCard(event.Card, event.Dest, event.BiYue)
				biyueCard = c
				return err
			},
			func() error {
				response := PlayResponse{
					SeatEvent:   SeatEvent{Seat: g.seatMap.GameToRoom[seat]},
					PlayRequest: event,
					BiYueCard:   biyueCard,
				}
				ctx.SendBack(host.TopicEvent{
					Topic:   topicPlay,
					Payload: response,
				})
				ctx.SendWatchers(host.TopicEvent{
					Topic:   topicPlay,
					Payload: response,
				})
				response.BiYueCard = -1
				ctx.SendExcludeSeat(host.TopicEvent{
					Topic:   topicPlay,
					Payload: response,
				}, g.seatMap.GameToRoom[seat])
				g.sendStatus(ctx)
				return nil
			},
			g.checkOver(ctx))
	case topicDraw:
		event := DrawRequest{}
		var seat int
		cards := make([]Card, 0)
		return util.DoUntilError(
			func() error { return g.checkPlaying(true) },
			func() error { seat = g.board.current; return nil },
			func() error { return g.checkCurrent(ctx) },
			func() error { return ctx.GetPayload(&event) },
			func() error {
				c, err := g.board.DrawCard(event.BiYue)
				cards = c
				return err
			},
			func() error {
				response := DrawResponse{
					SeatEvent:   SeatEvent{Seat: g.seatMap.GameToRoom[seat]},
					DrawRequest: event,
					Cards:       cards,
				}
				ctx.SendBack(host.TopicEvent{
					Topic:   topicDraw,
					Payload: response,
				})
				ctx.SendWatchers(host.TopicEvent{
					Topic:   topicDraw,
					Payload: response,
				})
				response.Cards = make([]Card, len(response.Cards))
				for i := range response.Cards {
					response.Cards[i] = -1
				}
				ctx.SendExcludeSeat(host.TopicEvent{
					Topic:   topicDraw,
					Payload: response,
				}, g.seatMap.GameToRoom[seat])
				g.sendStatus(ctx)
				return nil
			},
			g.checkOver(ctx))
	}
	return nil
}

func (g *Game) checkOver(ctx multi_player.Context) func() error {
	return func() error {
		if g.board != nil {
			if g.board.status == StatusOver {
				return ctx.TriggerEvent(host.TopicEvent{Topic: multi_player.TopicOverInner})
			}
		}
		return nil
	}
}

func (g *Game) MinPlayerCount() int {
	return 2
}

func (g *Game) MaxPlayerCount() int {
	return 4
}

func (g *Game) checkPlaying(expect bool) error {
	if g.board == nil || g.board.status == StatusOver {
		if expect {
			return errors.New("游戏尚未开始")
		}
	} else {
		if !expect {
			return errors.New("游戏已经开始")
		}
	}
	return nil
}

func (g *Game) checkCurrent(ctx multi_player.Context) error {
	player := ctx.GetPlayerById(ctx.ClientId)
	if player == nil {
		return errors.New("你不是该局玩家")
	}
	if g.board != nil && g.board.current != g.seatMap.RoomToGame[player.GetSeat()] {
		return errors.New("现在不是你的回合")
	}
	return nil
}

func (g *Game) sendStatus(ctx multi_player.Context) {
	ctx.SendAll(host.TopicEvent{
		Topic: topicStatus,
		Payload: StatusResponse{
			Status:  g.board.status,
			Current: g.seatMap.GameToRoom[g.board.current],
		},
	})
}

func (g *Game) toInfoEvent(ctx multi_player.Context, id string) host.TopicEvent {
	player := ctx.GetPlayerById(id)
	playerInfos := make([]PlayerInfo, 0)
	for gameSeat, p := range g.board.players {
		playerInfo := PlayerInfo{
			Seat:   g.seatMap.GameToRoom[gameSeat],
			Hand:   p.hand,
			Boats:  p.boats,
			Items:  p.items,
			Points: p.points,
		}
		if player != nil && g.seatMap.RoomToGame[player.GetSeat()] != gameSeat { // not this player
			playerInfo.Hand = map[Card]int{-1: p.HandCount()}
			playerInfo.Points = -1
		}
		playerInfos = append(playerInfos, playerInfo)
	}
	return host.TopicEvent{
		Topic: topicInfo,
		Payload: GameInfo{
			Status:      g.board.status,
			Current:     g.seatMap.GameToRoom[g.board.current],
			PlayerCount: g.board.playerCount,
			Deck:        len(g.board.deck),
			Items:       g.board.items,
			Goods:       g.board.goods,
			Board:       g.board.board,
			Players:     playerInfos,
		},
	}
}
