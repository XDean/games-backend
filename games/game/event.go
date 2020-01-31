package game

type (
	TopicEvent struct {
		Topic   string      `json:"topic"`
		Payload interface{} `json:"payload"`
	}

	ConnectEvent struct{}

	DisConnectEvent struct{}

	ReadyEvent struct {
		Ready bool
	}
)

func ErrorEvent(err string) TopicEvent {
	return TopicEvent{
		Topic:   "error",
		Payload: err,
	}
}
