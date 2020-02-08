package host

type (
	TopicEvent struct {
		Topic   string      `json:"topic"`
		Payload interface{} `json:"payload"`
	}

	Factory interface {
		NewHost() Host
	}

	Host struct {
		handler EventHandler
	}

	Context struct {
		Who          string
		Topic        string
		GetPayload   func(payload interface{}) error
		SendEvent    func(who string, event TopicEvent) // send event to client
		TriggerEvent func(event TopicEvent)             // trigger event let the host handle
	}

	EventHandler interface {
		Handle(ctx Context)
	}

	EventHandlerFunc func(ctx Context)

	Plugin interface {
		Plug(handler EventHandler) EventHandler
	}

	PluginFunc func(handler EventHandler) EventHandler
)

func (f EventHandlerFunc) Handle(ctx Context) {
	f(ctx)
}

func (f PluginFunc) Plug(handler EventHandler) EventHandler {
	return f(handler)
}
