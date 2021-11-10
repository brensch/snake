package mcts

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brensch/snake/rules"
)

func TestSearch(t *testing.T) {

	state := []byte(`{"Turn":23,"Height":11,"Width":11,"Food":[{"X":10,"Y":5}],"Snakes":[{"ID":"you","Body":[{"X":3,"Y":4},{"X":2,"Y":4},{"X":1,"Y":4},{"X":0,"Y":4}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1","Body":[{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1}],"Health":86,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"2","Body":[{"X":1,"Y":9},{"X":2,"Y":9},{"X":3,"Y":9},{"X":4,"Y":9},{"X":5,"Y":9},{"X":6,"Y":9}],"Health":89,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3","Body":[{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":6,"Y":5},{"X":7,"Y":5},{"X":8,"Y":5}],"Health":97,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
		t.Fail()
	}

	// Search(s, 10)
}
