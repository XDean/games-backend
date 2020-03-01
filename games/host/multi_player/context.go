package multi_player

import (
	"games-backend/games/host"
	"github.com/thoas/go-funk"
)

type (
	RoomGameSeat struct {
		RoomToGame []int
		GameToRoom []int
	}
)

func (c Context) GetRoomGameSeat() RoomGameSeat {
	roomToGame := make([]int, c.game.MaxPlayerCount())
	gameToRoom := make([]int, c.game.MaxPlayerCount())
	gameSeat := 0
	for seat, player := range c.players {
		if player != nil {
			roomToGame[seat] = gameSeat
			gameToRoom[gameSeat] = seat
			gameSeat++
		}
	}
	return RoomGameSeat{
		RoomToGame: roomToGame,
		GameToRoom: gameToRoom,
	}
}

func (c Context) Send(id string, event host.TopicEvent) {
	c.SendEvent(id, event)
}

func (c Context) SendAll(event host.TopicEvent) {
	c.SendPlayers(event)
	c.SendWatchers(event)
}

func (c Context) SendAllEach(event func(id string) host.TopicEvent) {
	for _, player := range c.players {
		if player != nil {
			c.SendEvent(player.id, event(player.id))
		}
	}
	for _, watcher := range c.watchers {
		c.SendEvent(watcher.id, event(watcher.id))
	}
}

func (c Context) SendPlayers(event host.TopicEvent) {
	for _, player := range c.players {
		if player != nil {
			c.Send(player.id, event)
		}
	}
}

func (c Context) SendWatchers(event host.TopicEvent) {
	for _, watcher := range c.watchers {
		c.Send(watcher.id, event)
	}
}

func (c Context) SendSeat(event host.TopicEvent, seats ...int) {
	for _, seat := range seats {
		player := c.players[seat]
		if player != nil {
			c.Send(player.id, event)
		}
	}
}

func (c *Context) SendExcludeSeat(event host.TopicEvent, seats ...int) {
	for seat, player := range c.players {
		if player != nil && !funk.ContainsInt(seats, seat) {
			c.Send(player.id, event)
		}
	}
}
