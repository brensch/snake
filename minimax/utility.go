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

func QuickHash(board *rules.BoardState) {

	buffer := make([]byte, 1000000)

	offset := 0
	// var thi
	// var thing rules.Point

	for _, snake := range board.Snakes {
		for _, bodyPiece := range snake.Body {

			buffer[offset] = bodyPiece.X
			buffer[offset+1] = bodyPiece.Y
			offset += 2
		}
	}
}
