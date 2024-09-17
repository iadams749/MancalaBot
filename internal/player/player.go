package player

import "github.com/iadams749/MancalaBot/internal/game"

// An ai takes in a game state and returns a move.
type Player interface {
	GetMove(game *game.Game) int
}