package lostcities

import (
	"fmt"
	"games-backend/games/host"
	"games-backend/games/host/multi_player"
)

type (
	Game struct {
		*Board
		history [][]GameEvent
	}

	InfoEvent struct {
	}

	GameEvent struct {
		Seat int `json:"seat"`
		// Play card
		Card Card `json:"card"`
		Drop bool `json:"drop"`
		// Draw card
		FromDeck     bool `json:"deck"`           // Or from drop
		Color        int  `json:"draw-color"`     // Only available when from drop
		DrawDropCard Card `json:"draw-drop-card"` // Only for response
	}

	GameInfo struct {
		Over        bool     `json:"over"`
		Seat        int      `json:"seat"`
		CurrentSeat int      `json:"current-seat"`
		Deck        int      `json:"deck"`
		MyBoard     [][]Card `json:"my-board"`
		OtherBoard  [][]Card `json:"other-board"`
		DropBoard   [][]Card `json:"drop-board"`
		Hand        []Card   `json:"hand"`
		Score       [2]int   `json:"score"`
	}
)

func (g *Game) Handle(ctx multi_player.Context) error {
	switch ctx.Topic {
	case "game-start":
		if g.Board == nil || g.over {
			g.Board = NewStandardBoard()
			ctx.SendEach(func(id string) host.TopicEvent {
				return g.gameInfo(ctx, "start", id)
			})
		}
	case "game-info":
		if g.Board != nil {
			ctx.SendBack(g.gameInfo(ctx, "game-info", ctx.ClientId))
		}
	case "play":
		event := GameEvent{}
		err := ctx.GetPayload(&event)
		if err != nil {
			return err
		}
		return g.Play(ctx, ctx.ClientId, event)
	}
	return nil
}

func (g *Game) PlayerCount() int {
	return 2
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
	defer func() {
		g.next()
		if g.over {
			ctx.SendAll(host.TopicEvent{
				Topic:   "over",
				Payload: g.score,
			})
		} else {
			ctx.SendAll(host.TopicEvent{
				Topic:   "turn",
				Payload: g.current,
			})
		}
	}()
	//g.history[g.current] = append(g.history[g.current], event)
	event.Seat = g.current
	// TODO Check deck has card
	if event.Drop {
		g.DropCard(g.current, event.Card)
	} else {
		g.PlayCard(g.current, event.Card)
	}
	if event.FromDeck {
		draw := g.DrawCard(g.current, 1)[0]
		defer ctx.SendSeat(host.TopicEvent{
			Topic:   "draw",
			Payload: draw,
		}, g.current)
	} else {
		draw := g.DrawDropCard(g.current, event.Color)
		event.DrawDropCard = draw
	}
	ctx.SendAll(event.asTopic())
	return nil
}

func (g *Game) gameInfo(ctx multi_player.Context, topic string, id string) host.TopicEvent {
	if seat, ok := ctx.GetSeat(id); ok {
		return host.TopicEvent{
			Topic: topic,
			Payload: GameInfo{
				Over:        g.over,
				Seat:        seat,
				CurrentSeat: g.current,
				Deck:        len(g.deck),
				Hand:        g.hand[seat],
				DropBoard:   g.drop,
				MyBoard:     g.board[seat],
				OtherBoard:  g.board[1-seat],
				Score:       g.score,
			},
		}
	} else {
		return host.TopicEvent{
			Topic: topic,
			Payload: GameInfo{
				Over:        g.over,
				CurrentSeat: g.current,
				Deck:        len(g.deck),
				DropBoard:   g.drop,
				MyBoard:     g.board[0],
				OtherBoard:  g.board[1],
				Score:       g.score,
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
