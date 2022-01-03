package minimax

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

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
		// 		{
		// 			explanation: "head towards space",
		// 			state:       []byte(`{"Turn":10,"Height":11,"Width":11,"Food":[{"X":0,"Y":4},{"X":6,"Y":2}],"Snakes":[{"ID":"91fd088f-938d-4c1a-a074-85dca741a019","Body":[{"X":6,"Y":7},{"X":5,"Y":7},{"X":4,"Y":7}],"Health":91,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3d4ac849-6cc3-4785-990a-9c6e619e538c","Body":[{"X":5,"Y":4},{"X":5,"Y":5},{"X":4,"Y":5},{"X":4,"Y":4},{"X":4,"Y":3}],"Health":99,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 			// okMoves:     []rules.Direction{rules.DirectionDown},
		// 		},
		// 		{
		// 			explanation: "don't die",
		// 			state:       []byte(`{"Turn":43,"Height":11,"Width":11,"Food":[{"X":5,"Y":5},{"X":4,"Y":8}],"Snakes":[{"ID":"b3ea9698-5dd1-4d5d-b070-dc3bc6a5fd51","Body":[{"X":8,"Y":2},{"X":9,"Y":2},{"X":9,"Y":3},{"X":8,"Y":3}],"Health":60,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"bc2aaf94-7822-4a51-b551-074f9c144da5","Body":[{"X":9,"Y":1},{"X":10,"Y":1},{"X":10,"Y":0},{"X":9,"Y":0},{"X":8,"Y":0}],"Health":95,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 			// okMoves:     []rules.Direction{rules.DirectionDown},
		// 		},
		// 		{
		// 			explanation: "test we don't keep iterating if we know we can't go deeper",
		// 			state:       []byte(`{"Turn":112,"Height":11,"Width":11,"Food":[{"X":2,"Y":8}],"Snakes":[{"ID":"5b51d8d5-d387-4d38-a93b-a13904d95972","Body":[{"X":4,"Y":7},{"X":3,"Y":7},{"X":2,"Y":7},{"X":1,"Y":7}],"Health":1,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"a9498ded-b27a-4b0b-8fe0-2d7c0b776168","Body":[{"X":4,"Y":5},{"X":4,"Y":4},{"X":4,"Y":3},{"X":4,"Y":2},{"X":4,"Y":1},{"X":4,"Y":0},{"X":3,"Y":0},{"X":2,"Y":0},{"X":2,"Y":0}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 			// okMoves:     []rules.Direction{rules.DirectionDown},
		// 		},
		// 		{
		// 			explanation: "test we count a draw as a loss",
		// 			state: []byte(`
		// {"Turn":22,"Height":11,"Width":11,"Food":[{"X":5,"Y":5},{"X":4,"Y":2},{"X":3,"Y":0}],"Snakes":[{"ID":"0061ca3b-c9c6-4053-9443-27ba2bccf63f","Body":[{"X":10,"Y":9},{"X":9,"Y":9},{"X":8,"Y":9},{"X":7,"Y":9}],"Health":81,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"93341b3b-b576-4e9a-a379-295de64c4965","Body":[{"X":9,"Y":10},{"X":8,"Y":10},{"X":7,"Y":10},{"X":6,"Y":10}],"Health":81,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 			// okMoves:     []rules.Direction{rules.DirectionDown},
		// 		},
		{
			explanation: "don't move somewhere we could be killed",
			state: []byte(`
{"Turn":69,"Height":11,"Width":11,"Food":[{"X":1,"Y":7},{"X":6,"Y":7}],"Snakes":[{"ID":"7e987b10-4ed5-4d58-ac31-32a0b02b2a91","Body":[{"X":3,"Y":5},{"X":3,"Y":4},{"X":3,"Y":3},{"X":4,"Y":3}],"Health":85,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"81e8c0bd-4a9a-4982-8ad7-528ed008ae5b","Body":[{"X":2,"Y":4},{"X":2,"Y":3},{"X":2,"Y":2},{"X":2,"Y":1},{"X":3,"Y":1},{"X":4,"Y":1},{"X":5,"Y":1},{"X":6,"Y":1},{"X":7,"Y":1},{"X":8,"Y":1},{"X":9,"Y":1},{"X":10,"Y":1},{"X":10,"Y":2}],"Health":91,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
			`),
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
		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		defer cancel()

		deepestDepth, _ := startingNode.Search(ctx, 16, 16, ruleset, nil)
		fmt.Println("got to depth", deepestDepth)
		fmt.Println("got startingnode children", len(startingNode.Children))
		fmt.Println("got score", *startingNode.Score)

		generator.PrintMap(startingNode.FindBestChild().State)

		// currentNode := startingNode
		// nextLevel := currentNode.Children
		// level := 0
		// for len(nextLevel) > 0 {
		// 	fmt.Println("------------------------LEVEL", level)

		// 	var newNextLevel []*Node
		// 	for _, child := range nextLevel {
		// 		child.Print()
		// 		generator.PrintMap(child.State)
		// 		newNextLevel = append(newNextLevel, child.Children...)

		// 	}

		// 	nextLevel = newNextLevel
		// 	level++

		// }

		nextChild := startingNode
		for {
			// fmt.Println("got score of", *nextChild.Score)
			temp := nextChild.FindBestChild()
			if temp == nil {
				break
			}
			nextChild = temp
			nextChild.Print()
			generator.PrintMap(nextChild.State)
		}
		// generator.PrintMap(bestChild)
	}
}

func TestDeepeningSearch(t *testing.T) {

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

		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()

		bestState := startingNode.DeepeningSearch(ctx, ruleset)
		// fmt.Println("got startingnode children", len(startingNode.Children))
		// fmt.Println("got score", *startingNode.Score)

		generator.PrintMap(&bestState)

		// currentNode := startingNode
		// nextLevel := currentNode.Children
		// level := 0
		// for len(nextLevel) > 0 {
		// 	fmt.Println("------------------------LEVEL", level)

		// 	var newNextLevel []*Node
		// 	for _, child := range nextLevel {
		// 		child.Print()
		// 		generator.PrintMap(child.State)
		// 		newNextLevel = append(newNextLevel, child.Children...)

		// 	}

		// 	nextLevel = newNextLevel
		// 	level++

		// }

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

func BenchmarkCopy(b *testing.B) {

	log.SetLevel(log.DebugLevel)

	for _, test := range testsMinimax {
		b.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			b.Error(err)
			b.FailNow()
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

		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()

		startingNode.Search(ctx, 19, 19, ruleset, nil)
		fmt.Println("got startingnode children", len(startingNode.Children))
		fmt.Println("got score", *startingNode.Score)

		generator.PrintMap(startingNode.FindBestChild().State)

		for n := 0; n < b.N; n++ {

			startingNode.CopyNode()
		}
		// currentNode := startingNode
		// nextLevel := currentNode.Children
		// level := 0
		// for len(nextLevel) > 0 {
		// 	fmt.Println("------------------------LEVEL", level)

		// 	var newNextLevel []*Node
		// 	for _, child := range nextLevel {
		// 		child.Print()
		// 		generator.PrintMap(child.State)
		// 		newNextLevel = append(newNextLevel, child.Children...)

		// 	}

		// 	nextLevel = newNextLevel
		// 	level++

		// }

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
