package lostcities

import (
	"games-backend/games/host"
	"games-backend/games/host/multi_player"
	"games-backend/games/host/plugin"
)

func init() {
	host.Register(host.Meta{
		Name:    "lostcities",
		Factory: Factory{},
	})
}

type Factory struct {
}

func (f Factory) NewHost() host.Host {
	return host.NewHost(multi_player.NewRoom(&Game{})).
		Plug(plugin.NewConnect()).
		Plug(plugin.NewChat())
}
