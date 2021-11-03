package pather

import (
	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func (p PathGrid) ExploreForLength(currentRoute []rules.Point, targetLength int) []rules.Point {
	if len(currentRoute) == targetLength {
		return currentRoute
	}
	// last value in array is where we're up to
	currentPoint := currentRoute[len(currentRoute)-1]

	width := len(p)
	height := len(p[0])
	neighbours := generator.NeighboursSafe(height, width, rules.Point{X: currentPoint.X, Y: currentPoint.Y})

	var newRoutes [][]rules.Point

	for _, neighbour := range neighbours {
		alreadyExplored := false
		for _, currentRoutePoint := range currentRoute {
			if generator.SamePoint(currentRoutePoint, neighbour) {
				alreadyExplored = true
				break
			}
		}

		if alreadyExplored {
			continue
		}

		if p[neighbour.X][neighbour.Y] == nil {
			p[neighbour.X][neighbour.Y] = &AStarCost{}
		}

		// check if this snake could possibly have travelled to this square
		// (not blocked)
		if p[neighbour.X][neighbour.Y].BlockedForTurns >= int32(len(currentRoute)) {
			continue
		}

		neighbourBestPath := p.ExploreForLength(append(currentRoute, neighbour), targetLength)

		newRoutes = append(newRoutes, neighbourBestPath)

		// shortcircuit the check here. no point exploring all neighbours if we've found a good path
		if len(neighbourBestPath) >= targetLength {
			return neighbourBestPath
		}

		// fmt.Printf("incrementing %+v, %d\n", p[neighbour.X][neighbour.Y], startingBlockedInValue+1)
	}

	if len(newRoutes) == 0 {
		return currentRoute
	}

	var longestRoute []rules.Point
	for _, route := range newRoutes {
		if len(route) > len(longestRoute) {
			longestRoute = route
		}
	}

	return longestRoute
}
