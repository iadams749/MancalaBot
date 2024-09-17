package player

import (
	"math/rand"

	"github.com/iadams749/MancalaBot/internal/game"
)

// A RandomPlayer takes in a game state and randomly selects a legal move.
type RandomPlayer struct {}

// GetMove picks a random legal move and returns it.
func (*RandomPlayer) GetMove(game *game.Game) int {
	moves := game.ValidMoves()

	idx := rand.Intn(len(moves))

	return moves[idx]
}