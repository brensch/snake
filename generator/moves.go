package generator

import (
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
	// generate the headsnaps for each snake
	smoothBrainArray := make([][]rules.Direction, len(state.Snakes))

	unsafeArray := make([][]bool, state.Width)
	for column := range unsafeArray {
		unsafeArray[column] = make([]bool, state.Height)
	}

	for _, snake := range state.Snakes {
		// check all except the last piece (will have moved by next turn)
		for _, bodyPiece := range snake.Body[:len(snake.Body)-1] {
			unsafeArray[bodyPiece.X][bodyPiece.Y] = true
		}
	}

	for snakePosition, snake := range state.Snakes {

		for direction := 0; direction < 4; direction++ {
			switch rules.Direction(direction) {
			case rules.DirectionRight:
				if snake.Body[0].X == state.Width-1 || unsafeArray[snake.Body[0].X+1][snake.Body[0].Y] {
					smoothBrainArray[snakePosition] = append(smoothBrainArray[snakePosition], rules.DirectionRight)
				}
			case rules.DirectionLeft:
				if snake.Body[0].X == 0 || unsafeArray[snake.Body[0].X-1][snake.Body[0].Y] {
					smoothBrainArray[snakePosition] = append(smoothBrainArray[snakePosition], rules.DirectionLeft)
				}
			case rules.DirectionUp:
				if snake.Body[0].Y == state.Height-1 || unsafeArray[snake.Body[0].X][snake.Body[0].Y+1] {
					smoothBrainArray[snakePosition] = append(smoothBrainArray[snakePosition], rules.DirectionUp)
				}
			case rules.DirectionDown:
				if snake.Body[0].Y == 0 || unsafeArray[snake.Body[0].X][snake.Body[0].Y-1] {
					smoothBrainArray[snakePosition] = append(smoothBrainArray[snakePosition], rules.DirectionDown)
				}
			}
		}
	}
	// fmt.Println(smoothBrainArray)
	// totalMoves := 4 * len(state.Snakes)
	totalMoves := 1
	for i := 0; i < len(state.Snakes); i++ {
		totalMoves = totalMoves * 4
	}
	var moves [][]rules.SnakeMove
	for i := 0; i < totalMoves; i++ {
		moveSet := baseConvert(i, state.Snakes)
		isSmooth := false
		// fmt.Println(moveSet, smoothBrainArray)
		for snake, move := range moveSet {
			for _, smoothMove := range smoothBrainArray[snake] {

				if move.Move == smoothMove {
					// fmt.Println("got smoth", move.Move, smoothBrainArray[snake])
					isSmooth = true
					break
				}
			}
		}
		if isSmooth {
			continue
		}
		// fmt.Println("galaxy")

		moves = append(moves, moveSet)

	}
	// fmt.Println("got ", len(moves))
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
