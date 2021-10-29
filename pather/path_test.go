package pather

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

var (
	addObstacleTests = [][]byte{
		[]byte(`{"Turn":0,"Height":11,"Width":11,"Food":[{"X":0,"Y":8},{"X":6,"Y":10},{"X":5,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":1,"Y":9},{"X":1,"Y":8},{"X":1,"Y":7}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"7dd375fc-c66e-413e-aa11-bd13d32bbef4","Body":[{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"X":2,"Y":9},{"X":2,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		[]byte(`{"Turn":179,"Height":11,"Width":11,"Food":[{"X":3,"Y":8},{"X":10,"Y":8},{"X":3,"Y":7},{"X":1,"Y":2}],"Snakes":[{"ID":"you","Body":[{"X":4,"Y":6},{"X":5,"Y":6},{"X":6,"Y":6},{"X":7,"Y":6},{"X":7,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"X":7,"Y":4},{"X":8,"Y":4},{"X":8,"Y":5},{"X":8,"Y":6},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":8,"Y":3},{"X":7,"Y":3},{"X":6,"Y":3},{"X":5,"Y":3}],"Health":87,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"dcf675fe-12ca-40d3-9268-07a2cc747866","Body":[{"X":3,"Y":5},{"X":2,"Y":5},{"X":2,"Y":4},{"X":3,"Y":4},{"X":4,"Y":4},{"X":4,"Y":3},{"X":4,"Y":2},{"X":5,"Y":2},{"X":6,"Y":2},{"X":7,"Y":2},{"X":8,"Y":2},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"X":9,"Y":1},{"X":8,"Y":1},{"X":7,"Y":1},{"X":6,"Y":1}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		[]byte(`{"Turn":203,"Height":11,"Width":11,"Food":[{"X":10,"Y":3},{"X":0,"Y":1}],"Snakes":[{"ID":"you","Body":[{"X":9,"Y":3},{"X":8,"Y":3},{"X":8,"Y":4},{"X":7,"Y":4},{"X":6,"Y":4},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5},{"X":8,"Y":6},{"X":8,"Y":7},{"X":8,"Y":8},{"X":8,"Y":9},{"X":8,"Y":10},{"X":9,"Y":10},{"X":9,"Y":9},{"X":9,"Y":8},{"X":9,"Y":7},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"a2f84b19-fe62-4c7c-b7a6-7c301c1b20ff","Body":[{"X":7,"Y":3},{"X":6,"Y":3},{"X":5,"Y":3},{"X":5,"Y":4},{"X":4,"Y":4},{"X":4,"Y":3},{"X":3,"Y":3},{"X":2,"Y":3},{"X":1,"Y":3},{"X":1,"Y":2},{"X":1,"Y":1},{"X":1,"Y":0},{"X":2,"Y":0},{"X":3,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":6,"Y":0},{"X":7,"Y":0},{"X":8,"Y":0},{"X":8,"Y":1},{"X":7,"Y":1}],"Health":80,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		[]byte(`{"Turn":117,"Height":11,"Width":11,"Food":[{"X":10,"Y":10},{"X":7,"Y":10}],"Snakes":[{"ID":"you","Body":[{"X":7,"Y":9},{"X":8,"Y":9},{"X":9,"Y":9},{"X":10,"Y":9},{"X":10,"Y":8},{"X":10,"Y":7},{"X":10,"Y":6},{"X":10,"Y":5},{"X":10,"Y":4},{"X":9,"Y":4},{"X":9,"Y":5},{"X":8,"Y":5},{"X":8,"Y":6}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"8f58be28-81aa-4046-b6f0-3b530d9fc4b6","Body":[{"X":5,"Y":9},{"X":5,"Y":8},{"X":6,"Y":8},{"X":7,"Y":8},{"X":7,"Y":7},{"X":7,"Y":6},{"X":6,"Y":6},{"X":6,"Y":7},{"X":5,"Y":7},{"X":4,"Y":7},{"X":4,"Y":8},{"X":3,"Y":8},{"X":3,"Y":9}],"Health":78,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		[]byte(`{"Hazards":null,"Food":[{"X":10,"Y":0},{"Y":8,"X":0},{"X":6,"Y":10}],"Turn":125,"Snakes":[{"EliminatedBy":"","Health":84,"ID":"you","EliminatedOnTurn":0,"Body":[{"X":9,"Y":0},{"Y":0,"X":8},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"X":5,"Y":1},{"Y":2,"X":5},{"Y":2,"X":6},{"Y":2,"X":7},{"X":8,"Y":2},{"X":8,"Y":1}],"EliminatedCause":""},{"EliminatedBy":"","ID":"gs_d8XdgWyYJQ37ddytQMhDGvh4","EliminatedOnTurn":0,"EliminatedCause":"","Body":[{"Y":7,"X":6},{"X":5,"Y":7},{"Y":7,"X":4},{"X":3,"Y":7},{"Y":7,"X":2},{"X":2,"Y":8},{"X":3,"Y":8},{"Y":8,"X":4},{"Y":8,"X":5},{"X":6,"Y":8},{"Y":8,"X":7},{"X":8,"Y":8},{"Y":7,"X":8},{"X":8,"Y":6}],"Health":93}],"Width":11,"Height":11}`),
		[]byte(`{"Height":11,"Turn":543,"Snakes":[{"EliminatedOnTurn":0,"Health":4,"EliminatedBy":"","ID":"you","Body":[{"Y":1,"X":8},{"X":7,"Y":1},{"Y":1,"X":6},{"X":5,"Y":1},{"Y":1,"X":4},{"X":3,"Y":1},{"Y":1,"X":2},{"Y":1,"X":1},{"Y":1,"X":0},{"Y":2,"X":0},{"Y":2,"X":1},{"Y":3,"X":1},{"Y":3,"X":0},{"Y":4,"X":0},{"Y":4,"X":1},{"X":2,"Y":4},{"Y":5,"X":2},{"Y":5,"X":1},{"X":1,"Y":6},{"X":1,"Y":7},{"X":1,"Y":8},{"Y":9,"X":1},{"Y":10,"X":1},{"X":2,"Y":10},{"Y":9,"X":2},{"X":2,"Y":8},{"Y":7,"X":2},{"X":2,"Y":6},{"Y":6,"X":3},{"Y":5,"X":3},{"Y":4,"X":3},{"Y":4,"X":4},{"Y":4,"X":5},{"Y":3,"X":5},{"X":4,"Y":3},{"X":4,"Y":2},{"Y":2,"X":5},{"Y":2,"X":6},{"Y":2,"X":7},{"X":8,"Y":2},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"Y":1,"X":9}],"EliminatedCause":""},{"Health":89,"EliminatedCause":"","Body":[{"Y":8,"X":5},{"Y":8,"X":4},{"X":4,"Y":9},{"X":4,"Y":10},{"X":3,"Y":10},{"Y":9,"X":3},{"X":3,"Y":8},{"Y":7,"X":3},{"Y":7,"X":4},{"Y":6,"X":4},{"Y":6,"X":5},{"X":6,"Y":6},{"Y":5,"X":6},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":3,"X":8},{"X":9,"Y":3},{"X":10,"Y":3},{"X":10,"Y":4},{"Y":5,"X":10},{"X":10,"Y":6},{"X":10,"Y":7},{"X":10,"Y":8},{"Y":9,"X":10},{"X":9,"Y":9},{"Y":8,"X":9},{"X":9,"Y":7},{"X":8,"Y":7},{"X":8,"Y":6},{"X":7,"Y":6},{"Y":7,"X":7},{"Y":7,"X":6},{"Y":8,"X":6},{"X":6,"Y":9},{"X":7,"Y":9},{"Y":8,"X":7},{"X":8,"Y":8},{"Y":9,"X":8},{"Y":10,"X":8}],"ID":"gs_JS7rvSbPKG9wmcmT7cScMKdX","EliminatedBy":"","EliminatedOnTurn":0}],"Food":[{"Y":8,"X":0},{"X":0,"Y":5},{"X":2,"Y":2},{"X":0,"Y":9},{"Y":0,"X":10},{"Y":7,"X":0},{"X":8,"Y":0},{"X":4,"Y":0},{"Y":3,"X":2},{"Y":2,"X":3}],"Hazards":null,"Width":11}`),
	}
	getRouteTests = []getRouteTest{
		{
			description: "correctly path to snack",
			state:       []byte(`{"Food":[{"Y":7,"X":9},{"X":10,"Y":0}],"Hazards":null,"Snakes":[{"EliminatedBy":"","Body":[{"Y":5,"X":3},{"Y":6,"X":3},{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"Y":9,"X":4},{"X":5,"Y":9},{"X":5,"Y":10},{"Y":10,"X":6},{"Y":9,"X":6},{"Y":9,"X":7},{"X":8,"Y":9},{"Y":9,"X":9},{"X":10,"Y":9},{"X":10,"Y":8},{"X":9,"Y":8},{"X":8,"Y":8},{"X":7,"Y":8}],"EliminatedOnTurn":0,"Health":92,"ID":"you","EliminatedCause":""},{"EliminatedCause":"","Health":98,"ID":"gs_XXp6TM7X8QXRBmXycrTVk7SP","Body":[{"X":1,"Y":7},{"Y":6,"X":1},{"X":0,"Y":6},{"Y":5,"X":0},{"Y":4,"X":0},{"Y":4,"X":1},{"X":2,"Y":4},{"Y":4,"X":3},{"Y":3,"X":3},{"X":3,"Y":2},{"X":3,"Y":1},{"Y":1,"X":2},{"Y":0,"X":2},{"Y":0,"X":3},{"X":4,"Y":0}],"EliminatedOnTurn":0,"EliminatedBy":""}],"Height":11,"Turn":138,"Width":11}`),
			target:      rules.Point{X: 9, Y: 7},
		},
		{
			description: "don't go into dead end",
			state:       []byte(`{"Turn":86,"Height":11,"Width":11,"Food":[{"X":10,"Y":0}],"Snakes":[{"ID":"you","Body":[{"X":1,"Y":0},{"X":1,"Y":1},{"X":1,"Y":2},{"X":1,"Y":3},{"X":1,"Y":4},{"X":0,"Y":4},{"X":0,"Y":5},{"X":0,"Y":6},{"X":1,"Y":6},{"X":2,"Y":6},{"X":3,"Y":6},{"X":4,"Y":6},{"X":5,"Y":6},{"X":6,"Y":6},{"X":6,"Y":7}],"Health":96,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"d480d032-4d70-4ad6-aa88-0b725a1e852e","Body":[{"X":7,"Y":0},{"X":7,"Y":1},{"X":8,"Y":1},{"X":8,"Y":2},{"X":8,"Y":3},{"X":8,"Y":4},{"X":9,"Y":4},{"X":10,"Y":4},{"X":10,"Y":3},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1}],"Health":99,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			target:      rules.Point{X: 9, Y: 7},
		},
	}
)

type getRouteTest struct {
	description  string
	state        []byte
	target       rules.Point
	expectedCost int
}

func TestAddObstacles(t *testing.T) {
	youID := "you"

	for _, test := range addObstacleTests {

		var s *rules.BoardState
		err := json.Unmarshal(test, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		you, err := generator.GetYou(s, youID)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		grid := initPathGrid(s)
		grid.AddObstacles(s, you.Body[0])

		grid.DebugPrint()

		// p.DebugPrint(you.Body[0], you.Body[len(you.Body)-1])

	}

}

func TestPath(t *testing.T) {
	youID := "you"

	for _, test := range getRouteTests {

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		you, err := generator.GetYou(s, youID)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		route, grid, err := GetRoute(s, &rules.StandardRuleset{}, you.Body[0], test.target)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("proposed route to target", route)

		grid.DebugPrint()

	}

}

// func TestAddObstacles(t *testing.T) {

// 	var state *rules.BoardState
// 	err := json.Unmarshal(start, &state)
// 	if err != nil {
// 		t.Error(err)
// 		t.FailNow()
// 	}

// 	generator.PrintMap(state)

// 	p := MakePathgrid(state, you.Body[0])

// }
