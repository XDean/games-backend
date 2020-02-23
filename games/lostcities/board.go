package lostcities

import (
	"github.com/thoas/go-funk"
	"math/rand"
)

const (
	DEVELOP_POINT = 20
	BONUS_COUNT   = 8
	BONUS_POINT   = 20

	COLOR       = 5
	CARD        = 12
	CARD_DOUBLE = 3
	HAND        = 8
)

type (
	Board struct {
		over    bool
		current int
		deck    []Card      // [index] from 0 (top)
		drop    [][]Card    // [color][index] from 0 (oldest)
		board   [2][][]Card // [player][color][index] from 0 (oldest)
		hand    [2][]Card   // [player][index] no order, by default 0 (oldest)
		score   [2]int
	}
)

func NewStandardBoard() *Board {
	deck := make([]Card, CARD*COLOR)
	for i := range deck {
		deck[i] = Card(i)
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	board := [2][][]Card{}
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

	hand := [2][]Card{}
	for i := range hand {
		hand[i] = make([]Card, 0)
	}
	g := &Board{
		current: 0,
		deck:    deck,
		board:   board,
		drop:    drop,
		hand:    hand,
		score:   [2]int{},
	}
	g.DrawCard(0, HAND)
	g.DrawCard(1, HAND)
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
	if len(g.deck) == 0 {
		g.over = true
		for player, board := range g.board {
			score := 0
			for _, cards := range board {
				sum := 0
				times := 1
				count := 0
				if len(board) != 0 {
					sum = -DEVELOP_POINT
					for _, card := range cards {
						if card.IsDouble() {
							times++
						} else {
							score += card.Point()
							count++
						}
					}
				}
				sum = times * sum
				if count >= BONUS_COUNT {
					sum += BONUS_POINT
				}
				score += sum
			}
			g.score[player] = score
		}
	} else {
		g.current = (g.current + 1) % 2
	}
}
