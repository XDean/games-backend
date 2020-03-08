package lostcities

import (
	"errors"
	"math/rand"
)

const (
	COLOR       = 6
	COLOR_COUNT = 11
)

const (
	StatusSet1 Status = iota
	StatusSet2
	StatusBuySwap
	StatusBanYun
	StatusDrawPlay
	StatusOver
)
const (
	GuanShui = iota
	BanYun
	BiYue
	Boat
)

type (
	Status int

	Card int // -1 means unknown
	Item int

	Board struct {
		status Status

		current     int
		playerCount int

		deck    []Card
		items   map[Item]int
		goods   map[Card]int
		board   [6]Card
		players []Player
	}

	Player struct {
		boats  []Card
		hand   map[Card]int
		items  map[Item]bool
		points int
	}
)

func NewStandardBoard(playerCount int) *Board {
	deck := make([]Card, COLOR*COLOR_COUNT)
	for i := range deck {
		deck[i] = Card(i % COLOR)
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	board := [6]Card{}
	copy(board[:], deck[:6])
	deck = deck[6:]

	goods := map[Card]int{}
	for i := 0; i < COLOR; i++ {
		goods[Card(i)] = 5
	}

	players := make([]Player, playerCount)
	for i := range players {
		hand := map[Card]int{}
		for _, c := range deck[:3] {
			hand[c]++
		}
		players[i] = Player{
			items:  map[Item]bool{},
			hand:   hand,
			boats:  []Card{-1, -1},
			points: 0,
		}
		deck = deck[3:]
	}

	g := &Board{
		status:      StatusSet1,
		current:     rand.Intn(playerCount),
		playerCount: playerCount,
		deck:        deck,
		goods:       goods,
		items:       map[Item]int{GuanShui: 2, BiYue: 2, BanYun: 2},
		board:       board,
		players:     players,
	}
	return g
}

func (g *Board) BuyItem(item Item, card Card) error {
	if g.status != StatusBuySwap {
		return errors.New("现在不是购买阶段")
	}
	if !item.IsValid() {
		return errors.New("参数不合法")
	}
	if g.players[g.current].points < item.Cost() {
		return errors.New("你没有足够的钱购买道具")
	}
	if item == Boat {
		if !card.IsValid() {
			return errors.New("上货类型不合法")
		}
		if g.goods[card] == 0 {
			return errors.New("该种货物已用完")
		}
		g.players[g.current].points -= item.Cost()
		g.players[g.current].boats = append(g.players[g.current].boats, card)
		g.status = StatusDrawPlay
	} else {
		if g.players[g.current].items[item] {
			return errors.New("不能重复购买相同的道具")
		}
		if g.items[item] == 0 {
			return errors.New("该道具已经售罄")
		}
		g.players[g.current].points -= item.Cost()
		g.players[g.current].items[item] = true
		g.items[item]--
		if item == BanYun {
			g.status = StatusBanYun
		} else {
			g.status = StatusDrawPlay
		}
	}
	return nil
}

func (g *Board) Swap(index1 int, card1 Card, index2 int, card2 Card) error {
	if g.status == StatusDrawPlay {
		return errors.New("现在不是换货阶段")
	}
	indexValid := func(index int) bool { return index >= 0 && index < len(g.players[g.current].boats) }
	if !indexValid(index1) || !card1.IsValid() || (indexValid(index2) != card2.IsValid()) || index1 == index2 {
		return errors.New("输入参数不合法")
	}
	isSwap2 := indexValid(index2)
	if isSwap2 {
		if g.players[g.current].items[BanYun] {
			if g.status == StatusBanYun {
				return errors.New("购买搬运工的回合只能换一船货")
			}
		} else {
			return errors.New("你没有搬运工，不能一次换两船货物")
		}
	}

	swap := func(index int, card Card) (func(), error) {
		if g.goods[card] == 0 {
			return nil, errors.New("该货物已经搬空")
		} else {
			return func() {
				old := g.players[g.current].boats[index]
				if old.IsValid() {
					g.goods[old]++
				}
				g.goods[card]--
				g.players[g.current].boats[index] = card
			}, nil
		}
	}
	f1, err1 := swap(index1, card1)
	if err1 != nil {
		return err1
	}
	if isSwap2 {
		f2, err2 := swap(index2, card2)
		if err2 != nil {
			return err2
		}
		f1()
		f2()
	} else {
		f1()
	}
	if g.status == StatusSet1 {
		if g.current == g.playerCount-1 {
			g.status = StatusSet2
		} else {
			g.current++
		}
	} else if g.status == StatusSet2 {
		if g.current == 0 {
			g.status = StatusBuySwap
		} else {
			g.current--
		}
	} else {
		g.status = StatusDrawPlay
	}
	return nil
}

func (g *Board) SkipSwap() error {
	if g.status == StatusBanYun || g.status == StatusBuySwap {
		return errors.New("该阶段无法跳过")
	}
	g.status = StatusDrawPlay
	return nil
}

func (g *Board) DrawCard(biyue bool) ([]Card, error) {
	if g.status != StatusDrawPlay {
		return nil, errors.New("现在不是抽牌阶段")
	}
	if !g.players[g.current].items[BiYue] && biyue {
		return nil, errors.New("你没有交易所牌，不能额外摸牌")
	}
	count := 1
	if biyue {
		count = 2
	}
	if len(g.deck) <= count {
		count = len(g.deck)
	}
	cards := g.deck[:count]
	for _, c := range cards {
		g.players[g.current].hand[c]++
	}
	g.deck = g.deck[count:]
	if len(g.deck) == 0 {
		g.status = StatusOver
	} else {
		g.status = StatusBuySwap
		g.current = (g.current + 1) % g.playerCount
	}
	return cards, nil
}

func (g *Board) PlayCard(card Card, dest [6]bool, biyue bool) (biyueCard Card, err error) {
	count := 0
	for _, d := range dest {
		if d {
			count++
		}
	}
	if g.status != StatusDrawPlay {
		return 0, errors.New("现在不是出牌阶段")
	}
	if count == 0 || !card.IsValid() {
		return 0, errors.New("参数不合法")
	}
	if !g.players[g.current].items[BiYue] && biyue {
		return 0, errors.New("你没有交易所牌，不能额外摸牌")
	}
	if g.players[g.current].hand[card] < count {
		return 0, errors.New("你没有足够的该种货物")
	}
	g.players[g.current].hand[card] -= count
	for i, d := range dest {
		if d {
			g.board[i] = card
		}
	}
	boardTotal := 0
	for _, c := range g.board {
		if c == card {
			boardTotal++
		}
	}
	for i := range g.players {
		playerTotal := 0
		for _, b := range g.players[i].boats {
			if b == card {
				playerTotal++
			}
		}
		if playerTotal != 0 {
			g.players[i].points += boardTotal * playerTotal
			if g.players[i].items[GuanShui] {
				g.players[i].points += 2
			}
		}
	}
	defer func() {
		if len(g.deck) == 0 {
			g.status = StatusOver
		} else {
			g.status = StatusBuySwap
			g.current = (g.current + 1) % g.playerCount
		}
	}()
	if biyue {
		card := g.deck[0]
		g.players[g.current].hand[card]++
		g.deck = g.deck[1:]
		return card, nil
	} else {
		return 0, nil
	}
}

func (c Card) IsValid() bool {
	return c >= 0 && c < COLOR
}

func (p Player) HandCount() int {
	sum := 0
	for _, c := range p.hand {
		sum += c
	}
	return sum
}

func (i Item) IsValid() bool {
	switch i {
	case Boat, GuanShui, BanYun, BiYue:
		return true
	default:
		return false
	}
}

func (i Item) Cost() int {
	switch i {
	case Boat:
		return 10
	case GuanShui:
		return 11
	case BanYun:
		return 12
	case BiYue:
		return 8
	default:
		return 0
	}
}
