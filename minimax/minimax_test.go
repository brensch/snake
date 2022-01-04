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
		// {
		// 	explanation: "head towards space",
		// 	state:       []byte(`{"Turn":10,"Height":11,"Width":11,"Food":[{"X":0,"Y":4},{"X":6,"Y":2}],"Snakes":[{"ID":"91fd088f-938d-4c1a-a074-85dca741a019","Body":[{"X":6,"Y":7},{"X":5,"Y":7},{"X":4,"Y":7}],"Health":91,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3d4ac849-6cc3-4785-990a-9c6e619e538c","Body":[{"X":5,"Y":4},{"X":5,"Y":5},{"X":4,"Y":5},{"X":4,"Y":4},{"X":4,"Y":3}],"Health":99,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "don't die by not knowing enemy head will kill us",
		// 	state:       []byte(`{"Turn":43,"Height":11,"Width":11,"Food":[{"X":5,"Y":5},{"X":4,"Y":8}],"Snakes":[{"ID":"b3ea9698-5dd1-4d5d-b070-dc3bc6a5fd51","Body":[{"X":8,"Y":2},{"X":9,"Y":2},{"X":9,"Y":3},{"X":8,"Y":3}],"Health":60,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"bc2aaf94-7822-4a51-b551-074f9c144da5","Body":[{"X":9,"Y":1},{"X":10,"Y":1},{"X":10,"Y":0},{"X":9,"Y":0},{"X":8,"Y":0}],"Health":95,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "test we don't keep iterating if we know we can't go deeper",
		// 	state:       []byte(`{"Turn":112,"Height":11,"Width":11,"Food":[{"X":2,"Y":8}],"Snakes":[{"ID":"5b51d8d5-d387-4d38-a93b-a13904d95972","Body":[{"X":4,"Y":7},{"X":3,"Y":7},{"X":2,"Y":7},{"X":1,"Y":7}],"Health":1,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"a9498ded-b27a-4b0b-8fe0-2d7c0b776168","Body":[{"X":4,"Y":5},{"X":4,"Y":4},{"X":4,"Y":3},{"X":4,"Y":2},{"X":4,"Y":1},{"X":4,"Y":0},{"X":3,"Y":0},{"X":2,"Y":0},{"X":2,"Y":0}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "test we count a draw as a loss",
		// 	state: []byte(`
		// {"Turn":22,"Height":11,"Width":11,"Food":[{"X":5,"Y":5},{"X":4,"Y":2},{"X":3,"Y":0}],"Snakes":[{"ID":"0061ca3b-c9c6-4053-9443-27ba2bccf63f","Body":[{"X":10,"Y":9},{"X":9,"Y":9},{"X":8,"Y":9},{"X":7,"Y":9}],"Health":81,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"93341b3b-b576-4e9a-a379-295de64c4965","Body":[{"X":9,"Y":10},{"X":8,"Y":10},{"X":7,"Y":10},{"X":6,"Y":10}],"Health":81,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "don't move somewhere we could be killed",
		// 	state: []byte(`
		// {"Turn":69,"Height":11,"Width":11,"Food":[{"X":1,"Y":7},{"X":6,"Y":7}],"Snakes":[{"ID":"7e987b10-4ed5-4d58-ac31-32a0b02b2a91","Body":[{"X":3,"Y":5},{"X":3,"Y":4},{"X":3,"Y":3},{"X":4,"Y":3}],"Health":85,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"81e8c0bd-4a9a-4982-8ad7-528ed008ae5b","Body":[{"X":2,"Y":4},{"X":2,"Y":3},{"X":2,"Y":2},{"X":2,"Y":1},{"X":3,"Y":1},{"X":4,"Y":1},{"X":5,"Y":1},{"X":6,"Y":1},{"X":7,"Y":1},{"X":8,"Y":1},{"X":9,"Y":1},{"X":10,"Y":1},{"X":10,"Y":2}],"Health":91,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "don't starve, caused by iterative depening giving strange output",
		// 	state: []byte(`
		// {"Turn":279,"Height":11,"Width":11,"Food":[{"X":9,"Y":1},{"X":4,"Y":0},{"X":3,"Y":2},{"X":1,"Y":2},{"X":4,"Y":9}],"Snakes":[{"ID":"6b81dd4c-fbd8-487b-b5bd-744d8b9f47cc","Body":[{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5},{"X":8,"Y":6},{"X":7,"Y":6},{"X":6,"Y":6},{"X":5,"Y":6}],"Health":6,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1f00c485-a941-4d1c-b03a-e53f58f9f429","Body":[{"X":5,"Y":7},{"X":6,"Y":7},{"X":6,"Y":8},{"X":7,"Y":8},{"X":8,"Y":8}],"Health":14,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "do not segfault",
		// 	state: []byte(`
		// {"Turn":615,"Height":11,"Width":11,"Food":[{"X":5,"Y":10},{"X":8,"Y":8},{"X":7,"Y":10},{"X":0,"Y":2}],"Snakes":[{"ID":"bbb49307-1b5a-4f26-904c-2df0766b96a2","Body":[{"X":1,"Y":1},{"X":1,"Y":2},{"X":1,"Y":3},{"X":1,"Y":4},{"X":1,"Y":5},{"X":2,"Y":5},{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":6},{"X":5,"Y":6},{"X":4,"Y":6},{"X":3,"Y":6},{"X":3,"Y":7},{"X":2,"Y":7},{"X":1,"Y":7},{"X":0,"Y":7},{"X":0,"Y":8},{"X":1,"Y":8}],"Health":78,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"f0f36033-e6f5-4ae6-a406-1a58015d6526","Body":[{"X":0,"Y":0},{"X":1,"Y":0},{"X":2,"Y":0},{"X":3,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":6,"Y":0},{"X":7,"Y":0},{"X":8,"Y":0},{"X":9,"Y":0},{"X":10,"Y":0},{"X":10,"Y":1},{"X":9,"Y":1},{"X":8,"Y":1},{"X":7,"Y":1},{"X":6,"Y":1},{"X":5,"Y":1},{"X":5,"Y":1}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "don't hit other snake going for food",
		// 	state: []byte(`
		// {"Turn":104,"Height":11,"Width":11,"Food":[{"X":0,"Y":2}],"Snakes":[{"ID":"c88f4037-48cc-47ff-b8d7-0028697a0dd6","Body":[{"X":8,"Y":7},{"X":7,"Y":7},{"X":7,"Y":6},{"X":6,"Y":6},{"X":6,"Y":7},{"X":5,"Y":7}],"Health":10,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"661374af-51f4-4fd1-96aa-cdf8a51ad61b","Body":[{"X":9,"Y":6},{"X":9,"Y":7},{"X":10,"Y":7},{"X":10,"Y":6}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "test go down to fill space.",
		// 	state: []byte(`
		// 			{"Turn":179,"Height":11,"Width":11,"Food":[{"X":9,"Y":8},{"X":9,"Y":6}],"Snakes":[{"ID":"a335824d-2942-4088-a79a-9ce28111e048","Body":[{"X":6,"Y":10},{"X":7,"Y":10},{"X":8,"Y":10},{"X":9,"Y":10},{"X":10,"Y":10},{"X":10,"Y":9},{"X":9,"Y":9}],"Health":31,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"23e3bcb5-9d9b-4c24-afea-c47bb55f8230","Body":[{"X":8,"Y":2},{"X":8,"Y":1},{"X":7,"Y":1},{"X":6,"Y":1},{"X":5,"Y":1},{"X":4,"Y":1},{"X":3,"Y":1},{"X":2,"Y":1},{"X":1,"Y":1},{"X":0,"Y":1},{"X":0,"Y":0},{"X":1,"Y":0},{"X":2,"Y":0},{"X":3,"Y":0}],"Health":46,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 			`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "don't go into tight space, symptom of not clearing old tree",
		// 	state: []byte(`
		// 	{"Turn":639,"Height":11,"Width":11,"Food":[{"X":9,"Y":10},{"X":10,"Y":10},{"X":1,"Y":8},{"X":3,"Y":7}],"Snakes":[{"ID":"4a4c77da-f10c-4064-bc48-3141562e788e","Body":[{"X":5,"Y":3},{"X":5,"Y":4},{"X":6,"Y":4},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5},{"X":8,"Y":6},{"X":9,"Y":6},{"X":9,"Y":7},{"X":8,"Y":7},{"X":7,"Y":7},{"X":6,"Y":7},{"X":5,"Y":7},{"X":4,"Y":7},{"X":4,"Y":6},{"X":5,"Y":6},{"X":5,"Y":5},{"X":4,"Y":5},{"X":3,"Y":5}],"Health":80,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"890785b2-4c44-4be8-914c-1d7cbda05612","Body":[{"X":2,"Y":6},{"X":1,"Y":6},{"X":1,"Y":5},{"X":1,"Y":4},{"X":0,"Y":4},{"X":0,"Y":3},{"X":1,"Y":3},{"X":1,"Y":2},{"X":2,"Y":2},{"X":3,"Y":2},{"X":4,"Y":2},{"X":5,"Y":2},{"X":5,"Y":1},{"X":6,"Y":1},{"X":7,"Y":1},{"X":8,"Y":1},{"X":8,"Y":2},{"X":8,"Y":2}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 	`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "check iterating too much",
		// 	state: []byte(`
		// 	{"Turn":27,"Height":11,"Width":11,"Food":[{"X":1,"Y":7},{"X":7,"Y":3}],"Snakes":[{"ID":"fc21e710-4a00-44fc-8703-a6e124cf295f","Body":[{"X":7,"Y":1},{"X":6,"Y":1},{"X":5,"Y":1},{"X":4,"Y":1},{"X":3,"Y":1}],"Health":84,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"7c17761e-e9de-4877-a2f3-0bf5310c2a22","Body":[{"X":6,"Y":0},{"X":5,"Y":0},{"X":4,"Y":0},{"X":3,"Y":0},{"X":2,"Y":0}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 	`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		// {
		// 	explanation: "don't panic",
		// 	state: []byte(`
		// 	{"Turn":8,"Height":11,"Width":11,"Food":[{"X":5,"Y":5},{"X":5,"Y":4}],"Snakes":[{"ID":"gs_WmdKBYdfcypWSbQSt3pgtTGH","Body":[{"X":5,"Y":7},{"X":4,"Y":7},{"X":3,"Y":7},{"X":2,"Y":7}],"Health":94,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"gs_TGYB6JcwF7c4fMCQBJDBpVpF","Body":[{"X":5,"Y":3},{"X":5,"Y":2},{"X":6,"Y":2},{"X":7,"Y":2}],"Health":94,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
		// 	`),
		// 	// okMoves:     []rules.Direction{rules.DirectionDown},
		// },
		{
			explanation: "doom on the left",
			state: []byte(`
			{"Turn":392,"Height":11,"Width":11,"Food":[{"X":1,"Y":0},{"X":0,"Y":0},{"X":1,"Y":6},{"X":9,"Y":5},{"X":1,"Y":1}],"Snakes":[{"ID":"gs_XTH4t36GjgDWYxH68cWkc3Cc","Body":[{"X":6,"Y":0},{"X":6,"Y":1},{"X":6,"Y":2},{"X":7,"Y":2},{"X":7,"Y":3},{"X":6,"Y":3},{"X":5,"Y":3},{"X":5,"Y":4},{"X":6,"Y":4},{"X":6,"Y":5},{"X":6,"Y":6},{"X":5,"Y":6},{"X":4,"Y":6},{"X":3,"Y":6},{"X":3,"Y":7},{"X":4,"Y":7},{"X":5,"Y":7},{"X":6,"Y":7},{"X":7,"Y":7},{"X":7,"Y":6},{"X":7,"Y":5},{"X":7,"Y":4},{"X":8,"Y":4},{"X":9,"Y":4},{"X":9,"Y":3},{"X":8,"Y":3},{"X":8,"Y":2},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"X":9,"Y":1},{"X":8,"Y":1},{"X":8,"Y":1}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"gs_3WTmFGhBX6jHPqQ7TF799cg4","Body":[{"X":1,"Y":9},{"X":1,"Y":10},{"X":2,"Y":10},{"X":3,"Y":10},{"X":4,"Y":10},{"X":5,"Y":10},{"X":6,"Y":10},{"X":7,"Y":10},{"X":8,"Y":10},{"X":9,"Y":10},{"X":10,"Y":10},{"X":10,"Y":9},{"X":9,"Y":9},{"X":9,"Y":8},{"X":8,"Y":8},{"X":7,"Y":8},{"X":7,"Y":9},{"X":6,"Y":9},{"X":5,"Y":9},{"X":4,"Y":9},{"X":3,"Y":9},{"X":3,"Y":8},{"X":2,"Y":8},{"X":2,"Y":7},{"X":2,"Y":6},{"X":2,"Y":5},{"X":2,"Y":4},{"X":2,"Y":3},{"X":2,"Y":2},{"X":3,"Y":2}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}
			`),
			// okMoves: []rules.Direction{rules.DirectionRight},
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

		deepestDepth, _ := startingNode.Search(ctx, 15, 15, ruleset, nil)
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

		// startingNode.ExploreBestPath()
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
			FoodSpawnChance:     15,
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
