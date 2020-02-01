package game

type (
	TopicEvent struct {
		Topic   string      `json:"topic"`
		Payload interface{} `json:"payload"`
	}

	ConnectEvent struct{}

	DisConnectEvent struct{}

	ReadyEvent bool

	ChatEvent string

	StartEvent struct{}
)

func ErrorEvent(err string) TopicEvent {
	return TopicEvent{
		Topic:   "error",
		Payload: err,
	}
}
