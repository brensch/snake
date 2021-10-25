package pather

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

const (
	CostFactor  = 10
	LargestCost = 100000
)

// PathGrid tracks the available
// is [x][y]
type PathGrid [][]*AStarCost

type AStarCost struct {
	DistanceToTarget int // the cost of the best current route between this point and the target
	DistanceToOrigin int // this is static and calculated for all points. it's pythag * cost factor.
	Origin           bool
	Explored         bool // we have already calculate this square
	Blocked          bool // This is an impassable square (ie snake)
}

func (p PathGrid) ScoreNeighbours(current rules.Point) {
	neighbours := generator.NeighboursSafe(len(p[0]), len(p), current)
	currentDistanceToTarget := p[current.X][current.Y].DistanceToTarget
	for _, neighbour := range neighbours {
		if p[neighbour.X][neighbour.Y].Blocked || p[neighbour.X][neighbour.Y].Explored {
			continue
		}
		p[neighbour.X][neighbour.Y].DistanceToTarget = currentDistanceToTarget + CostFactor
	}
}

// MakePathgrid performs all the steps to initialise a pathgrid to then be able to calculate routes.
// origin and yourhead are different because obstacle calculation should not take into account the
// snake's own head, but we may be initialising path grids not from the snake's own head. For instance
// the snack calculations
func MakePathgrid(s *rules.BoardState, origin, yourHead rules.Point) PathGrid {
	grid := initPathGrid(s, origin)
	grid.AddObstacles(s, origin, yourHead)
	grid.CalculateDistancesToOrigin(origin)

	return grid
}

func initPathGrid(s *rules.BoardState, origin rules.Point) PathGrid {
	grid := make(PathGrid, s.Width)

	for x := range grid {
		grid[x] = make([]*AStarCost, s.Height)
		for y := range grid[x] {
			grid[x][y] = &AStarCost{}
		}
	}

	grid[origin.X][origin.Y].Origin = true
	return grid
}

// func (p PathGrid) MarkExplored(coord rules.Point) {
// 	p[coord.X][coord.Y].Explored = true
// }

func Abs(num int32) int32 {
	if num < 0 {
		return -num
	}
	return num
}

// AddObstacles adds all snake body parts if they will still be there when you get there.
// ie it removes pieces from the ends of snakes depending on how far away it is from origin.
// It also creates a cloud around the head of the snake depending on if they could make it to that
// square theoretically by the time origin gets there.
// To avoid distant snakeheads causing enormous danger clouds, it will ignore them if they are more than
// 6 moves away.
func (p PathGrid) AddObstacles(s *rules.BoardState, origin, yourHead rules.Point) {

	for _, snake := range s.Snakes {

		// only include snakes segments if they will still be there when we get there
		for pointPosition, point := range snake.Body {
			distanceToOrigin := Abs(origin.X-point.X) + Abs(origin.Y-point.Y)

			if len(snake.Body)-pointPosition > int(distanceToOrigin) {
				p[point.X][point.Y].Blocked = true
			}
		}
	}

	// check every point for proximity to origin and a snake head
	for x := int32(0); x < s.Width; x++ {
		for y := int32(0); y < s.Height; y++ {
			distanceToOrigin := Abs(origin.X-x) + Abs(origin.Y-y)

			// anything more than 5 is too uncertain to be blocking squares because of it
			// TODO: actually subtly weight the path rather than just a binary good/bad
			if distanceToOrigin > 5 {
				continue
			}

			for _, snake := range s.Snakes {
				// use yourhead here to check if this is us
				if generator.SamePoint(snake.Body[0], yourHead) {
					continue
				}
				distanceToOtherSnake := Abs(snake.Body[0].X-x) + Abs(snake.Body[0].Y-y)
				if distanceToOtherSnake < distanceToOrigin {
					p[x][y].Blocked = true

				}
			}

		}
	}

	for _, hazard := range s.Hazards {
		p[hazard.X][hazard.Y].Blocked = true
	}
}

func (p PathGrid) AddFogOfUncertainty(head rules.Point, distance int32) {

	neighbours := generator.NeighboursSafe(len(p[0]), len(p), head)
	for _, neighbour := range neighbours {
		p[neighbour.X][neighbour.Y].Blocked = true
		p.AddFogOfUncertainty(neighbour, distance-1)
	}

}

