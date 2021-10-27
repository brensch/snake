package main

import (
	"context"
	"encoding/json"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

type testCase struct {
	explanation string
	state       []byte
	okMoves     []generator.Direction
}

var (
	tests = []testCase{
		{
			explanation: "check heading towards food",
			state:       []byte(`{"Turn":0,"Height":11,"Width":11,"Food":[{"X":0,"Y":8},{"X":6,"Y":10},{"X":5,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":1,"Y":9},{"X":1,"Y":8},{"X":1,"Y":7}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"7dd375fc-c66e-413e-aa11-bd13d32bbef4","Body":[{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"X":2,"Y":9},{"X":2,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionLeft},
		},
		{
			explanation: "check going for safe kill and snack",
			state:       []byte(`{"Turn":179,"Height":11,"Width":11,"Food":[{"X":3,"Y":8},{"X":10,"Y":8},{"X":3,"Y":7},{"X":1,"Y":2}],"Snakes":[{"ID":"you","Body":[{"X":4,"Y":6},{"X":5,"Y":6},{"X":6,"Y":6},{"X":7,"Y":6},{"X":7,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"X":7,"Y":4},{"X":8,"Y":4},{"X":8,"Y":5},{"X":8,"Y":6},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":8,"Y":3},{"X":7,"Y":3},{"X":6,"Y":3},{"X":5,"Y":3}],"Health":87,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"dcf675fe-12ca-40d3-9268-07a2cc747866","Body":[{"X":3,"Y":5},{"X":2,"Y":5},{"X":2,"Y":4},{"X":3,"Y":4},{"X":4,"Y":4},{"X":4,"Y":3},{"X":4,"Y":2},{"X":5,"Y":2},{"X":6,"Y":2},{"X":7,"Y":2},{"X":8,"Y":2},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"X":9,"Y":1},{"X":8,"Y":1},{"X":7,"Y":1},{"X":6,"Y":1}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionLeft},
		},
		{
			explanation: "check goes for snack",
			state:       []byte(`{"Turn":203,"Height":11,"Width":11,"Food":[{"X":10,"Y":3},{"X":0,"Y":1}],"Snakes":[{"ID":"you","Body":[{"X":9,"Y":3},{"X":8,"Y":3},{"X":8,"Y":4},{"X":7,"Y":4},{"X":6,"Y":4},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5},{"X":8,"Y":6},{"X":8,"Y":7},{"X":8,"Y":8},{"X":8,"Y":9},{"X":8,"Y":10},{"X":9,"Y":10},{"X":9,"Y":9},{"X":9,"Y":8},{"X":9,"Y":7},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"a2f84b19-fe62-4c7c-b7a6-7c301c1b20ff","Body":[{"X":7,"Y":3},{"X":6,"Y":3},{"X":5,"Y":3},{"X":5,"Y":4},{"X":4,"Y":4},{"X":4,"Y":3},{"X":3,"Y":3},{"X":2,"Y":3},{"X":1,"Y":3},{"X":1,"Y":2},{"X":1,"Y":1},{"X":1,"Y":0},{"X":2,"Y":0},{"X":3,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":6,"Y":0},{"X":7,"Y":0},{"X":8,"Y":0},{"X":8,"Y":1},{"X":7,"Y":1}],"Health":80,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "check going for food even when things seem tight",
			state:       []byte(`{"Turn":117,"Height":11,"Width":11,"Food":[{"X":10,"Y":10},{"X":7,"Y":10}],"Snakes":[{"ID":"you","Body":[{"X":7,"Y":9},{"X":8,"Y":9},{"X":9,"Y":9},{"X":10,"Y":9},{"X":10,"Y":8},{"X":10,"Y":7},{"X":10,"Y":6},{"X":10,"Y":5},{"X":10,"Y":4},{"X":9,"Y":4},{"X":9,"Y":5},{"X":8,"Y":5},{"X":8,"Y":6}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"8f58be28-81aa-4046-b6f0-3b530d9fc4b6","Body":[{"X":5,"Y":9},{"X":5,"Y":8},{"X":6,"Y":8},{"X":7,"Y":8},{"X":7,"Y":7},{"X":7,"Y":6},{"X":6,"Y":6},{"X":6,"Y":7},{"X":5,"Y":7},{"X":4,"Y":7},{"X":4,"Y":8},{"X":3,"Y":8},{"X":3,"Y":9}],"Health":78,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionUp},
		},
		{
			explanation: "test still get snack even when in corner",
			state:       []byte(`{"Hazards":null,"Food":[{"X":10,"Y":0},{"Y":8,"X":0},{"X":6,"Y":10}],"Turn":125,"Snakes":[{"EliminatedBy":"","Health":84,"ID":"you","EliminatedOnTurn":0,"Body":[{"X":9,"Y":0},{"Y":0,"X":8},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"X":5,"Y":1},{"Y":2,"X":5},{"Y":2,"X":6},{"Y":2,"X":7},{"X":8,"Y":2},{"X":8,"Y":1}],"EliminatedCause":""},{"EliminatedBy":"","ID":"gs_d8XdgWyYJQ37ddytQMhDGvh4","EliminatedOnTurn":0,"EliminatedCause":"","Body":[{"Y":7,"X":6},{"X":5,"Y":7},{"Y":7,"X":4},{"X":3,"Y":7},{"Y":7,"X":2},{"X":2,"Y":8},{"X":3,"Y":8},{"Y":8,"X":4},{"Y":8,"X":5},{"X":6,"Y":8},{"Y":8,"X":7},{"X":8,"Y":8},{"Y":7,"X":8},{"X":8,"Y":6}],"Health":93}],"Width":11,"Height":11}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "test not heading into dead end",
			state:       []byte(`{"Turn":174,"Height":11,"Width":11,"Food":[{"X":7,"Y":0}],"Snakes":[{"ID":"78c42638-b4d3-4677-a84a-3d7cd06573d7","Body":[{"X":4,"Y":5},{"X":3,"Y":5},{"X":2,"Y":5},{"X":2,"Y":6},{"X":2,"Y":7},{"X":2,"Y":8},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9},{"X":6,"Y":8},{"X":7,"Y":8},{"X":7,"Y":7},{"X":8,"Y":7},{"X":8,"Y":6},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5}],"Health":95,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"you","Body":[{"X":9,"Y":10},{"X":9,"Y":9},{"X":9,"Y":8},{"X":9,"Y":7},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":9,"Y":0},{"X":8,"Y":0},{"X":8,"Y":1},{"X":8,"Y":2},{"X":7,"Y":2},{"X":7,"Y":3},{"X":7,"Y":4},{"X":6,"Y":4},{"X":6,"Y":3},{"X":5,"Y":3},{"X":4,"Y":3},{"X":3,"Y":3},{"X":3,"Y":3}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionLeft},
		},
		{
			explanation: "test chasing tail when partially obstructed and not going for food in dead end",
			state:       []byte(`{"Turn":180,"Height":11,"Width":11,"Food":[{"X":0,"Y":10}],"Snakes":[{"EliminatedBy":"","Health":100,"ID":"gs_dKfSHcdGFSktdPbvqHRDPHb3","Body":[{"Y":1,"X":9},{"X":9,"Y":0},{"X":8,"Y":0},{"X":7,"Y":0},{"Y":0,"X":6},{"X":5,"Y":0},{"Y":0,"X":4},{"Y":0,"X":3},{"Y":0,"X":2},{"X":1,"Y":0},{"Y":1,"X":1},{"Y":1,"X":2},{"Y":1,"X":3},{"X":4,"Y":1},{"Y":1,"X":5},{"Y":1,"X":5}],"EliminatedCause":"","EliminatedOnTurn":0},{"EliminatedCause":"","EliminatedOnTurn":0,"ID":"you","Health":91,"EliminatedBy":"","Body":[{"Y":8,"X":10},{"Y":8,"X":9},{"Y":8,"X":8},{"Y":8,"X":7},{"Y":8,"X":6},{"Y":8,"X":5},{"X":5,"Y":7},{"X":4,"Y":7},{"Y":7,"X":3},{"Y":7,"X":2},{"X":2,"Y":8},{"X":2,"Y":9},{"Y":9,"X":1},{"X":0,"Y":9},{"Y":8,"X":0},{"Y":7,"X":0},{"Y":6,"X":0},{"Y":6,"X":1},{"X":2,"Y":6},{"Y":6,"X":3},{"X":4,"Y":6},{"X":5,"Y":6},{"X":6,"Y":6},{"X":7,"Y":6},{"X":8,"Y":6},{"Y":6,"X":9},{"Y":6,"X":10},{"X":10,"Y":5}]}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
		{
			explanation: "check you don't cut your available space into two too small chunks if there's food there",
			state:       []byte(`{"Width":11,"Turn":225,"Snakes":[{"EliminatedOnTurn":0,"EliminatedBy":"","EliminatedCause":"","Body":[{"Y":1,"X":6},{"X":5,"Y":1},{"X":5,"Y":2},{"Y":3,"X":5},{"Y":3,"X":4},{"Y":3,"X":3},{"X":2,"Y":3},{"X":1,"Y":3},{"X":0,"Y":3},{"Y":4,"X":0},{"X":0,"Y":5},{"X":0,"Y":6},{"X":0,"Y":7},{"X":0,"Y":8},{"Y":9,"X":0},{"Y":9,"X":1},{"X":2,"Y":9},{"X":3,"Y":9},{"X":3,"Y":8},{"X":3,"Y":7},{"Y":7,"X":4},{"Y":8,"X":4},{"Y":8,"X":5},{"Y":8,"X":6},{"X":6,"Y":7},{"X":5,"Y":7},{"Y":6,"X":5},{"X":4,"Y":6}],"Health":99,"ID":"you"},{"EliminatedBy":"","Health":92,"EliminatedOnTurn":0,"Body":[{"X":8,"Y":9},{"Y":9,"X":9},{"X":9,"Y":8},{"X":8,"Y":8},{"X":7,"Y":8},{"Y":7,"X":7},{"X":7,"Y":6},{"Y":6,"X":8},{"X":9,"Y":6},{"X":9,"Y":5},{"Y":4,"X":9},{"Y":4,"X":8},{"Y":4,"X":7},{"X":6,"Y":4},{"X":6,"Y":3},{"X":7,"Y":3},{"Y":3,"X":8},{"Y":3,"X":9},{"X":10,"Y":3},{"Y":4,"X":10},{"Y":5,"X":10},{"Y":6,"X":10},{"X":10,"Y":7},{"Y":8,"X":10},{"X":10,"Y":9},{"Y":10,"X":10},{"Y":10,"X":9}],"EliminatedCause":"","ID":"gs_gKDv4JhPCrKd4VvfvPXfrtQQ"}],"Hazards":null,"Height":11,"Food":[{"Y":0,"X":6},{"X":7,"Y":5}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight, generator.DirectionUp},
		},
		{
			explanation: "check don't chase snack into dead end",
			state:       []byte(`{"Width":11,"Snakes":[{"Body":[{"Y":8,"X":2},{"Y":8,"X":1},{"X":1,"Y":7},{"Y":7,"X":0},{"Y":6,"X":0},{"X":0,"Y":5},{"X":0,"Y":4},{"Y":3,"X":0},{"X":0,"Y":2},{"X":0,"Y":1},{"Y":1,"X":1},{"X":2,"Y":1},{"X":3,"Y":1},{"Y":1,"X":4},{"X":5,"Y":1},{"Y":2,"X":5},{"Y":3,"X":5},{"X":5,"Y":4},{"X":4,"Y":4}],"ID":"gs_Y9vVWDmfP3whjb6pGqFJrVH9","Health":91,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"Body":[{"X":6,"Y":0},{"X":6,"Y":1},{"X":6,"Y":2},{"Y":3,"X":6},{"X":6,"Y":4},{"X":6,"Y":5},{"Y":5,"X":7},{"Y":5,"X":8},{"Y":6,"X":8},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"Y":0,"X":9},{"Y":0,"X":8}],"ID":"you","EliminatedBy":"","Health":84,"EliminatedCause":"","EliminatedOnTurn":0}],"Turn":128,"Hazards":null,"Height":11,"Food":[{"Y":0,"X":1}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "check don't chase snack into small space",
			state:       []byte(`{"Height":11,"Food":[{"X":7,"Y":10}],"Snakes":[{"EliminatedOnTurn":0,"EliminatedCause":"","Body":[{"X":5,"Y":6},{"Y":7,"X":5},{"Y":8,"X":5},{"Y":9,"X":5},{"Y":10,"X":5},{"Y":10,"X":4},{"X":3,"Y":10},{"X":2,"Y":10},{"Y":10,"X":1},{"X":1,"Y":9},{"X":1,"Y":8},{"Y":8,"X":2},{"X":2,"Y":7},{"Y":6,"X":2},{"X":3,"Y":6}],"Health":89,"EliminatedBy":"","ID":"gs_W8T9gFYgRqhmwpbGMBXKmVHR"},{"Health":73,"ID":"you","Body":[{"X":6,"Y":9},{"X":7,"Y":9},{"X":7,"Y":8},{"Y":8,"X":8},{"X":9,"Y":8},{"Y":8,"X":10},{"Y":7,"X":10},{"Y":6,"X":10},{"Y":5,"X":10},{"X":9,"Y":5},{"X":9,"Y":6},{"X":9,"Y":7},{"Y":7,"X":8},{"X":8,"Y":6},{"X":7,"Y":6},{"Y":7,"X":7},{"X":6,"Y":7},{"X":6,"Y":8}],"EliminatedBy":"","EliminatedCause":"","EliminatedOnTurn":0}],"Width":11,"Hazards":null,"Turn":169}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
	}
)

func TestMove(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	youID := "you"

	for _, test := range tests {

		// if test.explanation != "check don't chase snack into dead end" {
		// 	continue
		// }

		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		you, err := generator.GetYou(s, youID)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		ruleset := &rules.StandardRuleset{
			FoodSpawnChance: 0,
			MinimumFood:     1,
		}

		move, reason := Move(context.Background(), s, ruleset, you, s.Turn, "test")
		moveOk := false
		for _, okMove := range test.okMoves {
			if move == okMove {
				moveOk = true
				break
			}
		}

		if moveOk {
			continue
		}

		t.Logf("got %s because %s. ok moves: %+v", move.String(), reason, test.okMoves)
		generator.PrintMap(s)
		t.Fail()

	}

}

// https://play.battlesnake.com/g/3f16d94d-b6a4-4c6d-940b-e6137db33789/

// kill test:
// https://play.battlesnake.com/g/44a88512-6fb6-4f8b-a0ae-4531441258d5/
// turn 89

// got stuck in hazard sauce:
// https://play.battlesnake.com/g/533cd7bf-390e-427c-851b-d5d8e0bc1bb6/

// should chase tail because kill is imminent left
// https://play.battlesnake.com/g/bfd4133d-1e7b-4e05-90a9-27fe220c0984/
// turn 102

// should get health when low...
// https://play.battlesnake.com/g/67e1fe80-2a45-4a68-932c-3c03b1118320/
// turn 543
