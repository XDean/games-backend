package game

import (
	"errors"
	"games-backend/games/host"
	"github.com/thoas/go-funk"
)

var (
	hosts = map[string]map[int]Server{}
)

func getServer(game string, id int) (Server, bool) {
	server, ok := hosts[game][id]
	return server, ok
}

func createServer(game string) (Server, error) {
	if meta, ok := host.GetMeta(game); !ok {
		return Server{}, errors.New("No such game")
	} else {
		newHost := meta.Factory.NewHost()
		err := newHost.Init()
		if err != nil {
			return Server{}, err
		}
		server := newServer(nextId(game), newHost)
		if _, ok := hosts[game]; !ok {
			hosts[game] = map[int]Server{}
		}
		hosts[game][server.id] = server
		return server, nil
	}
}

func nextId(game string) int {
	for {
		id := funk.RandomInt(1000, 9999)
		if ids, ok := hosts[game]; ok {
			if _, ok := ids[id]; ok {
				continue
			}
		}
		return id
	}
}
