package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

// AllMovesForState makes a list of move permutations excluding 75% headjobs and out of bounds
func AllMovesForState(state *rules.BoardState) [][]rules.Point {

	neighbourArray := make([][]rules.Point, len(state.Snakes))
	for snakeNumber, snake := range state.Snakes {
		neighbourArray[snakeNumber] = NeighboursSnake(state, snake.Body)
	}

	totalOptions := 1
	for _, neighbours := range neighbourArray {
		totalOptions = totalOptions * len(neighbours)
	}

	allMoves := make([][]rules.Point, totalOptions)

	for i := 0; i < totalOptions; i++ {
		divider := 1
		for j := 0; j < len(state.Snakes); j++ {
			allMoves[i] = append(allMoves[i], neighbourArray[j][i/divider%len(neighbourArray[j])])
			divider = divider * len(neighbourArray[j])
		}
	}

	return allMoves
}

func AllMovesForStateRaw(state *rules.BoardState) [][]rules.SnakeMove {
	// totalMoves := 4 * len(state.Snakes)
	totalMoves := 1
	for i := 0; i < len(state.Snakes); i++ {
		fmt.Println(totalMoves)
		totalMoves = totalMoves * 4
	}
	fmt.Println(totalMoves)
	var moves [][]rules.SnakeMove
	for i := 0; i < totalMoves; i++ {
		moves = append(moves, baseConvert(i, state.Snakes))

	}
	return moves
}

func baseConvert(x int, snakes []rules.Snake) []rules.SnakeMove {
	r := []rules.SnakeMove{}
	for _, snake := range snakes {
		r = append(r, rules.SnakeMove{
			ID:   snake.ID,
			Move: rules.Direction(x % 4),
		})

		x /= 4
	}
	return r
}
