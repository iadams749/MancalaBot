package game

import "fmt"

// constants for player turns
const (
	PlayerOneTurn = iota
	PlayerTwoTurn
)

// constants for the spots on the board
const (
	p1pot1 = iota
	p1pot2
	p1pot3
	p1pot4
	p1pot5
	p1pot6
	p1home
	p2pot1
	p2pot2
	p2pot3
	p2pot4
	p2pot5
	p2pot6
	p2home
)

var (
	ErrInvalidMove   = fmt.Errorf("invalid move")
	ErrInvalidMarker = fmt.Errorf("invalid marker")
)

// Game contains the model for a particular game state.
// It includes the state of the board and who's turn it is.
type Game struct {
	Board    []int
	Turn     uint8
	Finished bool
}

// New returns a pointer to a game in it's initial state.
func New() *Game {
	return &Game{
		Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
		Turn:     PlayerOneTurn,
		Finished: false,
	}
}

// DoMove applies a particular move to a game state.
func (g *Game) DoMove(spot int) error {
	// checking if the move is valid
	if !g.validateMove(spot) {
		return ErrInvalidMove
	}

	// picking up the marbles in the spot to start the move at
	holding := g.Board[spot]
	marker := spot + 1
	g.Board[spot] = 0

	// depositing the marbles in the spots on the board
	for holding > 0 {
		// checking if the current marker is in the opposite home
		if marker == p1home {
			if g.Turn == PlayerTwoTurn {
				marker++
				continue
			}
		} else if marker == p2home {
			if g.Turn == PlayerOneTurn {
				marker = p1pot1
				continue
			}
		}

		// incrementing the pot
		g.Board[marker] = g.Board[marker] + 1
		holding--

		// moving the marker if there are marbles left
		if holding != 0 {
			if marker == p2home {
				marker = p1pot1
			} else {
				marker++
			}
		}
	}

	// checking if the marker ended in an empty pot
	g.checkEmptyPot(marker)

	// checking if the game is over
	g.Finished = g.checkGameOver()

	if g.Finished {
		return nil
	}

	// checking whose turn it will be after the move
	if marker != p1home && marker != p2home {
		// flipping who's turn it is
		if g.Turn == PlayerOneTurn {
			g.Turn = PlayerTwoTurn
		} else if g.Turn == PlayerTwoTurn {
			g.Turn = PlayerOneTurn
		}

		return nil
	}

	return nil
}

// validateMove checks if a position is a valid move for a game state.
func (g *Game) validateMove(spot int) bool {
	if g.Board[spot] == 0 {
		return false
	}
	
	if g.Turn == PlayerOneTurn {
		if spot < p1pot1 || spot > p1pot6 {
			return false
		}
	} else if g.Turn == PlayerTwoTurn {
		if spot < p2pot1 || spot > p2pot6 {
			return false
		}
	}

	return true
}

// checkEmptyPot checks if the marker landed on an empty pot, and move the appropriate marbles.
func (g *Game) checkEmptyPot(marker int) {
	if marker == p1home || marker == p2home {
		return
	}

	if g.Board[marker] != 1 || g.Board[12-marker] == 0 {
		return
	}

	if g.Turn == PlayerOneTurn {
		if marker >= p1pot1 && marker <= p1pot6 {
			g.Board[6] = g.Board[6] + g.Board[marker] + g.Board[12-marker]
			g.Board[marker] = 0
			g.Board[12-marker] = 0
		} else {
			return
		}
	} else if g.Turn == PlayerTwoTurn {
		if marker >= p2pot1 && marker <= p2pot6 {
			g.Board[13] = g.Board[13] + g.Board[marker] + g.Board[12-marker]
			g.Board[marker] = 0
			g.Board[12-marker] = 0
		} else {
			return
		}
	}
}

// checkGameOver checks if the game is over, and updates the board accordingly
func (g *Game) checkGameOver() bool {
	// checking if player one's pots are empty
	playerOneEmpty := true
	for i := p1pot1; i <= p1pot6; i++ {
		if g.Board[i] != 0 {
			playerOneEmpty = false
			break
		}
	}

	// emptying the board and returning if player one's board is empty
	if playerOneEmpty {
		for i := p2pot1; i <= p2pot6; i++ {
			g.Board[13] += g.Board[i]
			g.Board[i] = 0
		}

		return true
	}

	// checking if player two's pots are empty
	playerTwoEmpty := true
	for i := p2pot1; i <= p2pot6; i++ {
		if g.Board[i] != 0 {
			playerTwoEmpty = false
			break
		}
	}

	//emptying the board and returning if player two's board is empty
	if playerTwoEmpty {
		for i := p1pot1; i <= p1pot6; i++ {
			g.Board[6] += g.Board[i]
			g.Board[i] = 0
		}

		return true
	}

	return false
}

// ValidMoves returns an array containing the positions of the valid moves for the game.
func (g *Game) ValidMoves() []int {
	var validMoves []int

	if g.Turn == PlayerOneTurn {
		for i := p1pot1; i <= p1pot6; i++ {
			if g.Board[i] != 0 {
				validMoves = append(validMoves, i)
			}
		}
	} else if g.Turn == PlayerTwoTurn {
		for i := p2pot1; i <= p2pot6; i++ {
			if g.Board[i] != 0 {
				validMoves = append(validMoves, i)
			}
		}
	}

	return validMoves
}

// GetScores returns the scores of the two players, and a bool indicating whether or not the game is over
func (g *Game) GetScores() (int, int, bool) {
	return g.Board[p1home], g.Board[p2home], g.Finished
}

// Print prints out the current game state to the console.
func (g *Game) Print() {

	turn := "P1"
	if g.Turn != PlayerOneTurn {
		turn = "P2"
	}

	fmt.Println("                   -->               P1  ")
	fmt.Println("-----------------------------------------")
	fmt.Printf("|    | %02d | %02d | %02d | %02d | %02d | %02d |    |\n", g.Board[0], g.Board[1], g.Board[2], g.Board[3], g.Board[4], g.Board[5])
	fmt.Printf("| %02d |-----------------------------| %02d |\n", g.Board[13], g.Board[6])
	fmt.Printf("|    | %02d | %02d | %02d | %02d | %02d | %02d |    |\n", g.Board[12], g.Board[11], g.Board[10], g.Board[9], g.Board[8], g.Board[7])
	fmt.Println("-----------------------------------------")
	fmt.Println("  P2               <--                   ")
	fmt.Printf("Current turn: %s\n", turn)
}
