package generator

import "github.com/brensch/snake/rules"

func ChaserNextMove(snake rules.Snake) rules.Direction {

	if snake.Body[0].X > snake.Body[len(snake.Body)-1].X {
		return rules.DirectionLeft
	}
	if snake.Body[0].X < snake.Body[len(snake.Body)-1].X {
		return rules.DirectionRight
	}
	if snake.Body[0].Y > snake.Body[len(snake.Body)-1].Y {
		return rules.DirectionDown
	}
	if snake.Body[0].Y < snake.Body[len(snake.Body)-1].Y {
		return rules.DirectionUp
	}

	return rules.DirectionUnknown
}
