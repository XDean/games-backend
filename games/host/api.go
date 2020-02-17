package host

import (
	"fmt"
	"games-backend/util"
)

type (
	TopicEvent struct {
		Topic   string      `json:"topic"`
		Payload interface{} `json:"payload"`
	}

	Factory interface {
		NewHost() Host
	}

	Host struct {
		Handler EventHandler
	}

	Context struct {
		ClientId     string
		Topic        string
		GetPayload   func(payload interface{}) error
		GetData      func(name string, receiver interface{}) error
		TriggerEvent func(event TopicEvent) error            // trigger event let the host handle
		SendEvent    func(clientId string, event TopicEvent) // send event to client
	}

	EventHandler interface {
		Handle(ctx Context) error
	}

	EventHandlerFunc func(ctx Context) error

	Plugin interface {
		Plug(handler EventHandler) EventHandler
	}

	PluginFunc func(handler EventHandler) EventHandler
)

func NewContext(ctx Context) Context {
	if ctx.GetPayload == nil {
		ctx.GetPayload = func(_ interface{}) error {
			return fmt.Errorf("No Payload")
		}
	}
	if ctx.GetData == nil {
		ctx.GetData = func(name string, _ interface{}) error {
			return fmt.Errorf("No Data: %s", name)
		}
	}
	if ctx.SendEvent == nil {
		ctx.SendEvent = func(_ string, _ TopicEvent) {}
	}
	if ctx.TriggerEvent == nil {
		ctx.TriggerEvent = func(_ TopicEvent) error { return nil }
	}
	return ctx
}

func (c Context) AddData(name string, data interface{}) Context {
	old := c.GetData
	reflectFunc := util.ReflectSetReceiver(data)
	c.GetData = func(n string, receiver interface{}) error {
		if name == n {
			return reflectFunc(receiver)
		}
		return old(n, receiver)
	}
	return c
}

func (c Context) SendBack(event TopicEvent) {
	c.SendEvent(c.ClientId, event)
}

func (f EventHandlerFunc) Handle(ctx Context) error {
	return f(ctx)
}

func (f PluginFunc) Plug(handler EventHandler) EventHandler {
	return f(handler)
}
