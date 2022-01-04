package generator

import (
	"bytes"
	"fmt"
	"log"

	"github.com/brensch/snake/rules"
)

var (
	characters = [8]rune{'☻', '■', '⌀', '●', '☺', '□', '⍟', '◘'}
)

// PrintMap prints a visual representation of the current map
func PrintMap(state *rules.BoardState) {
	var o bytes.Buffer
	board := make([][]rune, state.Width)
	for i := range board {
		board[i] = make([]rune, state.Height)
	}
	for y := byte(0); y < state.Height; y++ {
		for x := byte(0); x < state.Width; x++ {
			board[x][y] = '◦'
		}
	}
	for _, oob := range state.Hazards {
		board[oob.X][oob.Y] = '░'
	}
	o.WriteString(fmt.Sprintf("Hazards ░: %v\n", state.Hazards))
	for _, f := range state.Food {
		board[f.X][f.Y] = '⚕'
	}
	o.WriteString(fmt.Sprintf("Food ⚕: %v\n", state.Food))
	for numSnake, s := range state.Snakes {
		for _, b := range s.Body {
			if b.X >= 0 && b.X < state.Width && b.Y >= 0 && b.Y < state.Height {
				board[b.X][b.Y] = characters[numSnake]
			}
		}

		o.WriteString(fmt.Sprintf("%v %c: %v\n", state.Snakes[numSnake].ID, characters[numSnake], s))
	}
	for y := int32(state.Height) - 1; y >= 0; y-- {
		for x := int32(0); x < int32(state.Width); x++ {
			o.WriteRune(board[x][y])
		}
		o.WriteString("\n")
	}
	log.Print(o.String())
}
