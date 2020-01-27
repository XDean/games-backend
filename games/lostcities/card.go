package lostcities

type Card struct {
	int
}

func (c Card) Color() int {
	return c.int / CARD
}

func (c Card) Point() int {
	if c.IsDouble() {
		return 0
	} else {
		return c.int - CARD_DOUBLE + 2
	}
}

func (c Card) IsDouble() bool {
	return c.int < CARD_DOUBLE
}
