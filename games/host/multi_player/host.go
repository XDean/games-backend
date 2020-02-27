package multi_player

import (
	"fmt"
	"games-backend/games/host"
)

type (
	Game interface {
		Handle(ctx Context) error
		PlayerCount() int
	}

	Context struct {
		host.Context
		host *Host
	}

	Host struct {
		game Game

		playing  bool
		ready    map[string]bool
		players  []string
		watchers []string
	}
)

func NewHost(game Game) *Host {
	return &Host{
		game:     game,
		watchers: []string{},
		ready:    map[string]bool{},
		players:  make([]string, game.PlayerCount()),
	}
}

func (h *Host) Handle(ctx host.Context) error {
	id := ctx.ClientId
	multiContext := Context{host: h, Context: ctx}
	switch ctx.Topic {
	case "host-info":
		ctx.SendEvent(id, host.TopicEvent{
			Topic:   "host-info",
			Payload: h.toInfo(),
		})
	case "join":
		if h.isPlayer(id) {
			return fmt.Errorf("你已经加入了该房间")
		}
		if seat, ok := h.availableSeat(); ok {
			h.players[seat] = id
			multiContext.SendAll(host.TopicEvent{
				Topic: "join",
				Payload: playerInfo{
					Id:    id,
					Seat:  seat,
					Ready: false,
				},
			})
		} else {
			return fmt.Errorf("房间已满")
		}
	case "watch":
		if h.isWatcher(id) {
			return fmt.Errorf("你已在观战该房间")
		}
		h.watchers = append(h.watchers, id)
		multiContext.SendAll(host.TopicEvent{
			Topic: "watch",
			Payload: watcherInfo{
				Id: id,
			},
		})
	case "ready":
		if h.isPlayer(id) {
			seat, _ := h.getSeat(id)
			ready := false
			err := ctx.GetPayload(&ready)
			if err != nil {
				return err
			}
			h.ready[id] = ready
			multiContext.SendAll(host.TopicEvent{
				Topic: "ready",
				Payload: playerInfo{
					Id:    id,
					Seat:  seat,
					Ready: ready,
				},
			})
			if ready && h.isAllReady() {
				h.playing = true
				return ctx.TriggerEvent(host.TopicEvent{Topic: "game-start"})
			}
		}
	}
	return h.game.Handle(multiContext)
}

func (h *Host) isPlayer(id string) bool {
	for _, theId := range h.players {
		if theId == id {
			return true
		}
	}
	return false
}

func (h *Host) isWatcher(id string) bool {
	for _, theId := range h.watchers {
		if theId == id {
			return true
		}
	}
	return false
}

func (h *Host) getSeat(id string) (int, bool) {
	for seat, theId := range h.players {
		if theId == id {
			return seat, true
		}
	}
	return 0, false
}

func (h *Host) availableSeat() (int, bool) {
	for i := 0; i < h.game.PlayerCount(); i++ {
		if h.players[i] == "" {
			return i, true
		}
	}
	return 0, false
}

func (h *Host) isAllReady() bool {
	for _, id := range h.players {
		if id == "" {
			return false
		} else if !h.ready[id] {
			return false
		}
	}
	return true
}

func (h *Host) allPlayers() []string {
	return append(h.watchers, h.players...)
}
