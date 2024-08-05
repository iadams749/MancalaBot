package main

import (
	"fmt"

	"github.com/iadams749/MancalaBot/internal/game"
)

func main() {
	game := game.New()

	err := game.DoMove(2)
	if err != nil {
		fmt.Println(err.Error())
	}

	game.Print()
}