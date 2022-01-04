package generator

import "github.com/brensch/snake/rules"

// CentreMostPoint calculates the point closest to the middle of the board
func CentreMostPoint(s *rules.BoardState, points []rules.Point) rules.Point {
	centrePoint := rules.Point{X: s.Width / 2, Y: s.Height / 2}

	smallestDistance := byte(255)
	var smallestDistancePoint rules.Point
	for _, point := range points {
		distance := DistanceBetween(centrePoint, point)
		if distance >= smallestDistance {
			continue
		}

		smallestDistance = distance
		smallestDistancePoint = point
	}

	return smallestDistancePoint
}
