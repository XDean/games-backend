package multi_player

import (
	"games-backend/games/host"
	"github.com/thoas/go-funk"
)

const ContextKey = "multi_player.sender"

func (c Context) IsPlayer(id string) bool {
	return c.host.isPlayer(id)
}

func (c Context) IsWatcher(id string) bool {
	return c.host.isWatcher(id)
}

func (c Context) GetSeat(id string) (int, bool) {
	return c.host.getSeat(id)
}

func (c Context) Send(id string, event host.TopicEvent) {
	c.SendEvent(id, event)
}

func (c Context) SendAll(event host.TopicEvent) {
	for _, id := range c.host.allPlayers() {
		c.SendEvent(id, event)
	}
}

func (c Context) SendEach(event func(id string) host.TopicEvent) {
	for _, id := range c.host.allPlayers() {
		c.SendEvent(id, event(id))
	}
}

func (c Context) SendPlayers(event host.TopicEvent) {
	for _, id := range c.host.players {
		if id != "" {
			c.Send(id, event)
		}
	}
}

func (c Context) SendWatchers(event host.TopicEvent) {
	for _, id := range c.host.watchers {
		c.Send(id, event)
	}
}

func (c Context) SendSeat(event host.TopicEvent, seats ...int) {
	for _, seat := range seats {
		id := c.host.players[seat]
		if id != "" {
			c.Send(id, event)
		}
	}
}

func (c *Context) SendExcludeSeat(event host.TopicEvent, seats ...int) {
	for seat, id := range c.host.players {
		if id != "" && !funk.ContainsInt(seats, seat) {
			c.Send(id, event)
		}
	}
}
