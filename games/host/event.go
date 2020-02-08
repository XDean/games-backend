package host

func ErrorEvent(err string) TopicEvent {
	return TopicEvent{
		Topic:   "error",
		Payload: err,
	}
}
