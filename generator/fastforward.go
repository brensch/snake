package generator

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

func FastForward(s *rules.BoardState, ruleset rules.Ruleset, snake rules.Snake, proposedRoute []rules.Point) *rules.BoardState {
	// to find the routes from this point on, fastforward your position to what it would be
	// given this route, and make everyone else do their safest move
	// TODO: actually try and assume that they make good moves here
	nextState := s.Clone()

	previousHead := snake.Body[0]

	// need to traverse route backwards
	for routePosition := range proposedRoute {
		pointInRoute := proposedRoute[len(proposedRoute)-routePosition-1]

		moves := []rules.SnakeMove{
			{ID: snake.ID, Move: DirectionToPoint(previousHead, pointInRoute)},
		}
		previousHead = pointInRoute

		for _, snakeIter := range nextState.Snakes {
			if snake.ID == snakeIter.ID {
				continue
			}

			safestMoves := SafestMoves(nextState, ruleset, snakeIter)
			if len(safestMoves) == 0 {
				moves = append(moves, rules.SnakeMove{ID: snakeIter.ID, Move: rules.DirectionDown})
				continue
			}
			moves = append(moves, rules.SnakeMove{ID: snakeIter.ID, Move: safestMoves[0]})

		}

		var err error
		nextState, err = ruleset.CreateNextBoardState(nextState, moves)
		if err != nil {
			fmt.Println("very strange", err)
		}

	}

	return nextState
}
