package multi_player

const (
	TopicInfo      = "room-info"
	TopicJoin      = "room-join"
	TopicWatch     = "room-watch"
	TopicReady     = "room-ready"
	TopicSwap      = "room-swap-seat"
	TopicStart     = "room-game-start"
	TopicOver      = "room-game-over"
	TopicOverInner = "_room-game-over"
)

type (
	SwapSeatRequest struct {
		TargetSeat int `json:"target"`
	}

	SwapSeatResponse struct {
		FromSeat   int `json:"from"`
		TargetSeat int `json:"target"`
	}
)
