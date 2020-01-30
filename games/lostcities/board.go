package lostcities

import (
	"github.com/thoas/go-funk"
	"math/rand"
)

const (
	DEVELOP_POINT = 20
	BONUS_COUNT   = 8
	BONUS_POINT   = 20

	PLAYER      = 2
	COLOR       = 5
	CARD        = 12
	CARD_DOUBLE = 3
	CARD_POINT  = CARD - CARD_DOUBLE
)

type (
	Board struct {
		over    bool
		current int
		deck    []Card     // [index] from 0 (top)
		board   [][][]Card // [player][color][index] from 0 (oldest)
		drop    [][]Card   // [color][index] from 0 (oldest)
		hand    [][]Card   // [player][index] no order, by default 0 (oldest)
	}
)

func NewStandardBoard() *Board {
	deck := make([]Card, CARD*COLOR)
	for i := range deck {
		deck[i] = Card(i)
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	board := make([][][]Card, PLAYER)
	for i := range board {
		board[i] = make([][]Card, COLOR)
		for m := range board[i] {
			board[i][m] = make([]Card, 0)
		}
	}

	drop := make([][]Card, COLOR)
	for i := range drop {
		drop[i] = make([]Card, 0)
	}

	hand := make([][]Card, PLAYER)
	for i := range hand {
		hand[i] = make([]Card, 0)
	}
	g := &Board{
		current: 0,
		deck:    deck,
		board:   board,
		drop:    drop,
		hand:    hand,
	}
	g.DrawCard(0, 7)
	g.DrawCard(1, 7)
	return g
}

func (g *Board) DrawCard(player int, count int) []Card {
	card := g.deck[:count]
	g.deck = g.deck[count:]
	g.hand[player] = append(g.hand[player], card...)
	return card
}

func (g *Board) PlayCard(player int, card Card) {
	if g.removeHandCard(player, card) {
		g.board[player][card.Color()] = append(g.board[player][card.Color()], card)
	}
}

func (g *Board) DropCard(player int, card Card) {
	if g.removeHandCard(player, card) {
		g.drop[card.Color()] = append(g.drop[card.Color()], card)
	}
}

func (g *Board) DrawDropCard(player int, color int) Card {
	drop := g.drop[color]
	card := drop[len(drop)-1]
	g.drop[color] = drop[:len(drop)-1]
	g.hand[player] = append(g.hand[player], card)
	return card
}

func (g *Board) hasCard(player int, card Card) bool {
	return funk.Contains(g.hand[player], card)
}

func (g *Board) removeHandCard(player int, card Card) bool {
	index := funk.IndexOf(g.hand[player], card)
	if index != -1 {
		g.hand[player] = append(g.hand[player][:index], g.hand[player][index+1:]...)
		return true
	}
	return false
}

func (g *Board) next() {
	g.current = (g.current + 1) % PLAYER
}
