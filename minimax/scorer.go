package minimax

import (
	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

// type

func PercentageOfBoardControlled(board *rules.BoardState, you string) float64 {

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
		closestSnake := ""
		for j, snake := range board.Snakes {
			if allMovesFromSquares[j][i] < int32(closestDistance) {
				closestDistance = allMovesFromSquares[j][i]
				closestSnake = snake.ID
			}
		}

		if closestSnake == you {
			closestSquareCount++
		}

	}

	return float64(closestSquareCount) / float64(totalSpaces)

}
