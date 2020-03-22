package plugin

import (
	"games-backend/games/host"
)

type (
	Connect struct {
		connected map[string]bool
	}
)

func NewConnect() *Connect {
	return &Connect{connected: map[string]bool{}}
}

func (c *Connect) Plug(handler host.EventHandler) host.EventHandler {
	return host.EventHandlerFunc(func(ctx host.Context) error {
		switch ctx.Topic {
		case host.TopicConnect:
			c.connected[ctx.ClientId] = true
		case host.TopicDisConnect:
			delete(c.connected, ctx.ClientId)
		}
		return handler.Handle(ctx)
	})
}

func (c *Connect) GetAll() []string {
	res := make([]string, 0)
	for k := range c.connected {
		res = append(res, k)
	}
	return res
}
