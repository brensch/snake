package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

// AllMovesForState makes a list of move permutations excluding 75% headjobs and out of bounds
func AllMoveSetsForState(state *rules.BoardState) [][]rules.Point {

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

func AllMoveSetsForStateRaw(state *rules.BoardState) [][]rules.SnakeMove {
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

func AllMovesForSnake(b *rules.BoardState, snakePosition int) [4]bool {
	// get their neck
	head := b.Snakes[snakePosition].Body[0]
	neck := b.Snakes[snakePosition].Body[1]

	positions := [4]bool{}

	// go through each direction, check if it's ok

	// left
	if head.X > 0 && neck.X >= head.X {
		positions[rules.DirectionLeft] = true
	}

	// right
	if head.X < b.Width-1 && neck.X <= head.X {
		positions[rules.DirectionRight] = true
	}

	// down
	if head.Y > 0 && neck.Y >= head.Y {
		positions[rules.DirectionDown] = true
	}

	// up
	if head.Y < b.Height-1 && neck.Y <= head.Y {
		positions[rules.DirectionUp] = true
	}

	return positions

}

func AllMovesForSnakeSlow(b *rules.BoardState, snakePosition int) []rules.Direction {
	// get their neck
	head := b.Snakes[snakePosition].Body[0]
	neck := b.Snakes[snakePosition].Body[1]

	var positions []rules.Direction

	// go through each direction, check if it's ok

	// left
	if head.X > 0 && neck.X >= head.X {
		positions = append(positions, rules.DirectionLeft)
	}

	// right
	if head.X < b.Width-1 && neck.X <= head.X {
		positions = append(positions, rules.DirectionRight)
	}

	// down
	if head.Y > 0 && neck.Y >= head.Y {
		positions = append(positions, rules.DirectionDown)
	}

	// up
	if head.Y < b.Height-1 && neck.Y <= head.Y {
		positions = append(positions, rules.DirectionUp)
	}

	return positions

}
