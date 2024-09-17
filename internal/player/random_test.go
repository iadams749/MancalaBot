package player

import (
	"testing"

	"github.com/iadams749/MancalaBot/internal/game"
)

func TestRandomPlayer_GetMove(t *testing.T) {
	type args struct {
		game *game.Game
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "one possible move",
			args: args{
				game: &game.Game{
					Board:    []int{0, 4, 0, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0},
					Turn:     game.PlayerOneTurn,
					Finished: false,
				},
			},
			want: []int{1},
		},
		{
			name: "new game",
			args: args{
				game: &game.Game{
					Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
					Turn:     game.PlayerOneTurn,
					Finished: false,
				},
			},
			want: []int{0,1,2,3,4,5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RandomPlayer{}
			if got := r.GetMove(tt.args.game); !contains(tt.want, got) {
				t.Errorf("RandomPlayer.GetMove() = %v, want %v", got, tt.want)
			}
		})
	}
}

// contains checks if an int is in an array
func contains(slice []int, element int) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}