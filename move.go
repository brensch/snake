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

	// var closestSkippedSnack rules.Point

	for _, snack := range state.Food {
		// fmt.Println("checking snack", snack)

		route, _, err := pather.GetRoute(state, ruleset, you.Body[0], snack, you.ID)
		if err != nil {
			continue
		}

		// healthCost := routedGrid[snack.X][snack.Y].CostFromOrigin

		// if healthCost/pather.CostFactor > you.Health {
		// 	fmt.Println("too hungry for dat boy")
		// }

		ffState := generator.FastForward(state, ruleset, you, route)

		// squaresFromSnackOnwards := pather.CountSquaresReachableFromOrigin(ffState, route[0], you.ID)
		squaresFromSnackOnwards, snackOnwardsGrid := pather.GetReachablePoints(ffState, route[0], you.ID)

		snackOnwardsGrid.DebugPrint()
		//
		// fmt.Println(snack, squaresFromSnackOnwards)
		generator.PrintMap(ffState)
		if len(squaresFromSnackOnwards) < len(you.Body) {
			// if len(routesFromSnackOnwards) < len(you.Body) {
			// log.WithFields(log.Fields{
			// 	"reachable": squaresFromSnackOnwards,
			// 	"total":     len(you.Body),
			// 	"snack":     snack,
			// }).Debug("not enough room to fit ya boi if i chase that snack")
			fmt.Printf("can't fit if %+v\n", snack)

			continue
		}

		// check if other is likely to get this first
		otherSnakeCloser := false
		for _, snake := range state.Snakes {
			if snake.ID == you.ID {
				continue
			}

			opponentRouteToSnack, _, err := pather.GetRoute(state, ruleset, snake.Body[0], snack, snake.ID)
			if err != nil {
				continue
			}

			// scenario where we'd be the same size and kill them over it (very common)
			if len(opponentRouteToSnack) == len(route) && len(you.Body) > len(snake.Body) {
				continue
			}

			// if they are closer, don't chase
			if len(opponentRouteToSnack) <= len(route) {
				otherSnakeCloser = true
				fmt.Printf("other snake closer to %+v\n", snack)
				break
			}

		}

		if otherSnakeCloser {
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

	// if no snack, target the center most point from the reachable points
	reachablePoints, availableSquaresGrid := pather.GetReachablePoints(state, you.Body[0], you.ID)

	if len(reachablePoints) == 0 {
		return generator.DirectionDown, "no reachable points. GG"
	}

	// availableSquaresGrid.DebugPrint()

	// centrePoint := generator.CentreMostPoint(state, reachablePoints)
	furthestPoint := availableSquaresGrid.FurthestPoint()

	route, _, err := pather.GetRoute(state, ruleset, you.Body[0], furthestPoint, you.ID)
	if err == nil {
		return generator.DirectionToPoint(you.Body[0], route[len(route)-1]), "going to furthest point"
	}

	// // try to chase tail
	// route, _, err := pather.GetRoute(state, ruleset, you.Body[0], you.Body[len(you.Body)-1], you.ID)
	// if err == nil {
	// 	return generator.DirectionToPoint(you.Body[0], route[len(route)-1]), "chasing tail"
	// }

	// find the longest path to potentially find a better way.
	// approach is to find longest path, fast forward, and keep going for an amount of time i'll decide later
	// reachablePoints := pather.GetReachablePoints(state, you.Body[0], you.ID)
	// fmt.Println("furthest points", reachablePoints)
	// reachablePoints, availableSquaresGrid := pather.GetReachablePoints(state, you.Body[0], you.ID)

	// if len(reachablePoints) == 0 {
	// 	return generator.DirectionDown, "no reachable points. GG"
	// }

	// // if there are heaps of reachable points, we should just pick the furthest away point and go for there
	// if len(reachablePoints) > len(you.Body) {
	// 	var furthestPoint rules.Point
	// 	var furthestPointDistance int32
	// 	for _, point := range reachablePoints {
	// 		distance := availableSquaresGrid[point.X][point.Y].StepsFromOrigin
	// 		if distance > furthestPointDistance {
	// 			furthestPoint = point
	// 			furthestPointDistance = distance
	// 		}
	// 	}

	// 	// aim at that furthest point
	// 	route, _, err := pather.GetRoute(state, ruleset, you.Body[0], furthestPoint, you.ID)
	// 	if err == nil {
	// 		// fmt.Println(route)
	// 		return generator.DirectionToPoint(you.Body[0], route[len(route)-1]), "can't reach tail but plenty of room. going for furthest point."
	// 	}
	// }

	// if the reachable points with a* pathing is smaller than the length of our snake,
	// search for the longest path instead.

	// recursive function. needs to explore until finished
	longestPath := availableSquaresGrid.ExploreForLength([]rules.Point{you.Body[0]}, len(you.Body))

	// fmt.Println(longestPath)
	// for _,

	// old logic for fforwarding and checking return path. worked okayish

	// with the list of points, fastforward to each one and see if we can go further from there
	// longestRouteLength := 0
	// var longestRoute []rules.Point
	// for _, point := range reachablePoints {
	// 	route, _, err := pather.GetRoute(state, ruleset, you.Body[0], point, you.ID)
	// 	if err != nil {
	// 		fmt.Println("weird, this point should be reachable")
	// 		continue
	// 	}

	// 	// check route to start in case we can't actually make it back (avoids giving no routes)
	// 	if len(route) > longestRouteLength {
	// 		// Only need to store the initial route, since we only need first step anyway
	// 		longestRoute = route
	// 		longestRouteLength = len(route)
	// 	}

	// 	ffState := generator.FastForward(state, ruleset, you, route)

	// 	// try to chase tail again now we've fastforwarded to furthest point
	// 	// TODO: this is not your real tail. maybe subtract len of route.
	// 	routeToTail, _, err := pather.GetRoute(ffState, ruleset, point, you.Body[len(you.Body)-1], you.ID)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	// fmt.Println("could reach tail after", route)

	// 	potentialRouteLength := len(route) + len(routeToTail)
	// 	if potentialRouteLength > longestRouteLength {
	// 		// Only need to store the initial route, since we only need first step anyway
	// 		longestRoute = route
	// 		longestRouteLength = potentialRouteLength
	// 	}

	// }

	// if len(longestRoute) < len(you.Body) {
	// 	fmt.Println("doomed.")
	// }

	// fmt.Println("reachos", len(reachablePoints))
	// if len(reachablePoints) >
	// fmt.Println(longestPath)

	if len(longestPath) > 1 {
		// fmt.Println("got longest route", longestPath)
		return generator.DirectionToPoint(you.Body[0], longestPath[1]), "doing intentionally long path"
	}

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
	return generator.DirectionDown, "no routes left. GG."

}

func Move(ctx context.Context, state *rules.BoardState, ruleset rules.Ruleset, you rules.Snake, turn int32, gameID string) (generator.Direction, string) {
	galaxyBrain, reason := GalaxyBrain(ctx, state, ruleset, you, turn)
	safestMoves := generator.SafestMoves(state, ruleset, you)

	finalMove := galaxyBrain

	galaxyBrainSafe := false
	for _, smoothBrain := range safestMoves {
		if galaxyBrain == smoothBrain {
			galaxyBrainSafe = true
			break
		}
	}

	if len(safestMoves) != 0 && !galaxyBrainSafe {
		// fmt.Println("made a smooth brain move..............................................")
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
