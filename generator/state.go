package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

func GetYou(state *rules.BoardState, youID string) (rules.Snake, error) {
	for _, snake := range state.Snakes {
		if snake.ID == youID {
			return snake, nil
		}
	}

	return rules.Snake{}, fmt.Errorf("no snake found with id %s", youID)
}
