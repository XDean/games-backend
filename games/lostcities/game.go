package lostcities

import (
	"games-backend/games/game"
	"strings"
)

type (
	Game struct {
		*Board

		host    *game.Host
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
		Score       []int    `json:"score"`
	}
)

func (g *Game) PlayerCount() int {
	return 2
}

func (g *Game) Init(host *game.Host) {
	g.host = host
}

func (g *Game) NewEvent(topic string) interface{} {
	switch strings.ToLower(topic) {
	case "play":
		return &GameEvent{}
	case "game-info":
		return &InfoEvent{}
	default:
		return nil
	}
}

func (g *Game) HandleEvent(client *game.Client, event interface{}) {
	switch e := event.(type) {
	case *game.StartEvent:
		if g.Board == nil || g.over {
			g.Board = NewStandardBoard()
			g.host.SendEach(func(c *game.Client) game.TopicEvent {
				return g.gameInfo("start", c)
			})
		}
	case *game.ConnectEvent, *InfoEvent:
		if g.Board != nil {
			client.Send(g.gameInfo("game-info", client))
		}
	case *GameEvent:
		g.Play(client, *e)
	}
}

func (g *Game) Play(client *game.Client, event GameEvent) {
	if g.over {
		client.Error("游戏已经结束")
		return
	}
	if seat, ok := client.Seat(); !ok {
		client.Error("你不是该局玩家")
		return
	} else if g.current != seat {
		client.Error("现在不是你的回合")
		return
	}
	if event.Drop && !event.FromDeck && (event.Card.Color() == event.Color) {
		client.Error("你不能立刻摸起刚刚弃置的牌")
		return
	}
	if !g.hasCard(g.current, event.Card) {
		client.Error("卡牌不存在")
		return
	}
	cards := g.board[g.current][event.Card.Color()]
	if !event.Drop && len(cards) > 0 && !cards[0].IsDouble() && cards[0].Point() > event.Card.Point() {
		client.Error("每个系列的卡牌必须递增打出")
		return
	}
	if !event.FromDeck && len(g.drop[event.Color]) == 0 {
		client.Error("该弃牌堆中没有牌")
		return
	}
	defer func() {
		g.next()
		if g.over {
			g.sendAll(game.TopicEvent{
				Topic:   "over",
				Payload: g.score,
			})
		} else {
			g.sendAll(game.TopicEvent{
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
		defer g.sendToSeat(game.TopicEvent{
			Topic:   "draw",
			Payload: draw,
		}, g.current)
	} else {
		draw := g.DrawDropCard(g.current, event.Color)
		event.DrawDropCard = draw
	}
	g.sendAll(event.topic())
}

func (g *Game) gameInfo(topic string, client *game.Client) game.TopicEvent {
	if seat, ok := client.Seat(); ok {
		return game.TopicEvent{
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
		return game.TopicEvent{
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

func (g *Game) sendAll(event game.TopicEvent) {
	g.host.SendAll(event)
}

func (g *Game) sendToSeat(event game.TopicEvent, seats ...int) {
	g.host.SendToSeat(event, seats...)
}

func (g *Game) sendExcludeSeat(event game.TopicEvent, seats ...int) {
	g.host.SendExcludeSeat(event, seats...)
}

func (g *Game) sendError(err string) {
	g.host.SendToSeat(game.ErrorEvent(err), g.current)
}

func (e GameEvent) topic() game.TopicEvent {
	return game.TopicEvent{
		Topic:   "play",
		Payload: e,
	}
}
