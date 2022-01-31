import numpy as np
from model import Model
from random import randint


class Game:
    def __init__(self, p1_turn=True):
        self.model = Model()
        self.turn = p1_turn
        self.game_over = False

    # Getting all possible moves for a certain game
    def get_moves(self):
        row = 0

        if self.turn:
            row = 1

        moves = []

        for x in range(0, 6):
            if self.model.board[row][x] != 0:
                moves.append((x, row))

        return moves

    def do_turn(self, x, y):
        self.turn = self.model.doMove(x, y, self.turn)

        # Checking if the game is over
        self.game_over = True

        if self.turn:
            for x in range(0, 6):
                if self.model.board[1][x] != 0:
                    self.game_over = False
        elif not self.turn:
            for x in range(0, 6):
                if self.model.board[0][x] != 0:
                    self.game_over = False

        # Clearing the marbles if the game is over
        if self.game_over:
            self.model.clearBoard()

    def printout(self):
        print(self.model.board)
        print("Player One: " + str(self.model.p1_score))
        print("Player Two: " + str(self.model.p2_score))
        if self.turn:
            print("Current Turn: Player 1")
        else:
            print("Current Turn: Player 2")

    def random_game(self):
        while not self.game_over:
            moves = self.get_moves()
            x, y = moves[randint(0, len(moves) - 1)]
            self.do_turn(x, y)
            self.printout()
