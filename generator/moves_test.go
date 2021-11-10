package generator

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/brensch/snake/rules"
)

func TestAllMovesForState(t *testing.T) {
	state := []byte(`{"Hazards":null,"Food":[{"Y":9,"X":4},{"Y":10,"X":10},{"X":3,"Y":7}],"Width":11,"Snakes":[{"ID":"gs_Mryy8QHFPhDq3cxWJJpRFkg4","EliminatedCause":"","EliminatedOnTurn":0,"Health":78,"Body":[{"Y":6,"X":3},{"X":3,"Y":5},{"X":3,"Y":4},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":3,"X":1},{"X":0,"Y":3},{"Y":4,"X":0},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"X":0,"Y":8},{"Y":8,"X":1},{"X":2,"Y":8},{"Y":7,"X":2},{"Y":7,"X":1}],"EliminatedBy":""},{"Body":[{"X":9,"Y":2},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"Y":1,"X":5},{"X":5,"Y":2},{"Y":3,"X":5},{"X":5,"Y":4},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":4,"X":9},{"X":10,"Y":4},{"Y":5,"X":10},{"Y":5,"X":9},{"Y":5,"X":8},{"Y":5,"X":7},{"X":7,"Y":6}],"EliminatedOnTurn":0,"EliminatedCause":"","Health":97,"EliminatedBy":"","ID":"you"}],"Height":11,"Turn":163}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	moves := AllMoveSetsForState(s)

	fmt.Println(moves)

}

func TestSpeedAllMovesForState(t *testing.T) {
	state := []byte(`{"Hazards":null,"Food":[{"Y":9,"X":4},{"Y":10,"X":10},{"X":3,"Y":7}],"Width":11,"Snakes":[{"ID":"gs_Mryy8QHFPhDq3cxWJJpRFkg4","EliminatedCause":"","EliminatedOnTurn":0,"Health":78,"Body":[{"Y":6,"X":3},{"X":3,"Y":5},{"X":3,"Y":4},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":3,"X":1},{"X":0,"Y":3},{"Y":4,"X":0},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"X":0,"Y":8},{"Y":8,"X":1},{"X":2,"Y":8},{"Y":7,"X":2},{"Y":7,"X":1}],"EliminatedBy":""},{"Body":[{"X":9,"Y":2},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"Y":1,"X":5},{"X":5,"Y":2},{"Y":3,"X":5},{"X":5,"Y":4},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":4,"X":9},{"X":10,"Y":4},{"Y":5,"X":10},{"Y":5,"X":9},{"Y":5,"X":8},{"Y":5,"X":7},{"X":7,"Y":6}],"EliminatedOnTurn":0,"EliminatedCause":"","Health":97,"EliminatedBy":"","ID":"you"}],"Height":11,"Turn":163}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	start := time.Now()
	count := 0

	for {

		if time.Since(start) > 500*time.Millisecond {
			break
		}
		AllMoveSetsForState(s)
		count++
	}

	fmt.Println(count)

}

func TestSpeedAllMovesForStateRawSpeed(t *testing.T) {

	state := []byte(`{"Hazards":null,"Food":[{"Y":9,"X":4},{"Y":10,"X":10},{"X":3,"Y":7}],"Width":11,"Snakes":[{"ID":"gs_Mryy8QHFPhDq3cxWJJpRFkg4","EliminatedCause":"","EliminatedOnTurn":0,"Health":78,"Body":[{"Y":6,"X":3},{"X":3,"Y":5},{"X":3,"Y":4},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":3,"X":1},{"X":0,"Y":3},{"Y":4,"X":0},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"X":0,"Y":8},{"Y":8,"X":1},{"X":2,"Y":8},{"Y":7,"X":2},{"Y":7,"X":1}],"EliminatedBy":""},{"Body":[{"X":9,"Y":2},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"Y":1,"X":5},{"X":5,"Y":2},{"Y":3,"X":5},{"X":5,"Y":4},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":4,"X":9},{"X":10,"Y":4},{"Y":5,"X":10},{"Y":5,"X":9},{"Y":5,"X":8},{"Y":5,"X":7},{"X":7,"Y":6}],"EliminatedOnTurn":0,"EliminatedCause":"","Health":97,"EliminatedBy":"","ID":"you"}],"Height":11,"Turn":163}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	PrintMap(s)

	start := time.Now()
	for i := 0; i < 1000; i++ {

		AllMoveSetsForStateRaw(s)
	}
	fmt.Println("duration useconds", time.Since(start).Microseconds())

	moves := AllMoveSetsForStateRaw(s)

	fmt.Println(len(moves))

}

func TestSpeedAllMovesForStateRaw(t *testing.T) {
	state := []byte(`{"Turn":23,"Height":11,"Width":11,"Food":[{"X":10,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":3,"Y":4},{"X":2,"Y":4},{"X":1,"Y":4},{"X":0,"Y":4}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1","Body":[{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1}],"Health":86,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"2","Body":[{"X":1,"Y":9},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9}],"Health":89,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3","Body":[{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	// state := []byte(`{"Hazards":null,"Food":[{"Y":9,"X":4},{"Y":10,"X":10},{"X":3,"Y":7}],"Width":11,"Snakes":[{"ID":"gs_Mryy8QHFPhDq3cxWJJpRFkg4","EliminatedCause":"","EliminatedOnTurn":0,"Health":78,"Body":[{"Y":6,"X":3},{"X":3,"Y":5},{"X":3,"Y":4},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":3,"X":1},{"X":0,"Y":3},{"Y":4,"X":0},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"X":0,"Y":8},{"Y":8,"X":1},{"X":2,"Y":8},{"Y":7,"X":2},{"Y":7,"X":1}],"EliminatedBy":""},{"Body":[{"X":9,"Y":2},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"Y":1,"X":5},{"X":5,"Y":2},{"Y":3,"X":5},{"X":5,"Y":4},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":4,"X":9},{"X":10,"Y":4},{"Y":5,"X":10},{"Y":5,"X":9},{"Y":5,"X":8},{"Y":5,"X":7},{"X":7,"Y":6}],"EliminatedOnTurn":0,"EliminatedCause":"","Health":97,"EliminatedBy":"","ID":"you"}],"Height":11,"Turn":163}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	PrintMap(s)

	moves := AllMoveSetsForStateRaw(s)

	fmt.Println(len(moves))

	ruleset := &rules.StandardRuleset{
		FoodSpawnChance: 0,
		MinimumFood:     0,
	}

	for _, move := range moves {
		nextState, err := ruleset.CreateNextBoardState(s, move)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		for _, snake := range nextState.Snakes {
			if snake.EliminatedBy != "" {
				t.Log(snake)
				t.Fail()
			}
		}
	}

}
