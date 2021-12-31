package minimax

import (
	"math"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

// type

func PercentageOfBoardControlled(board *rules.BoardState, you int) float64 {

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

		if closestSnake == you {
			closestSquareCount++
		}

	}

	return (float64(closestSquareCount) / float64(totalSpaces))

}

// TODO: make this actually check the shortest path to the point
func PercentageOfBoardControlledSmart(board *rules.BoardState, you int) float64 {

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

		if closestSnake == you {
			closestSquareCount++
		}

	}

	return (float64(closestSquareCount) / float64(totalSpaces))

}

// GameFinished returns
// +inf: you won
// -inf: you lost
// 0: not finished
func GameFinished(board *rules.BoardState, you int) float64 {

	youSnake := board.Snakes[you]
	opponentSnake := board.Snakes[(you+1)%2]
	youHead := youSnake.Body[0]

	if youSnake.Health == 0 {
		return math.Inf(-1)
	}

	if opponentSnake.Health == 0 {
		return math.Inf(1)
	}

	for _, opponentPiece := range opponentSnake.Body {

		if opponentPiece.X == youHead.X && opponentPiece.Y == youHead.Y {
			return math.Inf(-1)
		}

	}

	for _, youPiece := range youSnake.Body[1:] {

		if youPiece.X == youHead.X && youPiece.Y == youHead.Y {
			return math.Inf(-1)
		}

	}

	return 0

}

func HeuristicAnalysis(board *rules.BoardState, you int) float64 {
	finishedCheck := GameFinished(board, you)
	if finishedCheck != 0 {
		return finishedCheck
	}

	return PercentageOfBoardControlled(board, you)
}
