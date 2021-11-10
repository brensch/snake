package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

func SafetyDance(state *rules.BoardState, ruleset rules.Ruleset, youID string) [4][][]rules.SnakeMove {

	moveSets := AllMoveSetsForStateRaw(state)
	var nextMoveSets [4][][]rules.SnakeMove
	var youPosition int
	// var you rules.Snake
	foundSnake := false

	// get your position
	for position, snake := range state.Snakes {
		if youID != snake.ID {
			foundSnake = true
			continue
		}

		youPosition = position
		// you = snake
		break
	}

	if !foundSnake {
		fmt.Println("no you bruh")
		return nextMoveSets
	}

	// go through all moveSets, generate them, and see which ones we don't die in
	for _, moveSet := range moveSets {

		nextState, err := ruleset.CreateNextBoardState(state, moveSet)
		if err != nil {
			// log.WithFields(log.Fields{
			// 	"state": state,
			// }).Error("failed to create next board state")
			continue
		}

		sadSnake := false
		for _, snake := range nextState.Snakes {
			if snake.EliminatedCause != "" {
				sadSnake = true
				// fmt.Println(snake.EliminatedCause)
				continue
			}
		}

		if sadSnake {
			continue
		}

		if nextState.Snakes[youPosition].EliminatedCause != "" {
			continue
		}

		yourMove := moveSet[youPosition]

		nextMoveSets[yourMove.Move] = append(nextMoveSets[yourMove.Move], moveSet)

	}

	return nextMoveSets
}

func SafeMoves(state *rules.BoardState, ruleset rules.Ruleset, youID string) [4]int {

	moves := AllMoveSetsForState(state)
	var safeMoves [4]int
	var youPosition int
	var you rules.Snake

	// get your position
	for position, snake := range state.Snakes {
		if youID != snake.ID {
			continue
		}

		youPosition = position
		you = snake
		break
	}

	// go through all moves, generate them, and see which ones we don't die in
	for _, move := range moves {

		var snakeMoves []rules.SnakeMove
		for moveSnakeNum, movePoint := range move {
			snakeMoves = append(snakeMoves, rules.SnakeMove{
				ID:   state.Snakes[moveSnakeNum].ID,
				Move: DirectionToPoint(state.Snakes[moveSnakeNum].Body[0], movePoint),
			})
		}

		nextState, err := ruleset.CreateNextBoardState(state, snakeMoves)
		if err != nil {
			// log.WithFields(log.Fields{
			// 	"state": state,
			// }).Error("failed to create next board state")
			continue
		}

		if nextState.Snakes[youPosition].EliminatedCause != "" {
			continue
		}

		direction := DirectionToPoint(you.Body[0], move[youPosition])
		safeMoves[direction]++

	}

	return safeMoves

}

func SafestMoves(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake) []rules.Direction {

	safeMoves := SafeMoves(state, ruleset, you.ID)

	var safestMoveCount int
	for _, move := range safeMoves {
		if move > safestMoveCount {
			safestMoveCount = move
		}
	}

	safestMoves := []rules.Direction{}
	for direction, move := range safeMoves {
		if move == safestMoveCount {
			safestMoves = append(safestMoves, rules.Direction(direction))
		}
	}

	return safestMoves
}

func MoveIsSafe(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, candidate rules.Direction) bool {
	safeMoves := SafeMoves(state, ruleset, you.ID)
	if safeMoves[candidate] > 0 {
		return true
	}

	return false
}
