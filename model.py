import numpy as np


class Model:
    def __init__(self):
        self.board = np.array([[4, 4, 4, 4, 4, 4], [4, 4, 4, 4, 4, 4]])
        self.p1_score = 0
        self.p2_score = 0

    def doMove(self, x, y, p1_turn=True):
        num_marbles = self.board[y][x]
        self.board[y][x] = 0

        spot_x = x
        spot_y = y

        # Loop for placing the marbles across the board
        while num_marbles > 0:

            # Determining  the spot where the next marble will be placed
            if spot_y == 0:
                spot_x -= 1

                if spot_x == -1 and p1_turn:
                    spot_y = 1
                    spot_x = 0

            elif spot_y == 1:
                spot_x += 1

                if x == 6 and not p1_turn:
                    spot_y = 0
                    spot_x = 5

            # Placing a marble in the appropriate slot
            if -1 < spot_x < 6:
                self.board[spot_y][spot_x] += 1
            elif spot_x == -1:
                spot_y = 1
                self.p2_score += 1
            elif spot_x == 6:
                spot_y = 0
                self.p1_score += 1

            num_marbles -= 1

        # Handling the end-of-turn actions

        # Checking if the last marble was in the pocket
        if spot_x == -1:
            return False
        elif spot_x == 6:
            return True

        # Checking if the last marble captures
        if self.board[spot_y][spot_x] == 1 and self.board[int(not spot_y)][spot_x] != 0 and spot_y == int(p1_turn):
            if p1_turn:
                self.p1_score += 1
                self.board[spot_y][spot_x] = 0
                self.p1_score += self.board[int(not spot_y)][spot_x]
                self.board[int(not spot_y)][spot_x] = 0
            elif not p1_turn:
                self.p2_score += 1
                self.board[spot_y][spot_x] = 0
                self.p2_score += self.board[int(not spot_y)][spot_x]
                self.board[int(not spot_y)][spot_x] = 0

        return not p1_turn
