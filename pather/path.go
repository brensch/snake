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
	BlockedInTurns   int32 // calculate where snakes may be in a given number of turns

	Hazard bool // costs extra health to traverse
	// Origin           bool // this is the square we started At
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

		// if neighbourPoint.X == 9 && neighbourPoint.Y == 5 {
		// 	fmt.Println("thing", neighbourPath.BlockedInTurns, neighbourPath.StepsFromOrigin, currentPoint)
		// }

		// check if this point is blocked by this point in our journey
		if neighbourPath.Explored ||
			currentPath.StepsFromOrigin < neighbourPath.BlockedForTurns ||
			(neighbourPath.BlockedInTurns != 0 && neighbourPath.BlockedInTurns <= currentPath.StepsFromOrigin+1) { // this calculates if another snake could boop us here
			continue
		}

		// // check if this point is blocked by this point in our journey
		// if neighbourPath.Explored ||
		// 	currentPath.StepsFromOrigin < neighbourPath.BlockedForTurns {
		// 	continue
		// }

		// weight this path more heavily if it's hazardous
		cost := CostFactor
		if neighbourPath.Hazard {
			cost += hazardCost
		}

		// update costs and steps independently
		// (steps are critical for calculating obstacle state)
		potentialCostFromOrigin := currentPath.CostFromOrigin + cost
		if neighbourPath.CostFromOrigin == 0 || neighbourPath.CostFromOrigin > potentialCostFromOrigin {
			neighbourPath.CostFromOrigin = potentialCostFromOrigin
		}

		potentialStepsFromOrigin := currentPath.StepsFromOrigin + 1
		if neighbourPath.StepsFromOrigin == 0 || neighbourPath.StepsFromOrigin > potentialStepsFromOrigin {
			neighbourPath.StepsFromOrigin = potentialStepsFromOrigin
		}

		// calculate the pythag distance to the targetPoint
		// using only ints and avoiding math library since performance is critical
		xDelta := generator.Abs(neighbourPoint.X-targetPoint.X) * CostFactor
		yDelta := generator.Abs(neighbourPoint.Y-targetPoint.Y) * CostFactor
		neighbourPath.DistanceToTarget = generator.IntSqrt(xDelta*xDelta + yDelta*yDelta)

		// check if we found the target
		if generator.SamePoint(neighbourPoint, targetPoint) {
			return true
		}

	}

	return false
}

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

func (p PathGrid) CalculatePointNeighbourBlockedInValues(x, y, startX, startY int32) []rules.Point {

	fmt.Println(x, y)
	if x < 0 || y < 0 {
		fmt.Println("out of bounds----------------", x, y)
		return nil
	}

	if p[x][y] == nil {
		fmt.Println("error in logic calculating neighbour blocked in values")
	}

	startingBlockedInValue := p[x][y].BlockedInTurns

	width := len(p)
	height := len(p[0])
	neighbours := generator.NeighboursSafe(height, width, rules.Point{X: x, Y: y})

	var newNeighbours []rules.Point
	// fmt.Println("got neighbours", neighbours)

	for _, neighbour := range neighbours {

		// we should not calculate a value for the start of the snake
		// (ie when it loops back on itself in the second iteration)
		if startX == neighbour.X && startY == neighbour.Y {
			continue
		}

		if p[neighbour.X][neighbour.Y] == nil {
			p[neighbour.X][neighbour.Y] = &AStarCost{}
		}

		// only want to update value if the current blockedin count is smaller
		// (ie this snake can get there quicker)
		if p[neighbour.X][neighbour.Y].BlockedInTurns != 0 &&
			p[neighbour.X][neighbour.Y].BlockedInTurns < startingBlockedInValue+1 {
			continue
		}

		// check if this snake could possibly have travelled to this square
		// (not blocked)
		if p[neighbour.X][neighbour.Y].BlockedForTurns >= startingBlockedInValue {
			continue
		}

		newNeighbours = append(newNeighbours, neighbour)
		p[neighbour.X][neighbour.Y].BlockedInTurns = startingBlockedInValue + 1
		// fmt.Printf("incrementing %+v, %d\n", p[neighbour.X][neighbour.Y], startingBlockedInValue+1)
	}

	return newNeighbours
}

