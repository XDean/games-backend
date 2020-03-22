package multi_player

import (
	"errors"
	"fmt"
	"games-backend/games/host"
	"games-backend/games/host/plugin"
)

type (
	Game interface {
		NewGame(ctx Context) error
		Handle(ctx Context) error
		MinPlayerCount() int
		MaxPlayerCount() int
	}

	Context struct {
		host.Context
		*Room
	}

	Room struct {
		game Game

		playing  bool
		players  []*Player // by seat
		watchers []*Watcher

		Connected *plugin.Connect `inject:""`
	}

	Player struct {
		id    string
		ready bool
		seat  int
		host  bool
	}

	Watcher struct {
		id string
	}
)

func NewRoom(game Game) *Room {
	return &Room{
		game:     game,
		watchers: []*Watcher{},
		players:  make([]*Player, game.MaxPlayerCount()),
	}
}

func (r *Room) Handle(ctx host.Context) error {
	id := ctx.ClientId
	multiContext := Context{Room: r, Context: ctx}
	switch ctx.Topic {
	case TopicInfo:
		ctx.SendEvent(id, host.TopicEvent{
			Topic:   TopicInfo,
			Payload: r.toInfo(),
		})
	case TopicJoin:
		if r.playing {
			return errors.New("游戏已经开始")
		}
		if r.IsPlayer(id) {
			return fmt.Errorf("你已经加入了该房间")
		}
		if seat, ok := r.availableSeat(); ok {
			r.players[seat] = &Player{
				id:    id,
				ready: false,
				seat:  seat,
				host:  r.GetHost() == nil,
			}
			multiContext.SendAllConnected(host.TopicEvent{
				Topic: TopicJoin,
				Payload: playerInfo{
					Id:    id,
					Seat:  seat,
					Ready: false,
					Host:  r.players[seat].host,
				},
			})
		} else {
			return fmt.Errorf("房间已满")
		}
	case TopicWatch:
		if r.IsWatcher(id) {
			return fmt.Errorf("你已在观战该房间")
		}
		r.watchers = append(r.watchers, &Watcher{id: id})
		multiContext.SendAllConnected(host.TopicEvent{
			Topic: TopicWatch,
			Payload: watcherInfo{
				Id: id,
			},
		})
	case TopicReady:
		if r.playing {
			return errors.New("游戏已经开始")
		}
		player := r.GetPlayerById(id)
		if player != nil {
			ready := false
			err := ctx.GetPayload(&ready)
			if err != nil {
				return err
			}
			player.ready = ready
			multiContext.SendAllConnected(host.TopicEvent{
				Topic: TopicReady,
				Payload: playerInfo{
					Id:    id,
					Seat:  player.seat,
					Ready: ready,
				},
			})
		}
	case TopicSwap:
		if r.playing {
			return errors.New("游戏已经开始")
		}
		player := r.GetPlayerById(id)
		if player != nil {
			if player.ready {
				return errors.New("已经准备不能换座位")
			}
			event := SwapSeatRequest{}
			err := ctx.GetPayload(&event)
			if err != nil {
				return err
			}
			if !r.ValidSeat(event.TargetSeat) {
				return errors.New("参数不合法")
			}
			fromSeat := player.seat
			targetPlayer := r.GetPlayerBySeat(event.TargetSeat)
			if targetPlayer != nil {
				if targetPlayer.ready {
					return errors.New("目标已经准备不能换座位")
				} else {
					player.seat, targetPlayer.seat = targetPlayer.seat, player.seat
					r.players[player.seat] = player
					r.players[targetPlayer.seat] = targetPlayer
				}
			} else {
				r.players[player.seat] = nil
				player.seat = event.TargetSeat
				r.players[event.TargetSeat] = player
			}
			multiContext.SendAllConnected(host.TopicEvent{
				Topic: TopicSwap,
				Payload: SwapSeatResponse{
					FromSeat:   fromSeat,
					TargetSeat: event.TargetSeat,
				},
			})
		}
	case TopicStart:
		hostPlayer := r.GetHost()
		if hostPlayer != nil && hostPlayer.id == id {
			if r.GetPlayerCount() < r.game.MinPlayerCount() {
				return errors.New("人数不足，无法开始游戏")
			}
			r.playing = true
			multiContext.SendAllConnected(host.TopicEvent{Topic: TopicStart})
			return r.game.NewGame(multiContext)
		} else {
			return errors.New("只有主机可以开始游戏")
		}
	case TopicOverInner:
		if r.playing {
			r.playing = false
			for _, player := range r.players {
				if player != nil {
					player.ready = false
				}
			}
			multiContext.SendAllConnected(host.TopicEvent{Topic: TopicOver})
		}
	}
	return r.game.Handle(multiContext)
}

func (r *Room) ValidSeat(seat int) bool {
	return seat >= 0 && seat < r.game.MaxPlayerCount()
}

func (r *Room) GetHost() *Player {
	for _, p := range r.players {
		if p != nil && p.host {
			return p
		}
	}
	return nil
}

func (r *Room) GetPlayerById(id string) *Player {
	for _, p := range r.players {
		if p != nil && p.id == id {
			return p
		}
	}
	return nil
}

func (r *Room) GetWatcherById(id string) *Watcher {
	for _, w := range r.watchers {
		if w.id == id {
			return w
		}
	}
	return nil
}

func (r *Room) GetPlayerBySeat(seat int) *Player {
	if seat >= 0 && seat < len(r.players) {
		return r.players[seat]
	} else {
		return nil
	}
}

func (r *Room) IsPlayer(id string) bool {
	return r.GetPlayerById(id) != nil
}

func (r *Room) IsWatcher(id string) bool {
	return r.GetWatcherById(id) != nil
}

func (r *Room) availableSeat() (int, bool) {
	for seat, player := range r.players {
		if player == nil {
			return seat, true
		}
	}
	return 0, false
}

func (r *Room) GetPlayers() []*Player {
	res := make([]*Player, 0)
	for _, p := range r.players {
		if p != nil {
			res = append(res, p)
		}
	}
	return res
}

func (r *Room) GetPlayerCount() int {
	return len(r.GetPlayers())
}

func (r *Room) isFull() bool {
	for _, player := range r.players {
		if player == nil {
			return false
		}
	}
	return true
}

func (r *Room) isAllReady() bool {
	for _, player := range r.players {
		if player != nil && !player.ready {
			return false
		}
	}
	return true
}

func (p *Player) GetSeat() int {
	return p.seat
}

func (p *Player) IsHost() bool {
	return p.host
}
