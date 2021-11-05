package pather

import (
	"encoding/json"
	"testing"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
	log "github.com/sirupsen/logrus"
)

type killTestCase struct {
	explanation string
	state       []byte
	okMoves     []rules.Direction
	isKillable  bool
}

// https://play.battlesnake.com/g/b6b113a8-711a-4e5e-905c-81674814afe2/
// turn 112, purple should attack in, but green can also still get out

var (
	tests = []killTestCase{
		{
			explanation: "check box into wall",
			state:       []byte(`{"Turn":49,"Height":11,"Width":11,"Hazards":null,"Food":[{"X":1,"Y":3},{"Y":9,"X":7}],"Snakes":[{"Body":[{"Y":5,"X":10},{"Y":6,"X":10},{"Y":7,"X":10},{"Y":8,"X":10},{"Y":9,"X":10},{"Y":9,"X":9},{"X":9,"Y":9}],"Health":100,"ID":"them","EliminatedOnTurn":0,"EliminatedCause":"","EliminatedBy":""},{"EliminatedCause":"","ID":"you","Body":[{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"Y":4,"X":5},{"X":4,"Y":4},{"Y":4,"X":3},{"X":2,"Y":4}],"Health":91,"EliminatedBy":"","EliminatedOnTurn":0},{"EliminatedCause":"","EliminatedOnTurn":0,"ID":"gs_88K6HT8qbvwV3x6rMcrDw3d3","Health":61,"Body":[{"X":2,"Y":9},{"X":3,"Y":9},{"Y":9,"X":4},{"X":4,"Y":8}],"EliminatedBy":""}]}`),
			okMoves:     []rules.Direction{rules.DirectionRight},
			isKillable:  true,
		},
		// {
		// 	explanation: "kill instantly",
		// 	state:       []byte(`{"Turn":49,"Height":11,"Width":11,"Hazards":null,"Food":[{"X":1,"Y":3},{"Y":9,"X":7}],"Snakes":[{"Body":[{"Y":5,"X":10},{"Y":6,"X":10},{"Y":7,"X":10},{"Y":8,"X":10},{"Y":9,"X":10},{"Y":9,"X":9},{"X":9,"Y":9}],"Health":100,"ID":"them","EliminatedOnTurn":0,"EliminatedCause":"","EliminatedBy":""},{"EliminatedCause":"","ID":"you","Body":[{"X":9,"Y":4},{"X":9,"Y":5},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"Y":4,"X":5},{"X":4,"Y":4},{"Y":4,"X":3},{"X":2,"Y":4}],"Health":91,"EliminatedBy":"","EliminatedOnTurn":0},{"EliminatedCause":"","EliminatedOnTurn":0,"ID":"gs_88K6HT8qbvwV3x6rMcrDw3d3","Health":61,"Body":[{"X":2,"Y":9},{"X":3,"Y":9},{"Y":9,"X":4},{"X":4,"Y":8}],"EliminatedBy":""}]}`),
		// 	okMoves:     []generator.Direction{generator.DirectionRight},
		// 	isKillable:  true,
		// },
		// {
		// 	explanation: "check follow along wall",
		// 	state:       []byte(`{"Turn":49,"Height":11,"Width":11,"Hazards":null,"Food":[{"X":1,"Y":3},{"Y":9,"X":7}],"Snakes":[{"Body":[{"Y":5,"X":10},{"Y":6,"X":10},{"Y":7,"X":10},{"Y":8,"X":10},{"Y":9,"X":10},{"Y":9,"X":9},{"X":9,"Y":9}],"Health":100,"ID":"them","EliminatedOnTurn":0,"EliminatedCause":"","EliminatedBy":""},{"EliminatedCause":"","ID":"you","Body":[{"X":9,"Y":3},{"X":9,"Y":4},{"X":9,"Y":5},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"Y":4,"X":5},{"X":4,"Y":4},{"Y":4,"X":3},{"X":2,"Y":4}],"Health":91,"EliminatedBy":"","EliminatedOnTurn":0},{"EliminatedCause":"","EliminatedOnTurn":0,"ID":"gs_88K6HT8qbvwV3x6rMcrDw3d3","Health":61,"Body":[{"X":2,"Y":9},{"X":3,"Y":9},{"Y":9,"X":4},{"X":4,"Y":8}],"EliminatedBy":""}]}`),
		// 	okMoves:     []generator.Direction{generator.DirectionRight},
		// 	isKillable:  true,
		// },
	}
)

func TestIdentifyKill(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	youID := "you"
	themID := "them"

	for _, test := range tests {

		// if test.explanation != "is your opponent hungry" {
		// 	continue
		// }

		t.Log("running test: ", test.explanation)

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

		them, err := generator.GetYou(s, themID)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		// ruleset := &rules.StandardRuleset{
		// 	FoodSpawnChance: 0,
		// 	MinimumFood:     1,
		// }

		killRoute, err := IdentifyKill(s, them, you)
		if err != nil {
			if !test.isKillable {
				continue
			}

			t.Log("failed to identify killable opponent")
			t.Fail()
			continue
		}

		direction := generator.DirectionToPoint(you.Body[0], killRoute[0])

		moveOk := false
		for _, okMove := range test.okMoves {
			if direction == okMove {
				moveOk = true
				break
			}
		}

		if moveOk {
			continue
		}

		t.Logf("got %s. ok moves: %+v", direction, test.okMoves)
		t.Fail()

	}

}
