package lostcities

const (
	topicInfo   = "hssl-info"
	topicSet    = "hssl-set"
	topicBuy    = "hssl-buy"
	topicSwap   = "hssl-swap"
	topicBanYun = "hssl-banyun"
	topicSkip   = "hssl-skip-swap"
	topicPlay   = "hssl-play"
	topicDraw   = "hssl-draw"
	topicStatus = "hssl-status"
)

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
		Seat   int           `json:"seat"`
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
		Round int `json:"round"` // 0 or 1
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
		Index int  `json:"index1"`
		Card  Card `json:"card1"`
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
		BiYueCard Card `json:"biyue-card"`
	}

	DrawRequest struct {
		BiYue bool `json:"biyue"`
	}

	DrawResponse struct {
		SeatEvent
		DrawRequest
		Cards []Card `json:"cards"`
	}
)
