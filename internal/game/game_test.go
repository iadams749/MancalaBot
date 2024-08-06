package game

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want *Game
	}{
		{
			name: "happy path",
			want: &Game{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DoMove() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGame_DoMove(t *testing.T) {
	t.Parallel()
	type fields struct {
		Board    []int
		Turn     uint8
		Finished bool
	}
	type args struct {
		spot int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Game
		wantErr error
	}{
		{
			name: "new game move 1",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 0,
			},
			want: &Game{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "new game move 2",
			fields: fields{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 7,
			},
			want: &Game{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 0, 5, 5, 5, 5, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p1 turn again",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 2,
			},
			want: &Game{
				Board:    []int{4, 4, 0, 5, 5, 5, 1, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p2 turn again",
			fields: fields{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 9,
			},
			want: &Game{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 4, 4, 0, 5, 5, 5, 1},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p1 around the bend",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 8, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 5,
			},
			want: &Game{
				Board:    []int{5, 4, 4, 4, 4, 0, 1, 5, 5, 5, 5, 5, 5, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p2 around the bend",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 8, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 12,
			},
			want: &Game{
				Board:    []int{5, 5, 5, 5, 5, 5, 0, 5, 4, 4, 4, 4, 0, 1},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p1 pot capture",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 0, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 0,
			},
			want: &Game{
				Board:    []int{0, 5, 5, 5, 0, 4, 5, 4, 0, 4, 4, 4, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p2 pot capture",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 0, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 7,
			},
			want: &Game{
				Board:    []int{4, 0, 4, 4, 4, 4, 0, 0, 5, 5, 5, 0, 4, 5},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p1 pseudo capture",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 3, 0, 0, 0, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 5,
			},
			want: &Game{
				Board:    []int{4, 4, 4, 4, 4, 0, 1, 1, 1, 4, 4, 4, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "p2 pseudo capture",
			fields: fields{
				Board:    []int{0, 0, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 3, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 12,
			},
			want: &Game{
				Board:    []int{1, 1, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 0, 1},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: nil,
		},
		{
			name: "finished no pot emptying",
			fields: fields{
				Board:    []int{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 5,
			},
			want: &Game{
				Board:    []int{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: true,
			},
			wantErr: nil,
		},
		{
			name: "finished p2 pot emptying",
			fields: fields{
				Board:    []int{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 5, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 5,
			},
			want: &Game{
				Board:    []int{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 5},
				Turn:     playerOneTurn,
				Finished: true,
			},
			wantErr: nil,
		},
		{
			name: "finished p1 pot emptying",
			fields: fields{
				Board:    []int{0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 12,
			},
			want: &Game{
				Board:    []int{0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 1},
				Turn:     playerTwoTurn,
				Finished: true,
			},
			wantErr: nil,
		},
		{
			name: "invalid move empty spot",
			fields: fields{
				Board:    []int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 1,
			},
			want: &Game{
				Board:    []int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: ErrInvalidMove,
		},
		{
			name: "invalid move p1 wrong spot",
			fields: fields{
				Board:    []int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			args: args{
				spot: 8,
			},
			want: &Game{
				Board:    []int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			wantErr: ErrInvalidMove,
		},
		{
			name: "invalid move p2 wrong spot",
			fields: fields{
				Board:    []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			args: args{
				spot: 0,
			},
			want: &Game{
				Board:    []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			wantErr: ErrInvalidMove,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := &Game{
				Board:    tt.fields.Board,
				Turn:     tt.fields.Turn,
				Finished: tt.fields.Finished,
			}

			err := g.DoMove(tt.args.spot)

			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DoMove() error mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want, g); diff != "" {
				t.Errorf("DoMove() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestGame_ValidMoves(t *testing.T) {
	t.Parallel()
	type fields struct {
		Board    []int
		Turn     uint8
		Finished bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "new game player 1 turn",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			want: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name: "in-progress game player 2 turn",
			fields: fields{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 4, 4, 0, 5, 5, 5, 1},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			want: []int{7, 8, 10, 11, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := &Game{
				Board:    tt.fields.Board,
				Turn:     tt.fields.Turn,
				Finished: tt.fields.Finished,
			}
			got := g.ValidMoves()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DoMove() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGame_GetScores(t *testing.T) {
	t.Parallel()
	type fields struct {
		Board    []int
		Turn     uint8
		Finished bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
		want1  int
		want2  bool
	}{
		{
			name: "new game",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			want:  0,
			want1: 0,
			want2: false,
		},
		{
			name: "finished game",
			fields: fields{
				Board:    []int{0, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0, 24},
				Turn:     playerOneTurn,
				Finished: true,
			},
			want:  24,
			want1: 24,
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := &Game{
				Board:    tt.fields.Board,
				Turn:     tt.fields.Turn,
				Finished: tt.fields.Finished,
			}
			got, got1, got2 := g.GetScores()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DoMove() P1 score mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want1, got1); diff != "" {
				t.Errorf("DoMove() P2 score mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want2, got2); diff != "" {
				t.Errorf("DoMove() Finished mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGame_Print(t *testing.T) {
	type fields struct {
		Board    []int
		Turn     uint8
		Finished bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "new game p1 turn",
			fields: fields{
				Board:    []int{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerOneTurn,
				Finished: false,
			},
			want: "                   -->               P1  \n-----------------------------------------\n|    | 04 | 04 | 04 | 04 | 04 | 04 |    |\n| 00 |-----------------------------| 00 |\n|    | 04 | 04 | 04 | 04 | 04 | 04 |    |\n-----------------------------------------\n  P2               <--                   \nCurrent turn: P1\n",
		},
		{
			name: "turn 2 p2 turn",
			fields: fields{
				Board:    []int{0, 5, 5, 5, 5, 4, 0, 4, 4, 4, 4, 4, 4, 0},
				Turn:     playerTwoTurn,
				Finished: false,
			},
			want: "                   -->               P1  \n-----------------------------------------\n|    | 00 | 05 | 05 | 05 | 05 | 04 |    |\n| 00 |-----------------------------| 00 |\n|    | 04 | 04 | 04 | 04 | 04 | 04 |    |\n-----------------------------------------\n  P2               <--                   \nCurrent turn: P2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board:    tt.fields.Board,
				Turn:     tt.fields.Turn,
				Finished: tt.fields.Finished,
			}

			// Create a pipe to capture stdout
			r, w, _ := os.Pipe()
			// Save the original stdout
			origStdout := os.Stdout
			// Redirect stdout to the pipe writer
			os.Stdout = w

			// Ensure stdout is restored and pipe writer is closed after the test
			defer func() {
				os.Stdout = origStdout
				w.Close()
			}()

			g.Print()

			// Close the pipe writer to allow reading from the pipe
			w.Close()

			// Read the captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			r.Close()

			// Check the buffer content
			got := buf.String()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DoMove() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
