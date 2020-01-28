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
		// Play card
		Card Card
		Drop bool
		// Draw card
		FromDeck bool // Or from drop
		Color    int  // Only available when from deck
	}

	PrivateInfo struct {
		Seat        int
		CurrentSeat int
		Deck        int
		MyBoard     [][]Card
		OtherBoard  [][]Card
		Drop        [][]Card
		Hand        []Card
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
		return GameEvent{}
	case "game-info":
		return InfoEvent{}
	default:
		return nil
	}
}

func (g *Game) HandleEvent(client *game.Client, event interface{}) {
	switch e := event.(type) {
	case game.ConnectEvent, InfoEvent:
		if seat, ok := client.Seat(); ok {
			client.Send(game.TopicEvent{
				Topic: "game-info",
				Payload: PrivateInfo{
					Seat:        seat,
					CurrentSeat: g.current,
					Deck:        len(g.deck),
					Hand:        g.hand[seat],
					Drop:        g.drop,
					MyBoard:     g.board[seat],
					OtherBoard:  g.board[1-seat],
				},
			})
		} else {
			client.Send(game.TopicEvent{
				Topic: "game-info",
				Payload: PrivateInfo{
					CurrentSeat: g.current,
					Deck:        len(g.deck),
					Drop:        g.drop,
					MyBoard:     g.board[0],
					OtherBoard:  g.board[1],
				},
			})
		}
	case GameEvent:
		g.Play(e)
	}
}

func (g *Game) Play(event GameEvent) {
	if event.Drop && !event.FromDeck && (event.Card.Color() == event.Color) {
		g.sendError("You can't draw the drop card immediately")
		return
	}
	if !g.hasCard(g.current, event.Card) {
		g.sendError("You don't have the card to play")
		return
	}
	cards := g.board[g.current][event.Card.Color()]
	if !event.Drop && len(cards) > 0 && cards[0].Point() > event.Card.Point() {
		g.sendError("You can't play the card because you already have a larger one")
		return
	}
	if !event.FromDeck && len(g.drop[event.Color]) == 0 {
		g.sendError("No card to draw in this color's drop area")
		return
	}
	defer g.next()
	g.history[g.current] = append(g.history[g.current], event)
	g.sendAll(event.topic())
	// TODO Check deck has card
	if event.Drop {
		g.DropCard(g.current, event.Card)
	} else {
		g.PlayCard(g.current, event.Card)
	}
	if event.FromDeck {
		draw := g.DrawCard(g.current, 1)[0]
		g.sendToSeat(game.TopicEvent{
			Topic:   "draw",
			Payload: draw,
		}, g.current)
	} else {
		draw := g.DrawDropCard(g.current, event.Color)
		g.sendToSeat(game.TopicEvent{
			Topic:   "draw",
			Payload: draw,
		}, g.current)
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
