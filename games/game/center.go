package game

import (
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

type gameMeta struct {
	name    string
	factory Factory
	hosts   map[int]*Host
}

var (
	gameMetas = map[string]gameMeta{}
)

func Register(name string, factory Factory) {
	if _, ok := gameMetas[name]; ok {
		logrus.Fatal("Duplicate game: " + name)
	} else {
		gameMetas[name] = gameMeta{
			name:    name,
			factory: factory,
			hosts:   map[int]*Host{},
		}
	}
}

func GetFactory(game string) Factory {
	if meta, ok := gameMetas[game]; ok {
		return meta.factory
	} else {
		return nil
	}
}

func GetHost(game string, id int) *Host {
	if meta, ok := gameMetas[game]; ok {
		return meta.hosts[id]
	} else {
		return nil
	}
}

func CreateHost(game string, g Game) *Host {
	if meta, ok := gameMetas[game]; ok {
		for {
			id := funk.RandomInt(1000, 9999)
			if _, ok := meta.hosts[id]; !ok {
				host := NewHost(id, g)
				host.Run()
				meta.hosts[id] = host
				return host
			}
		}
	} else {
		return nil
	}
}
