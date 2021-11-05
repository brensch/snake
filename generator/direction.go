package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

// DirectionToPoint expects OneOff() to be true to You.Head.
func DirectionToPoint(start, end rules.Point) rules.Direction {

	if SamePoint(start, rules.Point{X: end.X - 1, Y: end.Y}) {
		return rules.DirectionRight
	}
	if SamePoint(start, rules.Point{X: end.X + 1, Y: end.Y}) {
		return rules.DirectionLeft
	}
	if SamePoint(start, rules.Point{X: end.X, Y: end.Y - 1}) {
		return rules.DirectionUp
	}
	if SamePoint(start, rules.Point{X: end.X, Y: end.Y + 1}) {
		return rules.DirectionDown
	}

	fmt.Printf("bug in logic %+v %+v\n", start, end)

	return rules.DirectionDown

}
