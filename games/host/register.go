package host

import (
	"github.com/sirupsen/logrus"
)

type Meta struct {
	Name    string
	Factory Factory
}

var (
	hostMetas = map[string]Meta{}
)

func Register(meta Meta) {
	if _, ok := hostMetas[meta.Name]; ok {
		logrus.Fatal("Duplicate game: " + meta.Name)
	} else {
		hostMetas[meta.Name] = meta
	}
}

func GetMeta(name string) (Meta, bool) {
	meta, ok := hostMetas[name]
	return meta, ok
}
