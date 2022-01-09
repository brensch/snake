package minimax

import (
	"encoding/binary"
	"fmt"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func (n *Node) ExploreBestPath() {

	move := 0
	nextChild := n
	odd := false
	for {
		temp := nextChild.FindBestChild()
		if temp == nil {
			break
		}
		nextChild = temp
		odd = !odd
		if odd {
			continue
		}
		fmt.Println("-----------move", move)
		move++
		nextChild.Print()
		generator.PrintMap(nextChild.State)
		ShortestPathsBreadthPrint(nextChild.State)
	}
}

func Hash(board *rules.BoardState) uint64 {

	// maximum possible size for two snakes
	// 121*2 (for x and y of each location)
	// 2 delimiters
	// 2 health
	// 2 more delimiters
	// at this length there can be no food, so don't take into account (snakes will have consumed)
	// this totals 248 so 256 is laughing
	buffer := make([]byte, 256)

	// why not
	delimiter := byte(69)

	offset := 0

	for _, snake := range board.Snakes {
		buffer[offset] = snake.Health
		buffer[offset+1] = delimiter
		offset += 2
		for _, bodyPiece := range snake.Body {
			buffer[offset] = bodyPiece.X
			buffer[offset+1] = bodyPiece.Y
			offset += 2
		}
		buffer[offset] = delimiter
		offset++

	}

	// todo: decide if we should include food
	// for _, food := range board.Food {
	// }

	return binary.BigEndian.Uint64(buffer)
}