// AddObstacles adds all snake body parts if they will still be there when you get there.
// ie it removes pieces from the ends of snakes depending on how far away it is from origin.
// It also creates a cloud around the head of the snake depending on if they could make it to that
// square theoretically by the time origin gets there.
// To avoid distant snakeheads causing enormous danger clouds, it will ignore them if they are more than
// 6 moves away.
func (p PathGrid) AddObstacles(s *rules.BoardState, origin rules.Point, youID string) {

	you, err := generator.GetYou(s, youID)
	if err != nil {
		fmt.Println("there is no you snake in this board")
	}

	for _, snake := range s.Snakes {

		// get array of distances to snacks
		var distancesToSnacks []int32
		for _, snack := range s.Food {
			distancesToSnacks = append(distancesToSnacks,
				generator.Abs(snack.X-snake.Body[0].X)+generator.Abs(snack.Y-snake.Body[0].Y))
		}

		// only include snakes segments if they will still be there when we get there
		for pointNumber, point := range snake.Body {
			// distanceToOrigin := Abs(origin.X-point.X) + Abs(origin.Y-point.Y)

			if generator.OffBoard(s, point) {
				continue
			}

			if p[point.X][point.Y] == nil {
				p[point.X][point.Y] = &AStarCost{}
			}
			// fmt.Println(distancesToSnacks)

			// figure out how many snacks this snake may have encountered by this step
			potentialSnacks := 0
			if snake.ID != youID {
				for _, potentialSnackDistance := range distancesToSnacks {
					if int(potentialSnackDistance) <= len(snake.Body)-pointNumber {
						potentialSnacks++
					}
				}
			}

			// calculate how many turns we're blocked for based on the position of the block in the snake
			// and the length of the snake
			newBlockedForTurns := int32(len(snake.Body) - pointNumber - 1 + potentialSnacks)
			if p[point.X][point.Y].BlockedForTurns < newBlockedForTurns {
				p[point.X][point.Y].BlockedForTurns = newBlockedForTurns
			}

			// if generator.OffBoard(s, point) ||
			// 	len(snake.Body)-pointPosition <= int(distanceToOrigin) {
			// 	continue
			// }

			// p[point.X][point.Y].Blocked = true
		}

	}

	// do second loop to make sure all blocked in values are calculated
	for _, snake := range s.Snakes {
		// do not calculate the future state of a snake if it's you
		if snake.ID == youID {
			continue
		}

		// doing this above seems to lower life expectency. TODO: investigate why
		if snake.EliminatedCause != "" {
			continue
		}

		// if snake.Body[0].X < 0 || snake.Body[0].Y < 0 {
		// 	fmt.Println("got bad snake head", snake.Body[0])
		// }

		// TODO: figure out how many moves ahead i should add this
		pointsToCheck := []rules.Point{snake.Body[0]}

		// setting the starting blocked in value to 1 less since the head won't block you if they're a smaller snake
		startingBlockedInValue := int32(0)
		if len(snake.Body) < len(you.Body) {
			startingBlockedInValue = 1
		}

		if p[snake.Body[0].X][snake.Body[0].Y] == nil {
			p[snake.Body[0].X][snake.Body[0].Y] = &AStarCost{}
		}

		p[snake.Body[0].X][snake.Body[0].Y].BlockedInTurns = startingBlockedInValue

		//
		for i := 0; i < 10; i++ {
			var nextPointsToCheck []rules.Point
			for _, pointToCheck := range pointsToCheck {
				points := p.CalculatePointNeighbourBlockedInValues(pointToCheck.X, pointToCheck.Y, snake.Body[0].X, snake.Body[0].Y)
				nextPointsToCheck = append(nextPointsToCheck, points...)
			}
			if len(nextPointsToCheck) == 0 {
				break
			}
			pointsToCheck = nextPointsToCheck

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
		if p[hazard.X][hazard.Y] == nil {
			p[hazard.X][hazard.Y] = &AStarCost{}
		}
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

		// stop when we find the origin if we are 0 steps from origin
		if p[current.X][current.Y].StepsFromOrigin == 1 && generator.SamePoint(neighbour, origin) {
			return neighbour, nil
		}

		// find the smallest CostFromOrigin.
		// CostFromOrigin needs to not be zero since 0 implies that it was not explored.
		// Each step has to have 1 less StepsFromOrigin or it was not possible to get to that square At that time.
		if p[neighbour.X][neighbour.Y].CostFromOrigin > 0 &&
			p[neighbour.X][neighbour.Y].CostFromOrigin < lowestCostFromOrigin &&
			p[neighbour.X][neighbour.Y].StepsFromOrigin == p[current.X][current.Y].StepsFromOrigin-1 {
			lowestCostFromOrigin = p[neighbour.X][neighbour.Y].CostFromOrigin
			nextCoord = neighbour
			found = true
		}

	}

	if !found {
		// p.DebugPrint(current, target)
		return rules.Point{}, fmt.Errorf("no step takes us home from point %+v", current)
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
		// fmt.Println(nextSquare)

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

// GetRoute calculates a path between two points given a game state, mutating other snakes in a way that seems reasonable
// for each step.
// returns route, costToTarget, err
func GetRoute(s *rules.BoardState, ruleset rules.Ruleset, origin, target rules.Point, youID string) ([]rules.Point, PathGrid, error) {

	grid := initPathGrid(s)
	grid.AddObstacles(s, origin, youID)
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
		return nil, grid, err
	}

	route, err := grid.TraceRouteBackToOrigin(origin, target)
	if err != nil {
		// shouldn't be getting an error here tbh, should be caught above
		return nil, grid, err
	}

	return route, grid, nil
}

func CalculateAllAvailableSquares(s *rules.BoardState, origin rules.Point, youID string) PathGrid {
	grid := initPathGrid(s)
	grid.AddObstacles(s, origin, youID)

	// making a fake target outside the board does not lead to any illegal array dereferences, but does
	// ensure the algorithm will exhaust all squares available to it before giving up
	target := rules.Point{X: -1, Y: -1}
	grid.CalculateRouteToTarget(origin, target, 0)

	// grid.DebugPrint()
	return grid
}

// // TODO: make one which also calculates the cost based on hazard
// func CountSquaresReachableFromOrigin(s *rules.BoardState, origin rules.Point, youID string) int32 {

// 	grid := CalculateAllAvailableSquares(s, origin, youID)

// 	count := int32(0)

// 	for x := int32(0); x < s.Width; x++ {
// 		for y := int32(0); y < s.Height; y++ {
// 			if grid[x][y] != nil && grid[x][y].CostFromOrigin > 0 {
// 				count++
// 			}
// 		}
// 	}

// 	return count

// }

func GetReachablePoints(s *rules.BoardState, origin rules.Point, youID string) ([]rules.Point, PathGrid) {
	grid := CalculateAllAvailableSquares(s, origin, youID)
	// grid.DebugPrint()
	points := []rules.Point{}

	// TODO: don't smoothbrain this sorting algorithm
	// for {
	// 	var longestPath int32
	// 	var longestPathPoint rules.Point
	// 	found := false

	for x := int32(0); x < s.Width; x++ {
		for y := int32(0); y < s.Height; y++ {
			if grid[x][y] != nil && grid[x][y].StepsFromOrigin != 0 {
				points = append(points, rules.Point{X: x, Y: y})

				// longestPath = grid[x][y].StepsFromOrigin
				// longestPathPoint = rules.Point{X: x, Y: y}
				// found = true
			}
		}
	}
	return points, grid

	// 	points = append(points, longestPathPoint)
	// 	grid[longestPathPoint.X][longestPathPoint.Y].StepsFromOrigin = 0
	// }

}

func (p PathGrid) FurthestPoint() rules.Point {
	width := int32(len(p))
	height := int32(len(p[0]))

	furthestDistance := int32(0)
	var furthestDistancePoint rules.Point

	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {

			if p[x][y] == nil || p[x][y].StepsFromOrigin <= furthestDistance {
				continue
			}

			furthestDistance = p[x][y].StepsFromOrigin
			furthestDistancePoint = rules.Point{X: x, Y: y}

		}
	}

	return furthestDistancePoint
}

func (p PathGrid) At(point rules.Point) *AStarCost {
	return p[point.X][point.Y]
}

// func (p PathGrid) FindPointsBeyond(currentRoute []rules.Point, targetLength int) []rules.Point {

// 	// last value in array is where we're up to
// 	currentPoint := currentRoute[len(currentRoute)-1]

// 	width := len(p)
// 	height := len(p[0])
// 	neighbours := generator.NeighboursSafe(height, width, rules.Point{X: currentPoint.X, Y: currentPoint.Y})

// 	var newRoutes [][]rules.Point

// 	for _, neighbour := range neighbours {
// 		alreadyExplored := false
// 		for _, currentRoutePoint := range currentRoute {
// 			if generator.SamePoint(currentRoutePoint, neighbour) {
// 				alreadyExplored = true
// 				break
// 			}
// 		}

// 		if alreadyExplored {
// 			continue
// 		}

// 		if p[neighbour.X][neighbour.Y] == nil {
// 			p[neighbour.X][neighbour.Y] = &AStarCost{}
// 		}

// 		// check if this snake could possibly have travelled to this square
// 		// (not blocked)
// 		if p[neighbour.X][neighbour.Y].BlockedForTurns >= int32(len(currentRoute)) {
// 			continue
// 		}

// 		neighbourBestPath := p.ExploreForLength(append(currentRoute, neighbour), targetLength)

// 		newRoutes = append(newRoutes, neighbourBestPath)

// 		// shortcircuit the check here. no point exploring all neighbours if we've found a good path
// 		if len(neighbourBestPath) >= targetLength {
// 			return neighbourBestPath
// 		}

// 		// fmt.Printf("incrementing %+v, %d\n", p[neighbour.X][neighbour.Y], startingBlockedInValue+1)
// 	}

// 	if len(newRoutes) == 0 {
// 		return currentRoute
// 	}

// 	var longestRoute []rules.Point
// 	for _, route := range newRoutes {
// 		if len(route) > len(longestRoute) {
// 			longestRoute = route
// 		}
// 	}

// 	return longestRoute
// }
