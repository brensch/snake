package pather

import (
	"fmt"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func (p PathGrid) PrintBlocked() {

	width := len(p)
	height := len(p[0])

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {

			fmt.Printf("%2d ", p[x][y].BlockedForTurns)
			// if p[x][y].BlockedForTurns {
			// 	continue
			// }

			// fmt.Printf("0 ")
		}
		fmt.Println()
	}
}

func (p PathGrid) DebugPrint(start, end rules.Point) {

	width := len(p)
	height := len(p[0])

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {

			fmt.Printf("%3d:%3d ", p[x][y].DistanceToOrigin, p[x][y].DistanceToTarget)
			descriptor := "o"
			if p[x][y].Blocked {
				descriptor = "x"
			}
			if generator.SamePoint(start, rules.Point{X: int32(x), Y: int32(y)}) {
				descriptor = "s"
			}

			if generator.SamePoint(end, rules.Point{X: int32(x), Y: int32(y)}) {
				descriptor = "e"
			}

			fmt.Printf("%s |", descriptor)
		}
		fmt.Println()
	}
}
