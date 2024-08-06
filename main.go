package main

import (
	"fmt"
	"math/rand"

	"github.com/iadams749/MancalaBot/internal/game"
)

func main() {
	game := game.New()

	for !game.Finished {
		moves := game.ValidMoves()

		idx := rand.Intn(len(moves))

		 _ = game.DoMove(moves[idx])

		game.Print()
		fmt.Println()
	}
}
