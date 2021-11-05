package generator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brensch/snake/rules"
)

func TestSafetyDance(t *testing.T) {
	state := []byte(`{"Turn":29,"Height":11,"Width":11,"Food":[{"X":10,"Y":5}],"Snakes":[{"ID":"d317ceae-5ec7-4b10-80b7-89928d8d04ea","Body":[{"X":3,"Y":3},{"X":3,"Y":4},{"X":2,"Y":4},{"X":1,"Y":4},{"X":0,"Y":4}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"cf9f7cfe-2200-4433-a2b9-e5f8a8bc28f2","Body":[{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1}],"Health":86,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"c8b58917-e8f2-4b62-80f1-8a89411c8502","Body":[{"X":1,"Y":9},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9}],"Health":89,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"6e776816-e325-4572-af52-ecc0c6b9c8b6","Body":[{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	PrintMap(s)
	ruleset := &rules.StandardRuleset{
		FoodSpawnChance: 0,
		MinimumFood:     0,
	}

	for i := 0; i < 4; i++ {
		fmt.Printf("%s ", rules.Direction(i).String())
	}
	fmt.Println()

	fmt.Println(SafeMoves(s, ruleset, "you"))
}
