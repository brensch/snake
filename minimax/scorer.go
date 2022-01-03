package minimax

import (
	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/pather"
	"github.com/brensch/snake/rules"
)

// type

func PercentageOfBoardControlled(board *rules.BoardState) float64 {

	// you = 0

	// calculate how many squares i can reach first
	allMovesFromSquares := make([][]int32, len(board.Snakes))
	totalSpaces := int(board.Height * board.Width)

	// for each snake, find a naive distance to each point
	for snakeNumber, snake := range board.Snakes {

		movesFromSquares := make([]int32, totalSpaces)
		boardSpace := 0
		for x := int32(0); x < board.Width; x++ {
			for y := int32(0); y < board.Height; y++ {
				movesFromSquares[boardSpace] = generator.Abs(snake.Body[0].X-x) + generator.Abs(snake.Body[0].Y-y)
				boardSpace++
			}
		}

		allMovesFromSquares[snakeNumber] = movesFromSquares
	}

	closestSquareCount := 0

	for i := 0; i < totalSpaces; i++ {
		closestDistance := int32(1000)
		closestSnake := -1
		for snake := range board.Snakes {
			if allMovesFromSquares[snake][i] < int32(closestDistance) {
				closestDistance = allMovesFromSquares[snake][i]
				closestSnake = snake
			}
		}

		// check if index of closest snake is maxer index
		if closestSnake == 0 {
			closestSquareCount++
		}

	}

	return (float64(closestSquareCount) / float64(totalSpaces))

}

// TODO: make this actually check the shortest path to the point
func PercentageOfBoardControlledSmart(board *rules.BoardState) float64 {

	pather.CalculateAllAvailableSquares(board, board.Snakes[0].Body[0], board.Snakes[0].ID)

	return 0
	// return 0.5 - (float64(closestSquareCount) / float64(totalSpaces))

}

// GameFinished returns
// +1: maxer won
// -1: maxer lost
// 0: not finished
func GameFinished(board *rules.BoardState) float64 {

	if board.Turn == 1 {
		return 0
	}

	maxSnake := board.Snakes[0]
	minSnake := board.Snakes[1]
	maxHead := maxSnake.Body[0]
	minHead := minSnake.Body[0]

	if maxHead.X == minHead.X && maxHead.Y == minHead.Y {
		if len(maxSnake.Body) < len(minSnake.Body) {
			return -1
		}
		if len(maxSnake.Body) == len(minSnake.Body) {
			return -1
		}
		return 1
	}

	if maxSnake.Health == 0 {
		return -1
	}

	if minSnake.Health == 0 {
		return 1
	}

	for _, maxPiece := range maxSnake.Body[1:] {
		if maxPiece.X == maxHead.X && maxPiece.Y == maxHead.Y {
			return -1
		}
		if maxPiece.X == minHead.X && maxPiece.Y == minHead.Y {
			return 1
		}
	}

	// also check the head, but only for max vs min
	if maxSnake.Body[0].X == minHead.X && maxSnake.Body[0].Y == minHead.Y {
		return 1
	}

	for _, minPiece := range minSnake.Body[1:] {
		if minPiece.X == minHead.X && minPiece.Y == minHead.Y {
			return 1
		}
		if minPiece.X == maxHead.X && minPiece.Y == maxHead.Y {
			return -1
		}
	}

	// also check the head, but only for min vs max (improving performance)
	if minSnake.Body[0].X == maxHead.X && minSnake.Body[0].Y == maxHead.Y {
		return -1
	}

	return 0

}

func GameFinishedBits(snake1, snake2 int) float64 {
	if snake1^snake2 != 0 {
		return 1

	}

	return 0
}

func HeuristicAnalysis(board *rules.BoardState) float64 {

	return PercentageOfBoardControlled(board)
}