func (p PathGrid) GetNextLowestSquare() (rules.Point, error) {

	lowestCost := LargestCost
	var lowestX, lowestY int32
	found := false
	for x, xRow := range p {
		for y := range xRow {
			if p[x][y].Explored || p[x][y].Blocked || p[x][y].DistanceToTarget == 0 {
				continue
			}
			cost := p[x][y].DistanceToOrigin + p[x][y].DistanceToTarget
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

func (p PathGrid) NextStepTowardsTarget(current, target rules.Point) (rules.Point, error) {

	var nextCoord rules.Point
	// lowestCost := LargestCost
	// lowestDistanceToTarget := p[current.X][current.Y].DistanceToTarget
	lowestDistanceToTarget := LargestCost
	found := false

	neighbours := generator.NeighboursSafe(len(p[0]), len(p), current)
	for _, neighbour := range neighbours {

		// if thi
		if generator.SamePoint(neighbour, target) {
			return neighbour, nil
		}

		// if this neighbour is blocked or there is no route from it to the target, don't check its distance
		if p[neighbour.X][neighbour.Y].Blocked || p[neighbour.X][neighbour.Y].DistanceToTarget == 0 {
			continue
		}

		// you only need to check the distance to target, not the cost. cost is used in calculating the next
		// square to check in the distance to target sum.
		if p[neighbour.X][neighbour.Y].DistanceToTarget < lowestDistanceToTarget {
			lowestDistanceToTarget = p[neighbour.X][neighbour.Y].DistanceToTarget
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

func GetRoutesFromOrigin(state *rules.BoardState, origin, yourHead rules.Point) [][]rules.Point {
	p := MakePathgrid(state, origin, yourHead)
	var routes [][]rules.Point

	for x := int32(0); x < state.Width; x++ {
		for y := int32(0); y < state.Height; y++ {
			target := rules.Point{X: x, Y: y}
			if p[x][y].Blocked {
				continue
			}

			route, err := p.GetRoute(origin, target)
			if err != nil {
				continue
			}

			routes = append(routes, route)

		}
	}

	return routes
}

func (p PathGrid) FreeSquares(state *rules.BoardState) int {
	var freeSquares int
	for x := int32(0); x < state.Width; x++ {
		for y := int32(0); y < state.Height; y++ {
			if !p[x][y].Blocked {
				freeSquares++
			}
		}
	}
	return freeSquares
}

// TraceWeightedGridToTarget takes a weighted grid and finds the best next step on the route from origin to target
func (p PathGrid) TraceWeightedGridToTarget(origin, target rules.Point) ([]rules.Point, error) {

	route := []rules.Point{}
	currentPoint := origin

	if !p[origin.X][origin.Y].Blocked && p[origin.X][origin.Y].DistanceToTarget == 0 {
		return nil, fmt.Errorf("no path to target origin: %+v target %+v", origin, target)
	}

	for {
		nextPoint, err := p.NextStepTowardsTarget(currentPoint, target)
		if err != nil {
			return nil, err
		}

		route = append(route, nextPoint)

		if len(route) > 50 {
			log.WithFields(log.Fields{
				"pathgrid":   p,
				"origin":     origin,
				"target":     target,
				"next_point": nextPoint,
			}).Error("finding path entered a loop. indicates an issue with a* algorithm")
			return route, fmt.Errorf("got too many points in route.")
		}

		if generator.SamePoint(nextPoint, target) {
			break
		}

		currentPoint = nextPoint

	}

	return route, nil

}

// CalculateDistancesToOrigin performs the check of pythagorean distance from the origin to all other points on the board.
// Since this can be reused it reduces iterations to perform this once on all squares since we will be checking
// routes to multiple points from a single origin.
func (p PathGrid) CalculateDistancesToOrigin(origin rules.Point) {

	for x := 0; x < len(p); x++ {
		xDelta := origin.X - int32(x)
		for y := 0; y < len(p[x]); y++ {
			yDelta := origin.Y - int32(y)
			p[x][y].DistanceToOrigin = int(math.Sqrt(float64(xDelta*xDelta)+float64(yDelta*yDelta)) * CostFactor)
		}
	}
}

// CalculateDistancesToOrigin returns the shortest path from origin to target
func (p PathGrid) CalculateDistancesToTarget(origin, target rules.Point) {

	var err error
	nextSquare := target

	for {
		p[nextSquare.X][nextSquare.Y].Explored = true
		p.ScoreNeighbours(nextSquare)

		nextSquare, err = p.GetNextLowestSquare()
		if err != nil {
			return
		}

	}

}

func (p PathGrid) ResetDistancesToTarget() {
	for _, xColumn := range p {
		for _, y := range xColumn {
			y.DistanceToTarget = 0
			y.Explored = false
		}
	}
}

func (p PathGrid) GetRoute(origin, target rules.Point) ([]rules.Point, error) {
	p.CalculateDistancesToTarget(origin, target)
	defer p.ResetDistancesToTarget()
	return p.TraceWeightedGridToTarget(origin, target)
}
