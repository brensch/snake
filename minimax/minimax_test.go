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
		{
			explanation: "check heading towards food",
			state:       []byte(`{"Turn":0,"Height":11,"Width":11,"Food":[{"X":0,"Y":8},{"X":6,"Y":10},{"X":5,"Y":5}],"Snakes":[{"ID":"max","Body":[{"X":1,"Y":9},{"X":1,"Y":8},{"X":1,"Y":7}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"min","Body":[{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"X":2,"Y":9},{"X":2,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			scoreMin:    0.66,
			scoreMax:    0.67,
		},
	}
)

func TestMinimax(t *testing.T) {

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

		bestChild, score := Search(0, s, ruleset, 15, math.Inf(-1), math.Inf(1))
		fmt.Println("got score of next move", score)
		generator.PrintMap(bestChild)
	}
}
