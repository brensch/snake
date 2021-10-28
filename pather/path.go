package pather

import (
	"fmt"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

const (
	CostFactor  = int32(10)
	LargestCost = int32(100000)
)

// PathGrid tracks the available
// is [x][y]
type PathGrid [][]*AStarCost

type AStarCost struct {
	DistanceToTarget int32 // the straight line distance to the target calculated via pythag
	StepsFromOrigin  int32 // the lowest current known number of steps to this square from the origin
	CostFromOrigin   int32 // the lowest current known cost to get to this square from the origin. takes into account hazard cost
	Explored         bool  // we have already calculate this square
	BlockedForTurns  int32 // impassable unless the steps from origin is >= this value
	Hazard           bool  // costs extra health to traverse
	// Origin           bool // this is the square we started at
	// Blocked          bool // This is an impassable square (ie snake)
}

func (p PathGrid) ScoreNeighbours(currentPoint, targetPoint rules.Point, hazardCost int32) bool {
	neighbours := generator.NeighboursSafe(len(p[0]), len(p), currentPoint)
	currentPath := p[currentPoint.X][currentPoint.Y]

	for _, neighbourPoint := range neighbours {

		if p[neighbourPoint.X][neighbourPoint.Y] == nil {
			p[neighbourPoint.X][neighbourPoint.Y] = &AStarCost{}
		}

		neighbourPath := p[neighbourPoint.X][neighbourPoint.Y]

		// we see if by this point in time this
		if neighbourPath.Explored || currentPath.StepsFromOrigin < neighbourPath.BlockedForTurns {
			continue
		}

		// weight this path more heavily if it's hazardous
		cost := CostFactor
		if neighbourPath.Hazard {
			cost += hazardCost
		}

		// update costs and steps independently
		// (steps are critical for calculating obstacle state)
		neighbourPath.CostFromOrigin = currentPath.CostFromOrigin + cost
		neighbourPath.StepsFromOrigin = currentPath.StepsFromOrigin + 1

		// calculate the pythag distance to the targetPoint
		// using only ints and avoiding math library since performance is critical
		xDelta := Abs(neighbourPoint.X-targetPoint.X) * CostFactor
		yDelta := Abs(neighbourPoint.Y-targetPoint.Y) * CostFactor
		neighbourPath.DistanceToTarget = IntSqrt(xDelta*xDelta + yDelta*yDelta)

		// check if we found the target
		if generator.SamePoint(neighbourPoint, targetPoint) {
			return true
		}

	}

	return false
}

// // MakePathgrid performs all the steps to initialise a pathgrid to then be able to calculate routes.
// // origin and yourhead are different because obstacle calculation should not take into account the
// // snake's own head, but we may be initialising path grids not from the snake's own head. For instance
// // the snack calculations
// func MakePathgrid(s *rules.BoardState, origin, yourHead rules.Point) PathGrid {
// 	grid := initPathGrid(s)
// 	grid.AddObstacles(s, origin)
// 	// grid.CalculateDistancesToOrigin(origin)

// 	return grid
// }

func initPathGrid(s *rules.BoardState) PathGrid {
	grid := make(PathGrid, s.Width)

	for x := range grid {
		grid[x] = make([]*AStarCost, s.Height)
		// for y := range grid[x] {
		// 	grid[x][y] = &AStarCost{}
		// }
	}

	// grid[origin.X][origin.Y].Origin = true
	return grid
}

// func (p PathGrid) MarkExplored(coord rules.Point) {
// 	p[coord.X][coord.Y].Explored = true
// }

// Some int maths helpers to hopefully make things a bit faster
func Abs(num int32) int32 {
	if num < 0 {
		return -num
	}
	return num
}

func IntSqrt(n int32) int32 {

	x := n
	y := int32(1)
	for x > y {
		x = (x + y) / 2
		y = n / x
	}

	return x
}

// AddObstacles adds all snake body parts if they will still be there when you get there.
// ie it removes pieces from the ends of snakes depending on how far away it is from origin.
// It also creates a cloud around the head of the snake depending on if they could make it to that
// square theoretically by the time origin gets there.
// To avoid distant snakeheads causing enormous danger clouds, it will ignore them if they are more than
// 6 moves away.
func (p PathGrid) AddObstacles(s *rules.BoardState, origin rules.Point) {

	for _, snake := range s.Snakes {

		// only include snakes segments if they will still be there when we get there
		for pointNumber, point := range snake.Body {
			// distanceToOrigin := Abs(origin.X-point.X) + Abs(origin.Y-point.Y)

			if generator.OffBoard(s, point) {
				continue
			}

			if p[point.X][point.Y] == nil {
				p[point.X][point.Y] = &AStarCost{}
			}

			p[point.X][point.Y].BlockedForTurns = int32(len(snake.Body) - pointNumber - 1)

			// if generator.OffBoard(s, point) ||
			// 	len(snake.Body)-pointPosition <= int(distanceToOrigin) {
			// 	continue
			// }

			// p[point.X][point.Y].Blocked = true
		}
	}

	// // check every point for proximity to origin and a snake head
	// for x := int32(0); x < s.Width; x++ {
	// 	for y := int32(0); y < s.Height; y++ {
	// 		distanceToOrigin := Abs(origin.X-x) + Abs(origin.Y-y)

	// 		// anything more than 5 is too uncertain to be blocking squares because of it
	// 		// TODO: actually subtly weight the path rather than just a binary good/bad
	// 		if distanceToOrigin > 5 {
	// 			continue
	// 		}

	// 		for _, snake := range s.Snakes {
	// 			// use yourhead here to check if this is us
	// 			if generator.SamePoint(snake.Body[0], yourHead) {
	// 				continue
	// 			}
	// 			distanceToOtherSnake := Abs(snake.Body[0].X-x) + Abs(snake.Body[0].Y-y)
	// 			if distanceToOtherSnake < distanceToOrigin {
	// 				p[x][y].Blocked = true

	// 			}
	// 		}

	// 	}
	// }

	for _, hazard := range s.Hazards {
		p[hazard.X][hazard.Y].Hazard = true
	}
}

// func (p PathGrid) AddFogOfUncertainty(head rules.Point, distance int32) {

// 	neighbours := generator.NeighboursSafe(len(p[0]), len(p), head)
// 	for _, neighbour := range neighbours {
// 		p[neighbour.X][neighbour.Y].Blocked = true
// 		p.AddFogOfUncertainty(neighbour, distance-1)
// 	}

// }

func (p PathGrid) GetNextLowestSquare() (rules.Point, error) {

	lowestCost := LargestCost
	var lowestX, lowestY int32
	found := false
	for x, xRow := range p {
		for y := range xRow {
			if p[x][y] == nil || p[x][y].Explored || p[x][y].CostFromOrigin == 0 {
				continue
			}
			cost := p[x][y].CostFromOrigin + p[x][y].DistanceToTarget
			if cost < lowestCost {
				lowestCost = cost
				lowestX = int32(x)
				lowestY = int32(y)
				found = true
			}

		}
	}

	if !found {
		return rules.Point{}, fmt.Errorf("no next square to explore")
	}

	return rules.Point{X: lowestX, Y: lowestY}, nil
}

func (p PathGrid) NextStepBackToOrigin(current, origin rules.Point) (rules.Point, error) {

	var nextCoord rules.Point
	// lowestCost := LargestCost
	// lowestDistanceToTarget := p[current.X][current.Y].DistanceToTarget
	lowestCostFromOrigin := p[current.X][current.Y].CostFromOrigin
	found := false

	neighbours := generator.NeighboursSafe(len(p[0]), len(p), current)
	for _, neighbour := range neighbours {

		// if we didn't explore it it won't be part of the path
		if p[neighbour.X][neighbour.Y] == nil {
			continue
		}

		// stop when we find the origin
		if generator.SamePoint(neighbour, origin) {
			return neighbour, nil
		}

		// you only need to check the distance to target, not the cost. cost is used in calculating the next
		// square to check in the distance to target sum.
		if p[neighbour.X][neighbour.Y].CostFromOrigin < lowestCostFromOrigin {
			lowestCostFromOrigin = p[neighbour.X][neighbour.Y].CostFromOrigin
			nextCoord = neighbour
			found = true
		}

	}

	if !found {
		// p.DebugPrint(current, target)
		return rules.Point{}, fmt.Errorf("no step takes us home")
	}

	return nextCoord, nil

}

// func GetRoutesFromOrigin(state *rules.BoardState, origin, yourHead rules.Point, hazardCost int) [][]rules.Point {
// 	p := MakePathgrid(state, origin, yourHead)
// 	var routes [][]rules.Point

// 	for x := int32(0); x < state.Width; x++ {
// 		for y := int32(0); y < state.Height; y++ {
// 			target := rules.Point{X: x, Y: y}
// 			if p[x][y].Blocked {
// 				continue
// 			}

// 			route, _, err := p.GetRoute(origin, target, hazardCost)
// 			if err != nil {
// 				continue
// 			}

// 			routes = append(routes, route)

// 		}
// 	}

// 	return routes
// }

// func (p PathGrid) FreeSquares(state *rules.BoardState) int {
// 	var freeSquares int
// 	for x := int32(0); x < state.Width; x++ {
// 		for y := int32(0); y < state.Height; y++ {
// 			if !p[x][y].Blocked {
// 				freeSquares++
// 			}
// 		}
// 	}
// 	return freeSquares
// }

// TraceWeightedGridToTarget takes a weighted grid and finds the best next step on the route from origin to target
func (p PathGrid) TraceRouteBackToOrigin(origin, target rules.Point) ([]rules.Point, error) {

	currentPoint := target
	route := []rules.Point{currentPoint}

	// if !p[origin.X][origin.Y].Blocked && p[origin.X][origin.Y].DistanceToTarget == 0 {
	// 	return nil, fmt.Errorf("no path to target origin: %+v target %+v", origin, target)
	// }

	for {
		nextPoint, err := p.NextStepBackToOrigin(currentPoint, origin)
		if err != nil {
			return nil, err
		}

		if generator.SamePoint(nextPoint, origin) {
			break
		}

		route = append(route, nextPoint)
		currentPoint = nextPoint

		// if len(route) > 50 {
		// 	// log.WithFields(log.Fields{
		// 	// 	"pathgrid":   p,
		// 	// 	"origin":     origin,
		// 	// 	"target":     target,
		// 	// 	"next_point": nextPoint,
		// 	// }).Error("finding path entered a loop. indicates an issue with a* algorithm")
		// 	fmt.Println(route)
		// 	p.DebugPrint()
		// 	return route, fmt.Errorf("got too many points in route.")
		// }

	}

	return route, nil

}

// // CalculateDistancesToOrigin performs the check of pythagorean distance from the origin to all other points on the board.
// // Since this can be reused it reduces iterations to perform this once on all squares since we will be checking
// // routes to multiple points from a single origin.
// // It can be completed in O(N) time
// func (p PathGrid) CalculateDistancesToOrigin(origin rules.Point) {

// 	for x := 0; x < len(p); x++ {
// 		xDelta := origin.X - int32(x)
// 		for y := 0; y < len(p[x]); y++ {
// 			yDelta := origin.Y - int32(y)
// 			p[x][y].DistanceToOrigin = int(math.Sqrt(float64(xDelta*xDelta)+float64(yDelta*yDelta)) * CostFactor)
// 		}
// 	}
// }

// CalculateDistancesToOrigin returns the shortest path from origin to target
func (p PathGrid) CalculateRouteToTarget(origin, target rules.Point, hazardCost int32) error {

	var err error
	nextSquare := origin

	// cannot guarantee origin is initialised
	// (for instance in case of calculating from snack onwards)
	if p[nextSquare.X][nextSquare.Y] == nil {
		p[nextSquare.X][nextSquare.Y] = &AStarCost{}
	}

	for {
		p[nextSquare.X][nextSquare.Y].Explored = true
		foundTarget := p.ScoreNeighbours(nextSquare, target, hazardCost)
		if foundTarget {
			return nil
		}

		nextSquare, err = p.GetNextLowestSquare()
		if err != nil {
			return err
		}

	}

}

// func (p PathGrid) ResetDistancesToTarget() {
// 	for _, xColumn := range p {
// 		for _, y := range xColumn {
// 			y.DistanceToTarget = 0
// 			y.Explored = false
// 		}
// 	}
// }

// func (p PathGrid) GetRoute(origin, target rules.Point, hazardCost int) ([]rules.Point, int, error) {
// 	p.CalculateRouteToTarget(origin, target, hazardCost)
// 	healthCostToGetToSnack := p[origin.X][origin.Y]

// 	defer p.ResetDistancesToTarget()
// 	route, err := p.TraceWeightedGridToTarget(origin, target)
// 	return route, healthCostToGetToSnack.DistanceToTarget, err
// }

func GetRoute(s *rules.BoardState, ruleset rules.Ruleset, origin, target rules.Point) ([]rules.Point, int32, error) {

	grid := initPathGrid(s)
	grid.AddObstacles(s, origin)
	hazardCost := int32(0)

	switch r := ruleset.(type) {
	case *rules.RoyaleRuleset:
		hazardCost = r.HazardDamagePerTurn
	case *rules.SquadRuleset:
		hazardCost = r.HazardDamagePerTurn

	}

	//
	err := grid.CalculateRouteToTarget(origin, target, hazardCost)
	if err != nil {
		return nil, 0, err
	}

	route, err := grid.TraceRouteBackToOrigin(origin, target)
	if err != nil {
		// shouldn't be getting an error here tbh, should be caught above
		return nil, 0, err
	}

	return route, grid[target.X][target.Y].CostFromOrigin, nil
}

// TODO: make one which also calculates the cost based on hazard
func CountSquaresReachableFromOrigin(s *rules.BoardState, origin rules.Point) int32 {

	grid := initPathGrid(s)
	grid.AddObstacles(s, origin)

	// making a fake target outside the board does not lead to any illegal array dereferences, but does
	// ensure the algorithm will exhaust all squares available to it before giving up
	target := rules.Point{X: -1, Y: -1}
	grid.CalculateRouteToTarget(origin, target, 0)

	count := int32(0)

	for x := int32(0); x < s.Width; x++ {
		for y := int32(0); y < s.Height; y++ {
			if grid[x][y] != nil && grid[x][y].CostFromOrigin > 0 {
				count++
			}
		}
	}

	return count

}
