package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

type Direction uint8

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionUp
	DirectionDown
	DirectionUnknown

	directionCount = 4
)

func (d Direction) String() string {
	return [...]string{"left", "right", "up", "down", "unknown"}[d]
}

// DirectionToPoint expects OneOff() to be true to You.Head.
func DirectionToPoint(start, end rules.Point) Direction {

	if SamePoint(start, rules.Point{X: end.X - 1, Y: end.Y}) {
		return DirectionRight
	}
	if SamePoint(start, rules.Point{X: end.X + 1, Y: end.Y}) {
		return DirectionLeft
	}
	if SamePoint(start, rules.Point{X: end.X, Y: end.Y - 1}) {
		return DirectionUp
	}
	if SamePoint(start, rules.Point{X: end.X, Y: end.Y + 1}) {
		return DirectionDown
	}

	fmt.Printf("bug in logic %+v %+v\n", start, end)

	return DirectionDown

}

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
