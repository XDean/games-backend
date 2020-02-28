package lostcities

type Card int

func (c Card) Color() int {
	return int(c) / CARD
}

func (c Card) Point() int {
	if c.IsDouble() {
		return 0
	} else {
		return int(c)%CARD - CARD_DOUBLE + 2
	}
}

func (c Card) IsDouble() bool {
	return int(c)%CARD < CARD_DOUBLE
}
