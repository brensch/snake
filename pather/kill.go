package pather

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

func IdentifyKill(s *rules.BoardState, them, you rules.Snake) ([]rules.Point, error) {

	youReachablePoints, youGrid := GetReachablePoints(s, you.Body[0], you.ID)

	fmt.Println(youReachablePoints)
	youGrid.DebugPrint()

	themReachablePoints, themGrid := GetReachablePoints(s, them.Body[0], them.ID)

	fmt.Println(themReachablePoints)
	themGrid.DebugPrint()

	killPoints := themGrid.FindKillPoints()
	fmt.Println("killpoints:", killPoints)

	causedPoints := FindKillPointsYouCause(killPoints, youGrid, themGrid)

	fmt.Println("causedpoints", causedPoints)

	return nil, fmt.Errorf("in need of work")

}

// FindKillPoints finds all points where a snake could be caught and killed
func (p PathGrid) FindKillPoints() []rules.Point {

	width := len(p)
	height := len(p[0])

	var points []rules.Point

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// when the best number of turns you can reach a square is the same number of turns someone else can,
			// bad times
			if p[x][y] != nil && p[x][y].BlockedInTurns != 0 &&
				p[x][y].StepsFromOrigin == p[x][y].BlockedInTurns {
				points = append(points, rules.Point{X: int32(x), Y: int32(y)})
			}
		}
	}

	return points
}

// FindKillPoints finds all points where a snake could be caught and killed
func FindKillPointsYouCause(killPoints []rules.Point, youGrid, themGrid PathGrid) []rules.Point {

	var causedPoints []rules.Point
	for _, killPoint := range killPoints {
		stepsFromYou := youGrid.At(killPoint).StepsFromOrigin
		stepsFromThem := themGrid.At(killPoint).StepsFromOrigin
		if stepsFromThem == stepsFromYou {
			causedPoints = append(causedPoints, killPoint)
		}
	}

	return causedPoints
}

func (p PathGrid) GetFurthestReachablePoint(reachablePoints []rules.Point) rules.Point {
	largestStepCount := int32(0)
	var largestStepPoint rules.Point
	for _, point := range reachablePoints {
		if p.At(point).StepsFromOrigin > largestStepCount {
			largestStepCount = p.At(point).StepsFromOrigin
			largestStepPoint = point
		}
	}

	return largestStepPoint
}
