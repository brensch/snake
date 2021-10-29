package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/pather"
	"github.com/brensch/snake/rules"
)

const LargestCost = 10000

func GalaxyBrain(ctx context.Context, state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, turn int32) (generator.Direction, string) {

	var tastiestSnackPath []rules.Point
	foundSnack := false

	// p := pather.MakePathgrid(state, you.Body[0], you.Body[0])
	// freeSquares := p.FreeSquares(state)
	// hazardCost := 0

	// if ruleset.Name() == "royale" {
	// 	royaleRules, ok := ruleset.(*rules.RoyaleRuleset)
	// 	if ok {
	// 		hazardCost += int(royaleRules.HazardDamagePerTurn) * pather.CostFactor
	// 	}
	// }

	for _, snack := range state.Food {
		// fmt.Println("checking snack", snack)

		route, routedGrid, err := pather.GetRoute(state, ruleset, you.Body[0], snack)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		// fmt.Println(route)
		healthCost := routedGrid[snack.X][snack.Y].CostFromOrigin

		if healthCost/pather.CostFactor > you.Health {
			fmt.Println("too hungry for dat boy")
		}

		// to find the routes from this point on, fastforward your position to what it would be
		// given this route, and make everyone else do their safest move
		// TODO: actually try and assume that they make good moves here
		nextState := state.Clone()

		previousHead := you.Body[0]

		// need to traverse route backwards
		for routePosition := range route {
			pointInRoute := route[len(route)-routePosition-1]

			moves := []rules.SnakeMove{
				{ID: you.ID, Move: generator.DirectionToPoint(previousHead, pointInRoute).String()},
			}
			previousHead = pointInRoute

			for _, snake := range nextState.Snakes {
				if snake.ID == you.ID {
					continue
				}

				safestMoves := SafestMoves(nextState, ruleset, snake)
				if len(safestMoves) == 0 {
					moves = append(moves, rules.SnakeMove{ID: snake.ID, Move: generator.DirectionDown.String()})
					continue
				}
				moves = append(moves, rules.SnakeMove{ID: snake.ID, Move: safestMoves[0].String()})

			}

			nextState, err = ruleset.CreateNextBoardState(nextState, moves)

		}

		squaresFromSnackOnwards := pather.CountSquaresReachableFromOrigin(nextState, route[len(route)-1])

		if squaresFromSnackOnwards < int32(len(you.Body)) {
			// if len(routesFromSnackOnwards) < len(you.Body) {
			log.WithFields(log.Fields{
				"reachable": squaresFromSnackOnwards,
				"total":     len(you.Body),
				"snack":     snack,
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
		return generator.DirectionToPoint(you.Body[0], tastiestSnackPath[len(tastiestSnackPath)-1]), "chasing snack"
	}

	route, _, err := pather.GetRoute(state, ruleset, you.Body[0], you.Body[len(you.Body)-1])
	if err == nil {
		return generator.DirectionToPoint(you.Body[0], route[0]), "chasing tail"
	}

	// TODO: redo this
	// // if we can't reach our own tail, get all points and calculate longest route
	// allAvailableRoutes := pather.GetRoutesFromOrigin(state, you.Body[0], you.Body[0], hazardCost)

	// var longestRouteLength int
	// var longestRoute []rules.Point
	// for _, route := range allAvailableRoutes {
	// 	if len(route) <= longestRouteLength {
	// 		continue
	// 	}
	// 	longestRouteLength = len(route)
	// 	longestRoute = route
	// }

	// if len(longestRoute) == 0 {
	// 	return generator.DirectionDown, "no safe routes. GG"
	// }

	// TODO: make this longest route wind around into the space we have.
	// return generator.DirectionToPoint(you.Body[0], longestRoute), "doing longest route"
	return generator.DirectionDown, "yeet todo logic"

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

func Move(ctx context.Context, state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, turn int32, gameID string) (generator.Direction, string) {
	galaxyBrain, reason := GalaxyBrain(ctx, state, ruleset, you, turn)
	safestMoves := SafestMoves(state, ruleset, you)

	finalMove := galaxyBrain

	galaxyBrainSafe := false
	for _, smoothBrain := range safestMoves {
		if galaxyBrain == smoothBrain {
			galaxyBrainSafe = true
			break
		}
	}

	if len(safestMoves) != 0 && !galaxyBrainSafe {
		finalMove = safestMoves[0]
	}

	safeMoveStrings := []string{}
	for _, move := range safestMoves {
		safeMoveStrings = append(safeMoveStrings, move.String())
	}

	log.WithFields(log.Fields{
		"game":        gameID,
		"action":      "move",
		"galaxy":      galaxyBrain.String(),
		"galaxy_safe": galaxyBrainSafe,
		"safe":        safeMoveStrings,
		"actual":      finalMove.String(),
		"reason":      reason,
		"state":       state,
	}).Info("moved")

	return finalMove, reason

}
