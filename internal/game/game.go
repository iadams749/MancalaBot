package game

import "fmt"

// constants for player turns
const (
	playerOneTurn = iota
	playerTwoTurn
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
	Board []int
	Turn  uint8
}

// New returns a pointer to a game in it's initial state.
func New() *Game {
	return &Game{
		Board: []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
		Turn:  playerOneTurn,
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
			if g.Turn == playerTwoTurn {
				marker++
				continue
			}
		} else if marker == p2home {
			if g.Turn == playerOneTurn {
				marker = p1pot1
				continue
			}
		}

		// incrementing the pot
		g.Board[marker] = g.Board[marker] + 1
		holding--

		// moving the marker
		if marker == p2home {
			marker = p1pot1
		} else if holding != 0 {
			marker++
		}
	}

	// checking who's turn it will be after the move
	if marker == p1home {
		if g.Turn == playerOneTurn {
			return nil
		} else {
			return ErrInvalidMarker
		}
	} else if marker == p2home {
		if g.Turn == playerTwoTurn {
			return nil
		} else {
			return ErrInvalidMarker
		}
	}

	// checking if the marker ended in an empty pot
	g.checkEmptyPot(marker)
		
	// flipping who's turn it is
	if g.Turn == playerOneTurn {
		g.Turn = playerTwoTurn
	} else if g.Turn == playerTwoTurn {
		g.Turn = playerOneTurn
	}

	return nil
}

// ValidateMove checks if a position is a valid move for a game state.
func (g *Game) validateMove(spot int) bool {
	if g.Turn == playerOneTurn {
		if spot >= p1pot1 && spot <= p1pot6 {
			return true
		}
		return false
	} else if g.Turn == playerTwoTurn {
		if spot >= p2pot1 && spot <= p2pot6 {
			return true
		}
		return false
	}

	if g.Board[spot] == 0 {
		return false
	}

	return false
}

// checkEmptyPot checks if the marker landed on an empty pot, and move the appropriate marbles
func (g *Game) checkEmptyPot(marker int) {
	if g.Board[marker] != 1 {
		return
	}

	if g.Turn == playerOneTurn {
		if marker >= p1pot1 && marker <= p1pot6 {
			g.Board[6] = g.Board[6] + g.Board[marker] + g.Board[12-marker]
			g.Board[marker] = 0
			g.Board[12-marker] = 0
		} else {
			return
		}
	} else if g.Turn == playerTwoTurn {
		if marker >= p2pot1 && marker <= p2pot6 {
			g.Board[13] = g.Board[13] + g.Board[marker] + g.Board[12-marker]
			g.Board[marker] = 0
			g.Board[12-marker] = 0
		}
	}
}

// Print prints out the current game state to the console.
func (g *Game) Print() {

	turn := "P1"
	if g.Turn != playerOneTurn {
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
