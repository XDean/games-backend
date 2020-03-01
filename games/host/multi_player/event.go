package multi_player

type (
	SwapSeatRequest struct {
		TargetSeat int `json:"target"`
	}

	SwapSeatResponse struct {
		FromSeat   int `json:"from"`
		TargetSeat int `json:"target"`
	}
)
