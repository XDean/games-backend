package lostcities

import (
	"games-backend/games/game"
)

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

func (g *Game) PlayCard(player int, card Card) {
	if g.removeHandCard(player, card) {
		g.board[player][card.Color()] = append(g.board[player][card.Color()], card)
	}
}

func (g *Game) DropCard(player int, card Card) {
	if g.removeHandCard(player, card) {
		g.drop[card.Color()] = append(g.drop[card.Color()], card)
	}
}

func (g *Game) DrawCard(player int, count int) []Card {
	card := g.deck[:count]
	g.deck = g.deck[count:]
	g.hand[player] = append(g.hand[player], card...)
	return card
}

func (g *Game) DrawDropCard(player int, color int) Card {
	drop := g.drop[color]
	card := drop[len(drop)-1]
	g.drop[color] = drop[:len(drop)-1]
	g.hand[player] = append(g.hand[player], card)
	return card
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
