package minimax

import (
	"bytes"
	"encoding/gob"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func (n *Node) ExploreBestPath() {

	nextChild := n
	for {
		// fmt.Println("got score of", *nextChild.Score)
		temp := nextChild.FindBestChild()
		if temp == nil {
			break
		}
		nextChild = temp
		nextChild.Print()
		generator.PrintMap(nextChild.State)
		ShortestPathsBreadthPrint(nextChild.State)
	}
}

// func Hash(state *rules.BoardState) string {

// 	// todo: figure out if food is necessary
// 	// snakeString := fmt.Sprintf("%v-%d-%v-%d",
// 	// 	state.Snakes[0].Body,
// 	// 	state.Snakes[0].Health,
// 	// 	state.Snakes[1].Body,
// 	// 	state.Snakes[1].Health,
// 	// )

// 	snakeString := fmt.Sprint(state)
// 	return snakeString
// }

func Hash(board *rules.BoardState) []byte {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(board)
	return b.Bytes()
}

func QuickHash(board *rules.BoardState) []byte {

	// maximum possible size for two snakes
	// 121*2 (for x and y of each location)
	// 2 delimiters
	// 2 health
	// 2 more delimiters
	// at this length there can be no food, so don't take into account (snakes will have consumed)
	buffer := make([]byte, 248)

	// why not
	delimiter := byte(69)

	offset := 0
	// var thi
	// var thing rules.Point

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

	return buffer
}
