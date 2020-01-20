package game

type (
	Factory interface {
		NewConfig() interface{}
		NewGame(config interface{}) Game
	}
	Game interface {
		NewEvent(topic string) interface{}
		HandleEvent(client *Client, event interface{})
	}
)
