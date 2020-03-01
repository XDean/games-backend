package lostcities

import (
	"games-backend/games/host"
	"games-backend/games/host/multi_player"
	"games-backend/games/host/plugin"
)

func init() {
	host.Register(host.Meta{
		Name:    "hssl",
		Factory: Factory{},
	})
}

type Factory struct {
}

func (f Factory) NewHost() host.Host {
	return host.Host{
		Handler: multi_player.NewHost(&Game{}),
	}.Plug(plugin.NewChat())
}
