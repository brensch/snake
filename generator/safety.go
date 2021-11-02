package generator

import (
	"github.com/brensch/snake/rules"
	log "github.com/sirupsen/logrus"
)

func SafeMoves(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake) [4]int {

	moves := AllMovesForState(state)
	var safeMoves [4]int
	var youPosition int

	// get your position
	for position, snake := range state.Snakes {
		if you.ID != snake.ID {
			continue
		}

		youPosition = position
		break
	}

	// go through all moves, generate them, and see which ones we don't die in
	for _, move := range moves {

		var snakeMoves []rules.SnakeMove
		for moveSnakeNum, movePoint := range move {
			snakeMoves = append(snakeMoves, rules.SnakeMove{
				ID:   state.Snakes[moveSnakeNum].ID,
				Move: DirectionToPoint(state.Snakes[moveSnakeNum].Body[0], movePoint).String(),
			})
		}

		nextState, err := ruleset.CreateNextBoardState(state, snakeMoves)
		if err != nil {
			log.WithFields(log.Fields{
				"state": state,
			}).Error("failed to create next board state")
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

func SafestMoves(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake) []Direction {

	safeMoves := SafeMoves(state, ruleset, you)

	var safestMoveCount int
	for _, move := range safeMoves {
		if move > safestMoveCount {
			safestMoveCount = move
		}
	}

	safestMoves := []Direction{}
	for direction, move := range safeMoves {
		if move == safestMoveCount {
			safestMoves = append(safestMoves, Direction(direction))
		}
	}

	return safestMoves
}

func MoveIsSafe(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, candidate Direction) bool {
	safeMoves := SafeMoves(state, ruleset, you)
	if safeMoves[candidate] > 0 {
		return true
	}

	return false
}
