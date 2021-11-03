package pather

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

func IdentifyKill(s *rules.BoardState, them, you rules.Snake) ([]rules.Point, error) {

	youReachablePoints, youAvailableSquaresGrid := GetReachablePoints(s, you.Body[0], you.ID)

	fmt.Println(youReachablePoints)
	youAvailableSquaresGrid.DebugPrint()

	themReachablePoints, themAvailableSquaresGrid := GetReachablePoints(s, them.Body[0], them.ID)

	fmt.Println(themReachablePoints)
	themAvailableSquaresGrid.DebugPrint()

	killPoints := themAvailableSquaresGrid.FindKillPoints()
	fmt.Println("killpoints:", killPoints)

	// for _, point :=

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

// // FindKillPoints finds all points where a snake could be caught and killed
// func FindKillPointsYouCause(killPoints []rules.Point, youGrid, themGrid PathGrid) []rules.Point {

// 	for _, killPoint := range killPoints {
// 		stepsFromYou := youGrid[killPoint.X][killPoint.Y].StepsFromOrigin
// 		stepsFromTh
// 	}

// 	return points
// }
