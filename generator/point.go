package generator

import (
	"math"

	"github.com/brensch/snake/rules"
)

func Neighbours(point rules.Point) [4]rules.Point {

	return [4]rules.Point{left(point), right(point), up(point), down(point)}
}

// NeighboursSafe finds all the neighbours that are within the bounds of the game board.
func NeighboursSafe(height, width int, point rules.Point) []rules.Point {
	safeNeighbours := []rules.Point{}

	if point.X > 0 {
		safeNeighbours = append(safeNeighbours, rules.Point{X: point.X - 1, Y: point.Y})
	}

	if int(point.X) < width-1 {
		safeNeighbours = append(safeNeighbours, rules.Point{X: point.X + 1, Y: point.Y})
	}

	if point.Y > 0 {
		safeNeighbours = append(safeNeighbours, rules.Point{X: point.X, Y: point.Y - 1})
	}

	if int(point.Y) < height-1 {
		safeNeighbours = append(safeNeighbours, rules.Point{X: point.X, Y: point.Y + 1})
	}

	return safeNeighbours
}

// NeighboursSnake gets the neighbours viable for a snake (ie no headsnaps)
func NeighboursSnake(state *rules.BoardState, snake []rules.Point) []rules.Point {

	var possibleNeighbours []rules.Point

	for i := 0; i < 4; i++ {
		possibleNeighbour := Move(rules.Direction(i), snake[0])
		if SamePoint(possibleNeighbour, snake[1]) || OffBoard(state, possibleNeighbour) {
			continue
		}

		possibleNeighbours = append(possibleNeighbours, possibleNeighbour)
	}

	return possibleNeighbours
}

func Move(direction rules.Direction, p rules.Point) rules.Point {
	switch direction {
	case rules.DirectionLeft:
		return left(p)
	case rules.DirectionRight:
		return right(p)
	case rules.DirectionUp:
		return up(p)
	case rules.DirectionDown:
		return down(p)
	// if we get a rules.DirectionUnknown in here we should panic. indicates a bug in logic.
	case rules.DirectionUnknown:
		panic("got direction unknown")
	}

	return rules.Point{}
}

func up(p rules.Point) rules.Point {
	return rules.Point{X: p.X, Y: p.Y + 1}
}

func down(p rules.Point) rules.Point {
	return rules.Point{X: p.X, Y: p.Y - 1}
}

func left(p rules.Point) rules.Point {
	return rules.Point{X: p.X - 1, Y: p.Y}
}

func right(p rules.Point) rules.Point {
	return rules.Point{X: p.X + 1, Y: p.Y}
}

// SamePoint checks if two points have the same location
func SamePoint(p1, p2 rules.Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

// OneOff checks if two points are a single snake move away from each other
func OneOff(p1, p2 rules.Point) bool {
	return math.Abs(float64(p1.X-p2.X))+math.Abs(float64(p1.Y-p2.Y)) == 1
}

func OffBoard(state *rules.BoardState, p rules.Point) bool {
	return p.X < 0 ||
		p.X >= state.Width ||
		p.Y < 0 ||
		p.Y >= state.Height
}

// Headsnaps tells you which way you should avoid if you don't want 75% headjob
func Headsnaps(state *rules.BoardState) []rules.Direction {

	directions := make([]rules.Direction, len(state.Snakes))
	for snakePosition, snake := range state.Snakes {

		// assume snake always > 3
		if snake.Body[0].X > snake.Body[1].X {
			directions[snakePosition] = rules.DirectionLeft
			continue
		}
		if snake.Body[0].X < snake.Body[1].X {
			directions[snakePosition] = rules.DirectionRight
			continue
		}
		if snake.Body[0].Y > snake.Body[1].Y {
			directions[snakePosition] = rules.DirectionDown
			continue
		}
		if snake.Body[0].X > snake.Body[1].X {
			directions[snakePosition] = rules.DirectionUp
			continue
		}

		directions[snakePosition] = rules.DirectionUnknown
	}

	return directions
}

func DistanceBetween(a, b rules.Point) int32 {
	return Abs(a.X-b.X) + Abs(a.Y-b.Y)
}
