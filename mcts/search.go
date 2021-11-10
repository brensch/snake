package mcts

import (
	"context"
	"fmt"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

func Search(ctx context.Context, state *rules.BoardState, simulationDepth int) *Node {

	root := NewNode(nil, nil, state)

	for {
		if ctx.Err() != nil {
			return root
		}

		// always start at the root node
		node := root

		// Select.
		// traverse down fully explored routes until you get to a node that needs simulation
		for len(node.UntriedMovesets) == 0 && len(node.Children) > 0 {
			node = node.SelectNextChildToExplore()
		}

		// Expand.
		if len(node.UntriedMovesets) > 0 {
			node = node.MakeRandomUntriedMove()
		}

		ruleset := &rules.StandardRuleset{
			FoodSpawnChance: 0,
			MinimumFood:     0,
		}

		// Simulate.
		iterState := state
		var err error
		var j int
		for j = 0; j < simulationDepth; j++ {
			if ctx.Err() != nil {
				return root
			}
			moveSets := generator.AllMoveSetsForStateRaw(iterState)
			if len(moveSets) == 0 {
				break
			}
			randomMoveSet := moveSets[r1.Intn(len(moveSets))]
			iterState, err = ruleset.CreateNextBoardState(iterState, randomMoveSet)
			if err != nil {
				fmt.Println("shouldn't get err here", err, moveSets)
			}

			if WeDied(iterState) {
				break
			}
		}

		outcome := HeuristicScore(iterState, j)
		// Backpropagate
		node.UpdateOutcomeScore(outcome)
	}

	// fmt.Println(root)
	// for _, child := range root.Children {
	// 	fmt.Println(child.MoveSet, child.OutcomeScore)
	// }

}
