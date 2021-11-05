package generator_test

import (
	"fmt"
	"testing"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func TestAllMovesForState(t *testing.T) {

	state, err := rules.CreateDefaultBoardState(11, 11, []string{
		"1",
		"2",
		"3",
	})
	if err != nil {
		t.Fail()
		return
	}

	state.Food = nil
	state.Snakes = []rules.Snake{
		{
			ID:              "1",
			Body:            []rules.Point{{10, 10}, {10, 10}, {10, 10}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "2",
			Body:            []rules.Point{{7, 7}, {7, 7}, {7, 7}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "3",
			Body:            []rules.Point{{2, 2}, {2, 2}, {2, 2}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "3",
			Body:            []rules.Point{{4, 4}, {4, 4}, {4, 4}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
	}

	points := generator.AllMoveSetsForState(state)
	for _, row := range points {
		fmt.Println(row)
	}
	fmt.Println(len(points))

}
