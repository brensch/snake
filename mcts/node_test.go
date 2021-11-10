package mcts

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

// func TestPopulateChildMoveSets(t *testing.T) {
// 	state := []byte(`{"Turn":23,"Height":11,"Width":11,"Food":[{"X":10,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":3,"Y":3},{"X":3,"Y":4},{"X":2,"Y":4},{"X":1,"Y":4},{"X":0,"Y":4}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1","Body":[{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1}],"Health":86,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"2","Body":[{"X":1,"Y":9},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9}],"Health":89,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3","Body":[{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

// 	var s *rules.BoardState
// 	err := json.Unmarshal(state, &s)
// 	if err != nil {
// 		fmt.Print("err", err.Error())
// 		t.Fail()
// 	}

// 	generator.PrintMap(s)

// 	node := &Node{
// 		State: s,
// 	}

// 	node.PopulateChildMoveSets()

// 	if len(node.Children) != 81 {
// 		t.Log("node child length not correct", len(node.Children))
// 		t.Fail()
// 	}

// 	for _, child := range node.Children {
// 		if len(child.MoveSet) == 0 {
// 			t.Log("got empty moveset")
// 			t.Fail()
// 		}
// 	}

// }

func TestGenerateStateAndScoreChild(t *testing.T) {
	state := []byte(`{"Turn":23,"Height":11,"Width":11,"Food":[{"X":10,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":3,"Y":4},{"X":2,"Y":4},{"X":1,"Y":4},{"X":0,"Y":4}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1","Body":[{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1}],"Health":86,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"2","Body":[{"X":1,"Y":9},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9}],"Health":89,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3","Body":[{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	node := &Node{
		State: s,
	}
	generator.PrintMap(s)
	_ = node

	// node.PopulateChildMoveSets()

	// ruleset := &rules.StandardRuleset{
	// 	FoodSpawnChance: 0,
	// 	MinimumFood:     0,
	// }

	// start := time.Now()
	// for i := 0; i < len(node.Children); i++ {

	// 	err = node.GenerateStateAndScoreChild(ruleset, i)
	// 	if err != nil {
	// 		t.Log(err)
	// 		t.FailNow()
	// 	}

	// 	if node.Children[i].State == nil {
	// 		t.Log("got nil state after move")
	// 		t.Fail()
	// 	}

	// 	if node.Children[i].Score != 0 {
	// 		fmt.Println("got non zero score", node.Children[i].Score)
	// 		generator.PrintMap(node.Children[i].State)
	// 	}
	// }
	// fmt.Println("took useconds", time.Since(start).Microseconds())

}

func TestSearchNode(t *testing.T) {
	state := []byte(`{"Turn":23,"Height":11,"Width":11,"Food":[{"X":10,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":3,"Y":4},{"X":2,"Y":4},{"X":1,"Y":4},{"X":0,"Y":4}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1","Body":[{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1}],"Health":86,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"2","Body":[{"X":1,"Y":9},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9}],"Health":89,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3","Body":[{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	// node := &Node{
	// 	State: s,
	// }
	// generator.PrintMap(s)

	// start := time.Now()
	// for i := 0; i < 100; i++ {
	// 	node.Search(0, 20)
	// }

	// fmt.Println("took useconds", time.Since(start).Microseconds())

	// for _, child := range node.Children {
	// 	fmt.Println(child.Score)
	// }
}
