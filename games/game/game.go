package game

type (
	Factory interface {
		NewConfig() interface{}
		NewGame(config interface{}) Game
	}
	Game interface {
		NewEvent(topic string) interface{}
		Init(host *Host)
		HandleEvent(client *Client, event interface{})
		PlayerCount() int
	}

	Context struct {
		Id         string
		Topic      string
		GetPayload func(payload interface{})
		SendEvent  func(id string, event TopicEvent)
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
