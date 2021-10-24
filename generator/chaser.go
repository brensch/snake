package generator

import "github.com/brensch/snake/rules"

func ChaserNextMove(snake rules.Snake) Direction {

	if snake.Body[0].X > snake.Body[len(snake.Body)-1].X {
		return DirectionLeft
	}
	if snake.Body[0].X < snake.Body[len(snake.Body)-1].X {
		return DirectionRight
	}
	if snake.Body[0].Y > snake.Body[len(snake.Body)-1].Y {
		return DirectionDown
	}
	if snake.Body[0].Y < snake.Body[len(snake.Body)-1].Y {
		return DirectionUp
	}

	return DirectionUnknown
}
