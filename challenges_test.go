package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func TestAvoidWalls(t *testing.T) {
	snakeID := "brend"

	state, err := rules.CreateDefaultBoardState(5, 5, []string{snakeID})
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println("snek", state.Snakes)
	state.Snakes = []rules.Snake{
		{
			ID:              snakeID,
			Body:            []rules.Point{{2, 2}, {2, 2}, {2, 2}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
	}
	state.Food = []rules.Point{}

	standardGame := &rules.StandardRuleset{
		FoodSpawnChance: 0,
		MinimumFood:     0,
	}

	// generator.PrintMap(state)

	turn := int32(0)
	target := int32(100)

	for {
		t.Log("turn", turn)

		you, err := GetSnakeFromState(state, snakeID)
		if err != nil {
			t.Fail()
			return
		}

		if you.EliminatedCause != "" {
			fmt.Println("eliminated by", you.EliminatedCause)
			if turn < target {
				t.Fail()
			}
			return
		}

		move, _ := Move(context.Background(), state, standardGame, you, turn)

		nextMove := rules.SnakeMove{
			ID:   snakeID,
			Move: move.String(),
		}

		nextState, err := standardGame.CreateNextBoardState(state, []rules.SnakeMove{nextMove})
		if err != nil {
			t.Fail()
			return
		}

		generator.PrintMap(state)

		turn++
		state = nextState
		state.Turn = int32(turn)
	}
}

func TestAvoidBattleSnakes(t *testing.T) {

	snakeID := "brend"

	state, err := rules.CreateDefaultBoardState(11, 11, []string{
		snakeID,
		"bottomLeft",
		"topLeft",
		"topRight",
		"bottomRight",
	})
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println("snek", state.Snakes)
	state.Food = nil
	state.Snakes = []rules.Snake{
		{
			ID:              snakeID,
			Body:            []rules.Point{{5, 5}, {5, 5}, {5, 5}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "bottomLeft",
			Body:            []rules.Point{{3, 4}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "topLeft",
			Body:            []rules.Point{{1, 9}, {1, 8}, {1, 7}, {1, 6}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {5, 7}, {5, 8}, {5, 9}, {4, 9}, {3, 9}, {2, 9}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "topRight",
			Body:            []rules.Point{{10, 8}, {10, 7}, {10, 6}, {9, 6}, {8, 6}, {7, 6}, {7, 7}, {7, 8}, {7, 9}, {7, 10}, {8, 10}, {9, 10}, {10, 10}, {10, 9}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
		{
			ID:              "bottomRight",
			Body:            []rules.Point{{9, 3}, {9, 4}, {8, 4}, {7, 4}, {6, 4}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {9, 2}},
			Health:          100,
			EliminatedCause: "",
			EliminatedBy:    "",
		},
	}

	target := int32(100)
	turn := int32(0)
	standardGame := &rules.StandardRuleset{
		FoodSpawnChance: 0,
		MinimumFood:     0,
	}

	for {

		you, err := GetSnakeFromState(state, snakeID)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if you.EliminatedCause != "" {
			fmt.Println("eliminated by", you.EliminatedCause)
			if turn < target {
				t.Log("eliminated too early")
				t.Fail()
			}
			return
		}

		move, _ := Move(context.Background(), state, standardGame, you, turn)

		nextMove := rules.SnakeMove{
			ID:   snakeID,
			Move: move.String(),
		}

		var botMoves []rules.SnakeMove
		for _, bot := range state.Snakes[1:len(state.Snakes)] {

			botMoves = append(botMoves, rules.SnakeMove{
				ID:   bot.ID,
				Move: generator.ChaserNextMove(bot).String(),
			})
		}

		allMoves := append(botMoves, nextMove)
		fmt.Println(allMoves)

		nextState, err := standardGame.CreateNextBoardState(
			state,
			append(botMoves, nextMove),
		)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		generator.PrintMap(state)

		turn++
		state = nextState
		state.Turn = int32(turn)

		if turn > 20 {
			return
		}

		gameOver, err := standardGame.IsGameOver(state)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}
		if gameOver {
			t.Log("gameover?")
			t.Fail()
			return
		}
	}

}
