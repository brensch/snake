package minimax

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
	log "github.com/sirupsen/logrus"
)

type testCaseMinimax struct {
	explanation string
	state       []byte
	scoreMax    float64
	scoreMin    float64
}

var (
	testsMinimax = []testCaseMinimax{
		// {
		// 	explanation: "check heading towards food",
		// 	state:       []byte(`{"Turn":0,"Height":11,"Width":11,"Food":[{"X":0,"Y":8},{"X":6,"Y":10},{"X":5,"Y":5}],"Snakes":[{"ID":"max","Body":[{"X":1,"Y":9},{"X":1,"Y":8},{"X":1,"Y":7}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"min","Body":[{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"X":2,"Y":9},{"X":2,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 	scoreMin:    0.66,
		// 	scoreMax:    0.67,
		// },
		// {
		// 	explanation: "head towards space",
		// 	state:       []byte(`{"Turn":19,"Height":11,"Width":11,"Food":[{"X":0,"Y":4}],"Snakes":[{"ID":"you","Body":[{"X":6,"Y":8},{"X":7,"Y":8},{"X":8,"Y":8}],"Health":82,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"d66ed448-851d-416e-bc26-851906eaee5d","Body":[{"X":8,"Y":4},{"X":9,"Y":4},{"X":9,"Y":3},{"X":8,"Y":3},{"X":7,"Y":3},{"X":6,"Y":3},{"X":6,"Y":2}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		{
			explanation: "head towards space",
			state:       []byte(`{"Turn":10,"Height":11,"Width":11,"Food":[{"X":0,"Y":4},{"X":6,"Y":2}],"Snakes":[{"ID":"91fd088f-938d-4c1a-a074-85dca741a019","Body":[{"X":6,"Y":7},{"X":5,"Y":7},{"X":4,"Y":7}],"Health":91,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3d4ac849-6cc3-4785-990a-9c6e619e538c","Body":[{"X":5,"Y":4},{"X":5,"Y":5},{"X":4,"Y":5},{"X":4,"Y":4},{"X":4,"Y":3}],"Health":99,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			// okMoves:     []rules.Direction{rules.DirectionDown},
		},
	}
)

// func TestMinimax(t *testing.T) {

// 	log.SetLevel(log.DebugLevel)

// 	for _, test := range testsMinimax {
// 		t.Log("running test: ", test.explanation)

// 		var s *rules.BoardState
// 		err := json.Unmarshal(test.state, &s)
// 		if err != nil {
// 			t.Error(err)
// 			t.FailNow()
// 		}

// 		generator.PrintMap(s)
// 		ruleset := &rules.StandardRuleset{
// 			FoodSpawnChance:     50,
// 			MinimumFood:         0,
// 			HazardDamagePerTurn: 16,
// 		}

// 		bestChild, score := Search(0, s, ruleset, 15, math.Inf(-1), math.Inf(1))
// 		fmt.Println("got score of next move", score)
// 		generator.PrintMap(bestChild)
// 	}
// }

func TestMinimaxNode(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	for _, test := range testsMinimax {
		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)
		ruleset := &rules.StandardRuleset{
			FoodSpawnChance:     50,
			MinimumFood:         0,
			HazardDamagePerTurn: 16,
		}

		startingNode := &Node{
			Alpha:        math.Inf(-1),
			Beta:         math.Inf(1),
			IsMaximising: true,
			State:        s,
		}

		startingNode.Search(3, ruleset)
		fmt.Println("got startingnode children", len(startingNode.Children))
		fmt.Println("got score", *startingNode.Score)

		currentNode := startingNode
		nextLevel := currentNode.Children
		level := 0
		for len(nextLevel) > 0 {
			fmt.Println("------------------------LEVEL", level)

			var newNextLevel []*Node
			for _, child := range nextLevel {
				child.Print()
				generator.PrintMap(child.State)
				newNextLevel = append(newNextLevel, child.Children...)

			}

			nextLevel = newNextLevel
			level++

		}

		// nextChild := startingNode
		// for {
		// 	// fmt.Println("got score of", *nextChild.Score)
		// 	temp := nextChild.FindBestChild()
		// 	if temp == nil {
		// 		break
		// 	}
		// 	nextChild = temp
		// 	nextChild.Print()
		// 	generator.PrintMap(nextChild.State)
		// }
		// generator.PrintMap(bestChild)
	}
}
