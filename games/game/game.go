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
	}
)
