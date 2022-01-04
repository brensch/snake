package minimax

import (
	"encoding/json"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

type testCase struct {
	explanation string
	state       []byte
	scoreMax    float64
	scoreMin    float64
}

var (
	tests = []testCase{
		{
			explanation: "check heading towards food",
			state:       []byte(`{"Turn":0,"Height":11,"Width":11,"Food":[{"X":0,"Y":8},{"X":6,"Y":10},{"X":5,"Y":5}],"Snakes":[{"ID":"not_you","Body":[{"X":1,"Y":9},{"X":1,"Y":8},{"X":1,"Y":7}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"you","Body":[{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"X":2,"Y":9},{"X":2,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			scoreMin:    0.66,
			scoreMax:    0.67,
		},
	}
)

func TestPercentageOfBoardControlled(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	for _, test := range tests {
		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		score := PercentageOfBoardControlled(s)

		t.Log(score)

		if score > test.scoreMax {
			t.Log("score too high")
			t.FailNow()
		}

		if score < test.scoreMin {
			t.Log("score too low")
			t.FailNow()
		}
	}
}

func TestPercentageOfBoardControlledSmart(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	for _, test := range tests {
		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		score := PercentageOfBoardControlledSmart(s)

		t.Log(score)

		if score > test.scoreMax {
			t.Log("score too high")
			t.FailNow()
		}

		if score < test.scoreMin {
			t.Log("score too low")
			t.FailNow()
		}
	}
}

func BenchmarkPercentageOfBoardControlled(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		PercentageOfBoardControlled(s)

	}
}

func BenchmarkPercentageOfBoardControlledSmart(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		PercentageOfBoardControlledSmart(s)

	}
}

func BenchmarkGameFinished(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		GameFinished(s)

	}
}

func BenchmarkGameFinishedBits(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		GameFinishedBits(1, 2)

	}
}

func TestShortestPaths(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	for _, test := range tests {
		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		ShortestPaths(s)

		// PrintShortestPath(graph)

	}
}
func TestShortestPaths2(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	for _, test := range tests {
		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		ShortestPaths2(s)

		// PrintShortestPath(graph)

	}
}

func TestShortestPathsBreadth(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	for _, test := range tests {
		t.Log("running test: ", test.explanation)

		var s *rules.BoardState
		err := json.Unmarshal(test.state, &s)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		generator.PrintMap(s)

		controlledPercent := ShortestPathsBreadth(s)
		fmt.Println(controlledPercent)

		// PrintShortestPath(graph)

	}
}

func BenchmarkShortestPaths(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		ShortestPaths(s)

	}
}

func BenchmarkShortestPaths2(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		ShortestPaths2(s)

	}
}

func BenchmarkShortestPathsBreadth(b *testing.B) {
	log.SetLevel(log.DebugLevel)

	var s *rules.BoardState
	err := json.Unmarshal(tests[0].state, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for n := 0; n < b.N; n++ {

		ShortestPathsBreadth(s)

	}
}
