package lostcities

import (
	"fmt"
	"games-backend/games/host"
	"games-backend/games/host/multi_player"
)

type (
	Game struct {
		board *Board
	}
)

func (g *Game) NewGame(ctx multi_player.Context) error {
	if g.board == nil || g.status == Over {
		g.board = NewStandardBoard(ctx)
		ctx.SendEach(func(id string) host.TopicEvent {
			return g.gameInfo(ctx, "game-start", id)
		})
	}
	return nil
}

func (g *Game) Handle(ctx multi_player.Context) error {
	switch ctx.Topic {
	case "game-start":
		if g.Board == nil || g.over {
			g.Board = NewStandardBoard()
			ctx.SendEach(func(id string) host.TopicEvent {
				return g.gameInfo(ctx, "game-start", id)
			})
		}
	case "game-info":
		if g.Board != nil {
			ctx.SendBack(g.gameInfo(ctx, "game-info", ctx.ClientId))
		}
	case "play":
		if g.Board == nil {
			return fmt.Errorf("游戏尚未开始")
		}
		event := GameEvent{}
		err := ctx.GetPayload(&event)
		if err != nil {
			return err
		}
		return g.Play(ctx, ctx.ClientId, event)
	}
	return nil
}

func (g *Game) MinPlayerCount() int {
	return 3
}

func (g *Game) MaxPlayerCount() int {
	return 5
}

func (g *Game) Play(ctx multi_player.Context, id string, event GameEvent) error {
	if g.over {
		return fmt.Errorf("游戏已经结束")
	}
	if seat, ok := ctx.GetSeat(id); !ok {
		return fmt.Errorf("你不是该局玩家")
	} else if g.current != seat {
		return fmt.Errorf("现在不是你的回合")
	}
	if event.Drop && !event.FromDeck && (event.Card.Color() == event.Color) {
		return fmt.Errorf("你不能立刻摸起刚刚弃置的牌")
	}
	if !g.hasCard(g.current, event.Card) {
		return fmt.Errorf("卡牌不存在")
	}
	cards := g.board[g.current][event.Card.Color()]
	if !event.Drop && len(cards) > 0 && !cards[0].IsDouble() && cards[0].Point() > event.Card.Point() {
		return fmt.Errorf("每个系列的卡牌必须递增打出")
	}
	if !event.FromDeck && len(g.drop[event.Color]) == 0 {
		return fmt.Errorf("该弃牌堆中没有牌")
	}
	//g.history[g.current] = append(g.history[g.current], event)
	event.Seat = g.current
	if event.Drop {
		g.DropCard(g.current, event.Card)
	} else {
		g.PlayCard(g.current, event.Card)
	}
	if event.FromDeck {
		event.DeckCard = g.DrawCard(g.current, 1)[0]
	} else {
		event.DrawDropCard = g.DrawDropCard(g.current, event.Color)
	}
	ctx.SendSeat(event.asTopic(), g.current)
	ctx.SendWatchers(event.asTopic())
	event.DeckCard = -1
	ctx.SendExcludeSeat(event.asTopic(), g.current)

	g.next()
	if g.over {
		return ctx.TriggerEvent(host.TopicEvent{Topic: "game-over"})
	} else {
		ctx.SendAll(host.TopicEvent{
			Topic:   "turn",
			Payload: g.current,
		})
	}
	return nil
}

func (g *Game) gameInfo(ctx multi_player.Context, topic string, id string) host.TopicEvent {
	if seat, ok := ctx.GetSeat(id); ok {
		hand := g.hand
		hand[1-seat] = []Card{}
		return host.TopicEvent{
			Topic: topic,
			Payload: GameInfo{
				Over:        g.over,
				CurrentSeat: g.current,
				Deck:        len(g.deck),
				Hand:        hand,
				Drop:        g.drop,
				Board:       g.board,
			},
		}
	} else {
		return host.TopicEvent{
			Topic: topic,
			Payload: GameInfo{
				Over:        g.over,
				CurrentSeat: g.current,
				Deck:        len(g.deck),
				Hand:        g.hand,
				Drop:        g.drop,
				Board:       g.board,
			},
		}
	}
}

func (e GameEvent) asTopic() host.TopicEvent {
	return host.TopicEvent{
		Topic:   "play",
		Payload: e,
	}
}
