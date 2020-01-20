package lostcities

import (
	"games-backend/games/game"
	"strings"
)

type (
	GameEvent struct {
		// Play card
		Card Card
		Drop bool
		// Draw card
		FromDeck bool // Or from drop
		Color    int  // Only available when from deck
	}
)

func (g *Game) NewEvent(topic string) interface{} {
	switch strings.ToLower(topic) {
	case "play":
		return GameEvent{}
	default:
		return nil
	}
}

func (g *Game) HandleEvent(client *game.Client, event interface{}) {
	switch e := event.(type) {
	case GameEvent:
		g.Play(e)
	}
}
