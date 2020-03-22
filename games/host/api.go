package host

import (
	"fmt"
	"games-backend/util/inject"
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
		Inject  inject.Context
		Handler EventHandler
	}

	Context struct {
		ClientId     string
		Topic        string
		GetPayload   func(payload interface{}) error
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

func NewHost(handler EventHandler) *Host {
	h := Host{
		Inject:  inject.NewContext(),
		Handler: handler,
	}
	h.Inject.Register(handler)
	return &h
}

func NewContext(ctx Context) Context {
	if ctx.GetPayload == nil {
		ctx.GetPayload = func(_ interface{}) error {
			return fmt.Errorf("No Payload")
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

func (c Context) SendBack(event TopicEvent) {
	c.SendEvent(c.ClientId, event)
}

func (f EventHandlerFunc) Handle(ctx Context) error {
	return f(ctx)
}

func (f PluginFunc) Plug(handler EventHandler) EventHandler {
	return f(handler)
}
