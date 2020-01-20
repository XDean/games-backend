package lostcities

import "games-backend/games/game"

func init() {
	game.Register("lostcities", Factory{})
}

type Factory struct {
}

func (f Factory) NewConfig() interface{} {
	return struct{}{}
}

func (f Factory) NewGame(config interface{}) game.Game {
	return NewStandardGame()
}
