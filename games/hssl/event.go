package lostcities

type (
	// Common
	SeatEvent struct {
		Seat int `json:"seat"`
	}

	// Info
	GameInfo struct {
		Status      Status       `json:"status"`
		Current     int          `json:"current"`
		PlayerCount int          `json:"count"`
		Deck        int          `json:"deck"`
		Items       map[Item]int `json:"items"`
		Goods       map[Card]int `json:"goods"`
		Board       [6]Card      `json:"board"`
		Players     []PlayerInfo `json:"players"`
	}

	PlayerInfo struct {
		Boats  []Card        `json:"boats"`
		Hand   map[Card]int  `json:"hand"`
		Items  map[Item]bool `json:"items"`
		Points int           `json:"points"`
	}

	// Status
	StatusResponse struct {
		Status  Status `json:"status"`
		Current int    `json:"current"`
	}

	// Setting
	SettingRequest struct {
		Card Card `json:"card"`
	}

	SettingResponse struct {
		SeatEvent
		SettingRequest
	}

	// BuySwap + Banyun
	BuyRequest struct {
		Item Item `json:"item"`
		Card Card `json:"card"` // only for buy boat
	}

	BuyResponse struct {
		SeatEvent
		BuyRequest
	}

	SwapRequest struct {
		Index1 int  `json:"index1"`
		Card1  Card `json:"card1"`
		Index2 int  `json:"index2"`
		Card2  Card `json:"card2"`
	}

	SwapResponse struct {
		SeatEvent
		SwapRequest
	}

	BanYunRequest struct {
		Index int  `json:"index"`
		Card  Card `json:"card"`
	}

	BanYunResponse struct {
		SeatEvent
		BanYunRequest
	}

	// PlayDraw
	PlayRequest struct {
		Card  Card    `json:"card"`
		Dest  [6]bool `json:"dest"`
		BiYue bool    `json:"biyue"`
	}

	PlayResponse struct {
		SeatEvent
		PlayRequest
	}

	DrawRequest struct {
		BiYue bool `json:"biyue"`
	}

	DrawResponse struct {
		SeatEvent
		DrawRequest
	}

	DrawPrivateResponse struct {
		SeatEvent
		Cards []Card `json:"cards"`
	}
)
