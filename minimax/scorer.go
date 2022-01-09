package minimax

import (
	"encoding/json"
	"fmt"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

// type

func PercentageOfBoardControlled(board *rules.BoardState) float64 {

	// you = 0

	// calculate how many squares i can reach first
	allMovesFromSquares := make([][]byte, len(board.Snakes))
	totalSpaces := int(board.Height * board.Width)

	// for each snake, find a naive distance to each point
	for snakeNumber, snake := range board.Snakes {

		movesFromSquares := make([]byte, totalSpaces)
		boardSpace := 0
		for x := byte(0); x < board.Width; x++ {
			for y := byte(0); y < board.Height; y++ {
				movesFromSquares[boardSpace] = generator.Abs(snake.Body[0].X-x) + generator.Abs(snake.Body[0].Y-y)
				boardSpace++
			}
		}

		allMovesFromSquares[snakeNumber] = movesFromSquares
	}

	closestSquareCount := 0

	for i := 0; i < totalSpaces; i++ {
		closestDistance := byte(255)
		closestSnake := -1
		for snake := range board.Snakes {
			if allMovesFromSquares[snake][i] < byte(closestDistance) {
				closestDistance = allMovesFromSquares[snake][i]
				closestSnake = snake
			}
		}

		// check if index of closest snake is maxer index
		if closestSnake == 0 {
			closestSquareCount++
		}

	}

	return (float64(closestSquareCount) / float64(totalSpaces))

}

// GameFinished returns
// +1: maxer won
// -1: maxer lost
// 0: not finished
func GameFinished(board *rules.BoardState, isMax bool) float64 {

	// the cli starts at 1 but the online servers start at 0
	if board.Turn == 1 || board.Turn == 0 {
		return 0
	}

	maxSnake := board.Snakes[0]
	minSnake := board.Snakes[1]
	maxHead := maxSnake.Body[0]
	minHead := minSnake.Body[0]

	if maxHead.X == minHead.X && maxHead.Y == minHead.Y {
		// maxer can only win if the head collision happens on the min turn.
		// game is always calculated with maxer going first and taking turns, when in reality
		// the turns are made simultaneously.
		// we also want to avoid draws, so making draw state a loss
		if isMax || len(maxSnake.Body) <= len(minSnake.Body) {
			return -1
		}

		return 1
	}

	if maxSnake.Health == 0 {
		return -1
	}

	if minSnake.Health == 0 {
		return 1
	}

	for _, maxPiece := range maxSnake.Body[1:] {
		if maxPiece.X == maxHead.X && maxPiece.Y == maxHead.Y {
			return -1
		}
		if maxPiece.X == minHead.X && maxPiece.Y == minHead.Y {
			return 1
		}
	}

	// also check the head, but only for max vs min
	if maxSnake.Body[0].X == minHead.X && maxSnake.Body[0].Y == minHead.Y {
		return 1
	}

	for _, minPiece := range minSnake.Body[1:] {
		if minPiece.X == minHead.X && minPiece.Y == minHead.Y {
			return 1
		}
		if minPiece.X == maxHead.X && minPiece.Y == maxHead.Y {
			return -1
		}
	}

	// also check the head, but only for min vs max (improving performance)
	if minSnake.Body[0].X == maxHead.X && minSnake.Body[0].Y == maxHead.Y {
		return -1
	}

	return 0

}

func GameFinishedBits(snake1, snake2 int) float64 {
	if snake1^snake2 != 0 {
		return 1

	}

	return 0
}

func HeuristicAnalysis(board *rules.BoardState) float64 {

	// healthScore := 1.0
	// if board.Snakes[0].Health < 20 {
	// 	healthScore = 0.5
	// }

	// percentLengthOfOtherSnake := float64(len(board.Snakes[0].Body)) / float64(len(board.Snakes[1].Body))
	// lengthScore := percentLengthOfOtherSnake
	// if lengthScore > 1.1 {
	// 	lengthScore = 1.1
	// }

	areacontrol := ShortestPathsBreadth(board)

	// return areacontrol * lengthScore * healthScore

	// return PercentageOfBoardControlled(board) * lengthScore
	// return ShortestPathsBreadth(board) * lengthScore * healthScore
	// return PercentageOfBoardControlled(board)
	return areacontrol
}

func ShortestPaths(board *rules.BoardState) {

	// var obstacleGrid [11][11]int

	obstacleGrid := make([]int, 11*11)

	// snakeRoutes

	for _, snake := range board.Snakes {
		for snakePieceIndex, snakePiece := range snake.Body {
			obstacleGrid[snakePiece.Y*11+snakePiece.X] = len(snake.Body) - snakePieceIndex
		}
	}

	// PrintShortestPath(obstacleGrid)

	// iterate through each snake and do a dfs
	for snakeCount, snake := range board.Snakes {
		snakeRoute := make([]int, 11*11)

		// start at head
		head := snake.Body[0]

		snakeRoute[head.Y*11+head.X] = 1
		_ = snakeCount
		ExplorePoint(snakeRoute, obstacleGrid, int(head.X), int(head.Y))
		// ExplorePoint(snakeRoute, int(head.X+1), int(head.Y))
		// ExplorePoint(snakeRoute, int(head.X), int(head.Y+1))

		// fmt.Println("snake ", snakeCount)
		// PrintShortestPath(snakeRoute)
		// 	if

	}

	// return obstacleGrid

}

func ExplorePoint(graph, obstacles []int, x, y int) {
	// fmt.Println(x, y)
	// PrintShortestPath(graph)
	originScore := graph[y*11+x]
	// _ = originScore
	var directions [4]bool
	if x > 0 && (graph[y*11+(x-1)] > originScore+1 || graph[y*11+(x-1)] == 0) && obstacles[y*11+(x-1)] == 0 {
		graph[y*11+(x-1)] = originScore + 1
		directions[0] = true
		// ExplorePoint(graph, obstacles, x-1, y)
	}

	if x < 10 && (graph[y*11+(x+1)] > originScore+1 || graph[y*11+(x+1)] == 0) && obstacles[y*11+(x+1)] == 0 {
		graph[y*11+(x+1)] = originScore + 1
		directions[1] = true

		// ExplorePoint(graph, obstacles, x+1, y)
	}

	if y > 0 && (graph[(y-1)*11+x] > originScore+1 || graph[(y-1)*11+x] == 0) && obstacles[(y-1)*11+x] == 0 {
		graph[(y-1)*11+x] = originScore + 1
		directions[2] = true

		// ExplorePoint(graph, obstacles, x, y-1)
	}

	if y < 10 && (graph[(y+1)*11+x] > originScore+1 || graph[(y+1)*11+x] == 0) && obstacles[(y+1)*11+x] == 0 {
		graph[(y+1)*11+x] = originScore + 1
		directions[3] = true

		// ExplorePoint(graph, obstacles, x, y+1)
	}

	if directions[0] {
		ExplorePoint(graph, obstacles, x-1, y)
	}
	if directions[1] {
		ExplorePoint(graph, obstacles, x+1, y)
	}
	if directions[2] {
		ExplorePoint(graph, obstacles, x, y-1)
	}
	if directions[3] {
		ExplorePoint(graph, obstacles, x, y+1)
	}
}

func ShortestPaths2(board *rules.BoardState) {

	// var obstacleGrid [11][11]int

	obstacleGrid := make([]int, 11*11)

	// snakeRoutes

	for _, snake := range board.Snakes {
		for snakePieceIndex, snakePiece := range snake.Body {
			obstacleGrid[snakePiece.Y*11+snakePiece.X] = len(snake.Body) - snakePieceIndex
		}
	}

	// PrintShortestPath(obstacleGrid)

	// iterate through each snake and do a dfs
	for snakeCount, snake := range board.Snakes {
		snakeRoute := make([]int, 11*11)

		// start at head
		head := snake.Body[0]

		snakeRoute[head.Y*11+head.X] = 1
		_ = snakeCount
		ExplorePoint2(snakeRoute, obstacleGrid, int(head.X), int(head.Y))
		// ExplorePoint(snakeRoute, int(head.X+1), int(head.Y))
		// ExplorePoint(snakeRoute, int(head.X), int(head.Y+1))

		// fmt.Println("snake ", snakeCount)
		// PrintShortestPath(snakeRoute)
		// 	if

	}

	// return obstacleGrid

}

func ExplorePoint2(graph, obstacles []int, x, y int) {
	// fmt.Println(x, y)
	// PrintShortestPath(graph)
	originScore := graph[y*11+x]
	// _ = originScore
	if x > 0 && (graph[y*11+(x-1)] > originScore+1 || graph[y*11+(x-1)] == 0) && obstacles[y*11+(x-1)] == 0 {
		graph[y*11+(x-1)] = originScore + 1
		ExplorePoint(graph, obstacles, x-1, y)
	}

	if x < 10 && (graph[y*11+(x+1)] > originScore+1 || graph[y*11+(x+1)] == 0) && obstacles[y*11+(x+1)] == 0 {
		graph[y*11+(x+1)] = originScore + 1
		ExplorePoint(graph, obstacles, x+1, y)
	}

	if y > 0 && (graph[(y-1)*11+x] > originScore+1 || graph[(y-1)*11+x] == 0) && obstacles[(y-1)*11+x] == 0 {
		graph[(y-1)*11+x] = originScore + 1
		ExplorePoint(graph, obstacles, x, y-1)
	}

	if y < 10 && (graph[(y+1)*11+x] > originScore+1 || graph[(y+1)*11+x] == 0) && obstacles[(y+1)*11+x] == 0 {
		graph[(y+1)*11+x] = originScore + 1
		ExplorePoint(graph, obstacles, x, y+1)
	}

}

func ShortestPathsBreadth(board *rules.BoardState) float64 {

	// var obstacleGrid [11][11]int

	obstacleGrid := make([]int, 11*11)

	// snakeRoutes

	for _, snake := range board.Snakes {
		for snakePieceIndex, snakePiece := range snake.Body {
			obstacleGrid[snakePiece.Y*11+snakePiece.X] = len(snake.Body) - snakePieceIndex
		}
	}

	// PrintShortestPath(obstacleGrid)
	snakeRoutes := make([][]int, 2)

	// iterate through each snake and do a dfs
	for snakeCount, snake := range board.Snakes {
		snakeRoutes[snakeCount] = make([]int, 11*11)

		// start at head
		head := snake.Body[0]

		snakeRoutes[snakeCount][head.Y*11+head.X] = 1
		_ = snakeCount
		ExplorePoints(snakeRoutes[snakeCount], obstacleGrid, [][2]int{{int(head.X), int(head.Y)}}, 1)
		// ExplorePoint(snakeRoute, int(head.X+1), int(head.Y))
		// ExplorePoint(snakeRoute, int(head.X), int(head.Y+1))

		// fmt.Println("snake ", snakeCount)
		// PrintShortestPath(snakeRoutes[snakeCount])

	}

	controlledSquares := 0
	reachableSquares := 0

	for x := 0; x < 11; x++ {
		for y := 0; y < 11; y++ {

			if snakeRoutes[0][y*11+x] == 0 && snakeRoutes[1][y*11+x] > 0 {
				// matrix[y*11+x] = 2
				reachableSquares++
				continue
			}

			if snakeRoutes[1][y*11+x] == 0 && snakeRoutes[0][y*11+x] > 0 {
				// matrix[y*11+x] = 1
				reachableSquares++
				controlledSquares++
				continue
			}

			if snakeRoutes[0][y*11+x] < snakeRoutes[1][y*11+x] {
				reachableSquares++
				controlledSquares++
				// matrix[y*11+x] = 1
				continue
			}
		}
	}

	// if 1-(float64(controlledSquares)/float64(121)+0.5) == 0.12809917355371903 {
	// 	fmt.Println(float64(controlledSquares), float64(121))
	// }

	if controlledSquares == 0 {
		return 0.0000001
	}
	return float64(controlledSquares) / float64(reachableSquares)

	// return obstacleGrid

}

func ShortestPathsBreadthPrint(board *rules.BoardState) float64 {

	// var obstacleGrid [11][11]int

	obstacleGrid := make([]int, 11*11)

	// snakeRoutes

	for _, snake := range board.Snakes {
		for snakePieceIndex, snakePiece := range snake.Body {
			obstacleGrid[snakePiece.Y*11+snakePiece.X] = len(snake.Body) - snakePieceIndex
		}
	}

	// PrintShortestPath(obstacleGrid)
	snakeRoutes := make([][]int, 2)

	// iterate through each snake and do a dfs
	for snakeCount, snake := range board.Snakes {
		snakeRoutes[snakeCount] = make([]int, 11*11)

		// start at head
		head := snake.Body[0]

		snakeRoutes[snakeCount][head.Y*11+head.X] = 1
		_ = snakeCount
		ExplorePoints(snakeRoutes[snakeCount], obstacleGrid, [][2]int{{int(head.X), int(head.Y)}}, 1)
		// ExplorePoint(snakeRoute, int(head.X+1), int(head.Y))
		// ExplorePoint(snakeRoute, int(head.X), int(head.Y+1))

		fmt.Println("snake ", snakeCount)
		PrintShortestPath(snakeRoutes[snakeCount])

	}

	controlledSquares := 0
	reachableSquares := 0
	matrix := make([]int, 121)

	for x := 0; x < 11; x++ {
		for y := 0; y < 11; y++ {

			if snakeRoutes[0][y*11+x] == 0 && snakeRoutes[1][y*11+x] > 0 {
				matrix[y*11+x] = 2
				reachableSquares++
				continue
			}

			if snakeRoutes[1][y*11+x] == 0 && snakeRoutes[0][y*11+x] > 0 {
				matrix[y*11+x] = 1
				reachableSquares++
				controlledSquares++
				continue
			}

			if snakeRoutes[0][y*11+x] < snakeRoutes[1][y*11+x] {
				reachableSquares++
				controlledSquares++
				matrix[y*11+x] = 1
				continue
			}
		}
	}

	fmt.Println("controlled area by us", float64(controlledSquares)/float64(reachableSquares))
	PrintShortestPath(matrix)
	boardbytes, _ := json.Marshal(board)
	fmt.Println(string(boardbytes))

	if controlledSquares == 0 {
		fmt.Println("got no square")
		return 0.0000001
	}
	return float64(controlledSquares) / float64(reachableSquares)

	// return obstacleGrid

}

func ExplorePoints(graph, obstacles []int, coords [][2]int, prevScore int) {

	// the size selected it the largest possible size the next iteration can be.
	// doing this doubles performance (!). appends are very slow. speed is critical in this function even if things get unidiomatic.
	nextPoints := make([][2]int, len(coords)*3+1)
	// use this to know where we are in the nextPoints array
	pointCount := 0

	for _, coord := range coords {
		x := coord[0]
		y := coord[1]

		// _ = originScore
		if x > 0 && (graph[y*11+(x-1)] > prevScore+1 || graph[y*11+(x-1)] == 0) && obstacles[y*11+(x-1)] < prevScore+1 {
			graph[y*11+(x-1)] = prevScore + 1
			nextPoints[pointCount] = [2]int{x - 1, y}
			pointCount++
		}

		if x < 10 && (graph[y*11+(x+1)] > prevScore+1 || graph[y*11+(x+1)] == 0) && obstacles[y*11+(x+1)] < prevScore+1 {
			graph[y*11+(x+1)] = prevScore + 1
			nextPoints[pointCount] = [2]int{x + 1, y}
			pointCount++
		}

		if y > 0 && (graph[(y-1)*11+x] > prevScore+1 || graph[(y-1)*11+x] == 0) && obstacles[(y-1)*11+x] < prevScore+1 {
			graph[(y-1)*11+x] = prevScore + 1
			nextPoints[pointCount] = [2]int{x, y - 1}
			pointCount++

		}

		if y < 10 && (graph[(y+1)*11+x] > prevScore+1 || graph[(y+1)*11+x] == 0) && obstacles[(y+1)*11+x] < prevScore+1 {
			graph[(y+1)*11+x] = prevScore + 1
			nextPoints[pointCount] = [2]int{x, y + 1}
			pointCount++

		}
	}

	if pointCount > 0 {
		ExplorePoints(graph, obstacles, nextPoints[:pointCount], prevScore+1)
	}

}

func PrintShortestPath(graph []int) {

	for y := 10; y >= 0; y-- {
		for x := 0; x < 11; x++ {
			fmt.Printf("%3d", graph[y*11+x])
		}
		fmt.Println()
	}

}
