package minimax

import (
	"encoding/json"
	"testing"

	"github.com/brensch/snake/rules"
	log "github.com/sirupsen/logrus"
)

func TestHash(t *testing.T) {
	state := []byte(`{"Turn":631,"Height":11,"Width":11,"Food":[{"X":8,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":2,"Y":0},{"X":0,"Y":10},{"X":2,"Y":1},{"X":10,"Y":2}],"Snakes":[{"ID":"f92aec56-1423-4a13-8a79-77b7df6d55ab","Body":[{"X":2,"Y":10},{"X":3,"Y":10},{"X":4,"Y":10},{"X":5,"Y":10},{"X":6,"Y":10},{"X":7,"Y":10},{"X":8,"Y":10},{"X":9,"Y":10},{"X":10,"Y":10},{"X":10,"Y":9},{"X":9,"Y":9},{"X":8,"Y":9},{"X":7,"Y":9},{"X":6,"Y":9},{"X":5,"Y":9},{"X":4,"Y":9},{"X":3,"Y":9},{"X":2,"Y":9},{"X":1,"Y":9},{"X":0,"Y":9},{"X":0,"Y":8},{"X":1,"Y":8},{"X":2,"Y":8},{"X":3,"Y":8},{"X":4,"Y":8}],"Health":94,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1a285af4-be1c-4135-9c85-acaaeeb59f6e","Body":[{"X":8,"Y":6},{"X":9,"Y":6},{"X":10,"Y":6},{"X":10,"Y":5},{"X":9,"Y":5},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"X":4,"Y":5},{"X":3,"Y":5},{"X":2,"Y":5},{"X":1,"Y":5},{"X":0,"Y":5},{"X":0,"Y":4},{"X":1,"Y":4},{"X":2,"Y":4},{"X":3,"Y":4}],"Health":55,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	differentState := []byte(`{"Turn":631,"Height":11,"Width":11,"Food":[{"X":8,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":2,"Y":0},{"X":0,"Y":10},{"X":2,"Y":1},{"X":10,"Y":2}],"Snakes":[{"ID":"f92aec56-1423-4a13-8a79-77b7df6d55ab","Body":[{"X":3,"Y":9},{"X":3,"Y":10},{"X":4,"Y":10},{"X":5,"Y":10},{"X":6,"Y":10},{"X":7,"Y":10},{"X":8,"Y":10},{"X":9,"Y":10},{"X":10,"Y":10},{"X":10,"Y":9},{"X":9,"Y":9},{"X":8,"Y":9},{"X":7,"Y":9},{"X":6,"Y":9},{"X":5,"Y":9},{"X":4,"Y":9},{"X":3,"Y":9},{"X":2,"Y":9},{"X":1,"Y":9},{"X":0,"Y":9},{"X":0,"Y":8},{"X":1,"Y":8},{"X":2,"Y":8},{"X":3,"Y":8},{"X":4,"Y":8}],"Health":94,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1a285af4-be1c-4135-9c85-acaaeeb59f6e","Body":[{"X":8,"Y":6},{"X":9,"Y":6},{"X":10,"Y":6},{"X":10,"Y":5},{"X":9,"Y":5},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"X":4,"Y":5},{"X":3,"Y":5},{"X":2,"Y":5},{"X":1,"Y":5},{"X":0,"Y":5},{"X":0,"Y":4},{"X":1,"Y":4},{"X":2,"Y":4},{"X":3,"Y":4}],"Health":55,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var s2 *rules.BoardState
	err = json.Unmarshal(state, &s2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	mappo := make(map[uint64]bool)

	mappo[Hash(s)] = true

	if !mappo[Hash(s2)] {
		t.Log("hashes didn't match")
		t.FailNow()
	}

	var s3 *rules.BoardState
	err = json.Unmarshal(differentState, &s3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if mappo[Hash(s3)] {
		t.Log("this should be different")
		t.FailNow()
	}

}

func BenchmarkHash(b *testing.B) {
	b.StopTimer()
	state := []byte(`{"Turn":631,"Height":11,"Width":11,"Food":[{"X":8,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":2,"Y":0},{"X":0,"Y":10},{"X":2,"Y":1},{"X":10,"Y":2}],"Snakes":[{"ID":"f92aec56-1423-4a13-8a79-77b7df6d55ab","Body":[{"X":2,"Y":10},{"X":3,"Y":10},{"X":4,"Y":10},{"X":5,"Y":10},{"X":6,"Y":10},{"X":7,"Y":10},{"X":8,"Y":10},{"X":9,"Y":10},{"X":10,"Y":10},{"X":10,"Y":9},{"X":9,"Y":9},{"X":8,"Y":9},{"X":7,"Y":9},{"X":6,"Y":9},{"X":5,"Y":9},{"X":4,"Y":9},{"X":3,"Y":9},{"X":2,"Y":9},{"X":1,"Y":9},{"X":0,"Y":9},{"X":0,"Y":8},{"X":1,"Y":8},{"X":2,"Y":8},{"X":3,"Y":8},{"X":4,"Y":8}],"Health":94,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"1a285af4-be1c-4135-9c85-acaaeeb59f6e","Body":[{"X":8,"Y":6},{"X":9,"Y":6},{"X":10,"Y":6},{"X":10,"Y":5},{"X":9,"Y":5},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"X":4,"Y":5},{"X":3,"Y":5},{"X":2,"Y":5},{"X":1,"Y":5},{"X":0,"Y":5},{"X":0,"Y":4},{"X":1,"Y":4},{"X":2,"Y":4},{"X":3,"Y":4}],"Health":55,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`)

	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	b.StartTimer()

	for n := 0; n < b.N; n++ {

		Hash(s)
	}

}
