package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/pather"
	"github.com/brensch/snake/rules"
)

const LargestCost = 10000

func Move(ctx context.Context, state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, turn int32) (generator.Direction, string) {

	var tastiestSnackPath []rules.Point
	foundSnack := false

	p := pather.MakePathgrid(state, you.Body[0], you.Body[0])
	freeSquares := p.FreeSquares(state)

	for _, snack := range state.Food {

		route, err := p.GetRoute(you.Body[0], snack)
		if err != nil {
			continue
		}

		routesFromSnackOnwards := pather.GetRoutesFromOrigin(state, snack, you.Body[0])

		if len(routesFromSnackOnwards) < freeSquares/2 {
			// if len(routesFromSnackOnwards) < len(you.Body) {
			log.WithFields(log.Fields{
				"reachable": len(routesFromSnackOnwards),
				"total":     freeSquares,
			}).Debug("not enough room to fit ya boi if i chase that snack")
			continue
		}

		if (len(tastiestSnackPath) == 0 && len(route) > 0) ||
			len(route) < len(tastiestSnackPath) {
			tastiestSnackPath = route
			foundSnack = true
		}
	}

	// needMoreGirth := false
	// for _, snake := range state.Snakes {
	// 	if snake.ID == you.ID {
	// 		continue
	// 	}
	// 	if len(snake.Body) >= len(you.Body)-1 {
	// 		// fmt.Println(you.ID, "need more girth")
	// 		needMoreGirth = true
	// 	}
	// }

	// just fed is important because our tail stays in the same place for an extra move
	// justFed := generator.SamePoint(you.Body[len(you.Body)-1], you.Body[len(you.Body)-2])

	// chase food if we just ate, think we aren't girthy enough, or are going to starve.
	// if foundSnack && (justFed || needMoreGirth || (int(you.Health)-10 < len(tastiestSnackPath))) {
	if foundSnack {
		return generator.DirectionToPoint(you.Body[0], tastiestSnackPath[0]), "chasing snack"
	}

	route, err := p.GetRoute(you.Body[0], you.Body[len(you.Body)-1])
	if err == nil {
		return generator.DirectionToPoint(you.Body[0], route[0]), "chasing tail"
	}

	// if we can't reach our own tail, get all points and calculate longest route
	allAvailableRoutes := pather.GetRoutesFromOrigin(state, you.Body[0], you.Body[0])

	var longestRouteLength int
	var longestRoute []rules.Point
	for _, route := range allAvailableRoutes {
		if len(route) <= longestRouteLength {
			continue
		}
		longestRouteLength = len(route)
		longestRoute = route
	}

	if len(longestRoute) == 0 {
		return generator.DirectionDown, "no safe routes. GG"

	}

	// TODO: make this longest route wind around into the space we have.
	return generator.DirectionToPoint(you.Body[0], longestRoute[0]), "doing longest route"

}

func SafeMoves(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake) [4]int {

	moves := generator.AllMovesForState(state)
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
				Move: generator.DirectionToPoint(state.Snakes[moveSnakeNum].Body[0], movePoint).String(),
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

		direction := generator.DirectionToPoint(you.Body[0], move[youPosition])
		safeMoves[direction]++

	}

	return safeMoves

}

func SafestMoves(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake) []generator.Direction {

	safeMoves := SafeMoves(state, ruleset, you)

	var safestMoveCount int
	for _, move := range safeMoves {
		if move > safestMoveCount {
			safestMoveCount = move
		}
	}

	safestMoves := []generator.Direction{}
	for direction, move := range safeMoves {
		if move == safestMoveCount {
			safestMoves = append(safestMoves, generator.Direction(direction))
		}
	}

	return safestMoves
}

func MoveIsSafe(state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, candidate generator.Direction) bool {
	safeMoves := SafeMoves(state, ruleset, you)
	if safeMoves[candidate] > 0 {
		return true
	}

	return false
}
