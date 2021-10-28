package pather

import (
	"fmt"
)

// func (p PathGrid) PrintBlocked() {

// 	width := len(p)
// 	height := len(p[0])

// 	for y := height - 1; y >= 0; y-- {
// 		for x := 0; x < width; x++ {

// 			fmt.Printf("%2d ", p[x][y].BlockedForTurns)
// 			// if p[x][y].BlockedForTurns {
// 			// 	continue
// 			// }

// 			// fmt.Printf("0 ")
// 		}
// 		fmt.Println()
// 	}
// }

func (p PathGrid) DebugPrint() {

	width := len(p)
	height := len(p[0])

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {

			if p[x][y] == nil {
				fmt.Printf("---:---:--- |")
				continue
			}

			fmt.Printf("%3d:%3d:%3d |", p[x][y].CostFromOrigin, p[x][y].DistanceToTarget, p[x][y].BlockedForTurns)
			// descriptor := "o"
			// if p[x][y].Blocked {
			// 	descriptor = "x"
			// }
			// if generator.SamePoint(start, rules.Point{X: int32(x), Y: int32(y)}) {
			// 	descriptor = "s"
			// }

			// if generator.SamePoint(end, rules.Point{X: int32(x), Y: int32(y)}) {
			// 	descriptor = "e"
			// }

			// fmt.Printf("%s |", descriptor)
		}
		fmt.Println()
	}
}
