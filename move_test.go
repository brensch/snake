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
			okMoves:     []generator.Direction{generator.DirectionLeft, generator.DirectionUp},
		},
		// TODO: kill
		// {
		// 	explanation: "check going for safe kill and snack - should turn into opponent path to snack",
		// 	state:       []byte(`{"Turn":179,"Height":11,"Width":11,"Food":[{"X":3,"Y":8},{"X":10,"Y":8},{"X":3,"Y":7},{"X":1,"Y":2}],"Snakes":[{"ID":"you","Body":[{"X":4,"Y":6},{"X":5,"Y":6},{"X":6,"Y":6},{"X":7,"Y":6},{"X":7,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"X":7,"Y":4},{"X":8,"Y":4},{"X":8,"Y":5},{"X":8,"Y":6},{"X":9,"Y":6},{"X":9,"Y":5},{"X":9,"Y":4},{"X":9,"Y":3},{"X":8,"Y":3},{"X":7,"Y":3},{"X":6,"Y":3},{"X":5,"Y":3}],"Health":87,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"dcf675fe-12ca-40d3-9268-07a2cc747866","Body":[{"X":3,"Y":5},{"X":2,"Y":5},{"X":2,"Y":4},{"X":3,"Y":4},{"X":4,"Y":4},{"X":4,"Y":3},{"X":4,"Y":2},{"X":5,"Y":2},{"X":6,"Y":2},{"X":7,"Y":2},{"X":8,"Y":2},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"X":9,"Y":1},{"X":8,"Y":1},{"X":7,"Y":1},{"X":6,"Y":1}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
		// 	okMoves:     []generator.Direction{generator.DirectionLeft},
		// },
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
		// TODO: add longest path algorithm done properly
		// {
		// 	explanation: "check you don't cut your available space into two too small chunks if there's food there",
		// 	state:       []byte(`{"Width":11,"Turn":225,"Snakes":[{"EliminatedOnTurn":0,"EliminatedBy":"","EliminatedCause":"","Body":[{"Y":1,"X":6},{"X":5,"Y":1},{"X":5,"Y":2},{"Y":3,"X":5},{"Y":3,"X":4},{"Y":3,"X":3},{"X":2,"Y":3},{"X":1,"Y":3},{"X":0,"Y":3},{"Y":4,"X":0},{"X":0,"Y":5},{"X":0,"Y":6},{"X":0,"Y":7},{"X":0,"Y":8},{"Y":9,"X":0},{"Y":9,"X":1},{"X":2,"Y":9},{"X":3,"Y":9},{"X":3,"Y":8},{"X":3,"Y":7},{"Y":7,"X":4},{"Y":8,"X":4},{"Y":8,"X":5},{"Y":8,"X":6},{"X":6,"Y":7},{"X":5,"Y":7},{"Y":6,"X":5},{"X":4,"Y":6}],"Health":99,"ID":"you"},{"EliminatedBy":"","Health":92,"EliminatedOnTurn":0,"Body":[{"X":8,"Y":9},{"Y":9,"X":9},{"X":9,"Y":8},{"X":8,"Y":8},{"X":7,"Y":8},{"Y":7,"X":7},{"X":7,"Y":6},{"Y":6,"X":8},{"X":9,"Y":6},{"X":9,"Y":5},{"Y":4,"X":9},{"Y":4,"X":8},{"Y":4,"X":7},{"X":6,"Y":4},{"X":6,"Y":3},{"X":7,"Y":3},{"Y":3,"X":8},{"Y":3,"X":9},{"X":10,"Y":3},{"Y":4,"X":10},{"Y":5,"X":10},{"Y":6,"X":10},{"X":10,"Y":7},{"Y":8,"X":10},{"X":10,"Y":9},{"Y":10,"X":10},{"Y":10,"X":9}],"EliminatedCause":"","ID":"gs_gKDv4JhPCrKd4VvfvPXfrtQQ"}],"Hazards":null,"Height":11,"Food":[{"Y":0,"X":6},{"X":7,"Y":5}]}`),
		// 	okMoves:     []generator.Direction{generator.DirectionRight, generator.DirectionUp},
		// },
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
		// https://play.battlesnake.com/g/67e1fe80-2a45-4a68-932c-3c03b1118320/
		// turn 543
		{
			explanation: "should seek food when hungry",
			state:       []byte(`{"Height":11,"Turn":543,"Snakes":[{"EliminatedOnTurn":0,"Health":4,"EliminatedBy":"","ID":"you","Body":[{"Y":1,"X":8},{"X":7,"Y":1},{"Y":1,"X":6},{"X":5,"Y":1},{"Y":1,"X":4},{"X":3,"Y":1},{"Y":1,"X":2},{"Y":1,"X":1},{"Y":1,"X":0},{"Y":2,"X":0},{"Y":2,"X":1},{"Y":3,"X":1},{"Y":3,"X":0},{"Y":4,"X":0},{"Y":4,"X":1},{"X":2,"Y":4},{"Y":5,"X":2},{"Y":5,"X":1},{"X":1,"Y":6},{"X":1,"Y":7},{"X":1,"Y":8},{"Y":9,"X":1},{"Y":10,"X":1},{"X":2,"Y":10},{"Y":9,"X":2},{"X":2,"Y":8},{"Y":7,"X":2},{"X":2,"Y":6},{"Y":6,"X":3},{"Y":5,"X":3},{"Y":4,"X":3},{"Y":4,"X":4},{"Y":4,"X":5},{"Y":3,"X":5},{"X":4,"Y":3},{"X":4,"Y":2},{"Y":2,"X":5},{"Y":2,"X":6},{"Y":2,"X":7},{"X":8,"Y":2},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"Y":1,"X":9}],"EliminatedCause":""},{"Health":89,"EliminatedCause":"","Body":[{"Y":8,"X":5},{"Y":8,"X":4},{"X":4,"Y":9},{"X":4,"Y":10},{"X":3,"Y":10},{"Y":9,"X":3},{"X":3,"Y":8},{"Y":7,"X":3},{"Y":7,"X":4},{"Y":6,"X":4},{"Y":6,"X":5},{"X":6,"Y":6},{"Y":5,"X":6},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":3,"X":8},{"X":9,"Y":3},{"X":10,"Y":3},{"X":10,"Y":4},{"Y":5,"X":10},{"X":10,"Y":6},{"X":10,"Y":7},{"X":10,"Y":8},{"Y":9,"X":10},{"X":9,"Y":9},{"Y":8,"X":9},{"X":9,"Y":7},{"X":8,"Y":7},{"X":8,"Y":6},{"X":7,"Y":6},{"Y":7,"X":7},{"Y":7,"X":6},{"Y":8,"X":6},{"X":6,"Y":9},{"X":7,"Y":9},{"Y":8,"X":7},{"X":8,"Y":8},{"Y":9,"X":8},{"Y":10,"X":8}],"ID":"gs_JS7rvSbPKG9wmcmT7cScMKdX","EliminatedBy":"","EliminatedOnTurn":0}],"Food":[{"Y":8,"X":0},{"X":0,"Y":5},{"X":2,"Y":2},{"X":0,"Y":9},{"Y":0,"X":10},{"Y":7,"X":0},{"X":8,"Y":0},{"X":4,"Y":0},{"Y":3,"X":2},{"Y":2,"X":3}],"Hazards":null,"Width":11}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
		{
			explanation: "should not smoothbrain into corner",
			state:       []byte(`{"Height":11,"Hazards":null,"Snakes":[{"ID":"you","EliminatedOnTurn":0,"EliminatedBy":"","Body":[{"Y":0,"X":1},{"X":1,"Y":1},{"X":0,"Y":1},{"X":0,"Y":2},{"X":1,"Y":2},{"Y":2,"X":2},{"X":3,"Y":2},{"Y":2,"X":4},{"X":5,"Y":2},{"X":6,"Y":2},{"X":7,"Y":2},{"X":8,"Y":2},{"X":8,"Y":3},{"Y":4,"X":8}],"Health":96,"EliminatedCause":""},{"ID":"gs_wc78BfCt3r9cDHHwpGKm7Ypd","EliminatedCause":"","Health":100,"EliminatedOnTurn":0,"EliminatedBy":"","Body":[{"X":2,"Y":3},{"Y":4,"X":2},{"X":2,"Y":5},{"X":2,"Y":6},{"X":2,"Y":7},{"X":1,"Y":7},{"Y":8,"X":1},{"Y":8,"X":2},{"X":3,"Y":8},{"X":4,"Y":8},{"X":5,"Y":8},{"X":6,"Y":8},{"X":7,"Y":8},{"X":8,"Y":8},{"X":9,"Y":8},{"Y":8,"X":10},{"X":10,"Y":7},{"X":10,"Y":7}]}],"Turn":147,"Width":11,"Food":[{"X":6,"Y":5}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "avoid an awkward corner when i'm shorter going for snack",
			state:       []byte(`{"Width":11,"Height":11,"Snakes":[{"Health":85,"EliminatedCause":"","Body":[{"X":10,"Y":4},{"Y":3,"X":10},{"X":10,"Y":2},{"X":9,"Y":2},{"Y":2,"X":8},{"Y":2,"X":7}],"EliminatedBy":"","ID":"you","EliminatedOnTurn":0},{"ID":"gs_VymmYHMkQVwqDy3BP8Pv6hWR","Health":95,"EliminatedCause":"","Body":[{"X":8,"Y":6},{"Y":6,"X":7},{"X":7,"Y":7},{"Y":8,"X":7},{"X":7,"Y":9},{"Y":10,"X":7},{"Y":10,"X":6},{"X":5,"Y":10}],"EliminatedOnTurn":0,"EliminatedBy":""}],"Turn":52,"Food":[{"Y":6,"X":10}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionLeft},
		},
		// todo: add snack in top left corner. current failure mode is trivial mispath followed by unsafe move
		{
			explanation: "avoid entering a space where i could be easily cut off",
			state:       []byte(`{"Food":[{"Y":7,"X":9},{"X":10,"Y":0}],"Hazards":null,"Snakes":[{"EliminatedBy":"","Body":[{"Y":5,"X":3},{"Y":6,"X":3},{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"Y":9,"X":4},{"X":5,"Y":9},{"X":5,"Y":10},{"Y":10,"X":6},{"Y":9,"X":6},{"Y":9,"X":7},{"X":8,"Y":9},{"Y":9,"X":9},{"X":10,"Y":9},{"X":10,"Y":8},{"X":9,"Y":8},{"X":8,"Y":8},{"X":7,"Y":8}],"EliminatedOnTurn":0,"Health":92,"ID":"you","EliminatedCause":""},{"EliminatedCause":"","Health":98,"ID":"gs_XXp6TM7X8QXRBmXycrTVk7SP","Body":[{"X":1,"Y":7},{"Y":6,"X":1},{"X":0,"Y":6},{"Y":5,"X":0},{"Y":4,"X":0},{"Y":4,"X":1},{"X":2,"Y":4},{"Y":4,"X":3},{"Y":3,"X":3},{"X":3,"Y":2},{"X":3,"Y":1},{"Y":1,"X":2},{"Y":0,"X":2},{"Y":0,"X":3},{"X":4,"Y":0}],"EliminatedOnTurn":0,"EliminatedBy":""}],"Height":11,"Turn":138,"Width":11}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "don't go into tomb of you",
			state:       []byte(`{"Hazards":null,"Width":11,"Food":[{"Y":8,"X":0}],"Turn":90,"Snakes":[{"EliminatedBy":"","Health":81,"ID":"gs_MfYkjTrrjBYMBqf4mxM9xRjM","EliminatedOnTurn":0,"EliminatedCause":"","Body":[{"Y":3,"X":3},{"X":4,"Y":3},{"Y":3,"X":5},{"Y":3,"X":6},{"Y":4,"X":6},{"X":7,"Y":4},{"Y":4,"X":8},{"X":9,"Y":4},{"Y":5,"X":9},{"X":8,"Y":5},{"Y":5,"X":7}]},{"Health":95,"EliminatedBy":"","EliminatedCause":"","EliminatedOnTurn":0,"ID":"you","Body":[{"X":6,"Y":8},{"Y":8,"X":7},{"X":8,"Y":8},{"X":8,"Y":9},{"X":9,"Y":9},{"Y":10,"X":9},{"Y":10,"X":8},{"Y":10,"X":7},{"Y":10,"X":6},{"Y":10,"X":5},{"X":5,"Y":9},{"Y":8,"X":5},{"Y":7,"X":5},{"Y":7,"X":4}]}],"Height":11}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
		{
			explanation: "don't go into tomb of you after chasing tail in pursuit of snak",
			state:       []byte(`{"Snakes":[{"EliminatedBy":"","Body":[{"Y":4,"X":7},{"Y":4,"X":8},{"X":8,"Y":5},{"Y":5,"X":9},{"X":10,"Y":5},{"Y":6,"X":10},{"Y":7,"X":10},{"Y":7,"X":9},{"Y":7,"X":8},{"X":7,"Y":7},{"X":6,"Y":7},{"X":6,"Y":6},{"Y":5,"X":6},{"X":6,"Y":4}],"ID":"you","EliminatedOnTurn":0,"Health":97,"EliminatedCause":""},{"ID":"gs_Y87CWDQHbcxPkrtJStypp4TT","EliminatedOnTurn":0,"Body":[{"Y":5,"X":2},{"X":2,"Y":6},{"X":2,"Y":7},{"X":2,"Y":8},{"X":3,"Y":8},{"Y":8,"X":4},{"X":5,"Y":8},{"X":5,"Y":7},{"Y":7,"X":4},{"Y":6,"X":4},{"X":4,"Y":5},{"Y":4,"X":4},{"Y":4,"X":3},{"X":3,"Y":3},{"X":4,"Y":3},{"Y":2,"X":4}],"EliminatedBy":"","Health":86,"EliminatedCause":""}],"Turn":93,"Height":11,"Hazards":null,"Food":[{"X":7,"Y":6}],"Width":11}`),
			okMoves:     []generator.Direction{generator.DirectionDown, generator.DirectionLeft},
		},
		{
			explanation: "don't smoothbrain into the corner",
			state:       []byte(`{"Hazards":null,"Height":11,"Turn":206,"Width":11,"Snakes":[{"Health":98,"EliminatedCause":"","Body":[{"Y":10,"X":8},{"X":7,"Y":10},{"Y":10,"X":6},{"Y":9,"X":6},{"Y":8,"X":6},{"X":6,"Y":7},{"X":7,"Y":7},{"Y":6,"X":7},{"X":7,"Y":5},{"Y":4,"X":7},{"X":7,"Y":3},{"X":7,"Y":2},{"Y":2,"X":8},{"X":8,"Y":1},{"Y":1,"X":7},{"X":6,"Y":1},{"Y":1,"X":5},{"X":4,"Y":1},{"Y":2,"X":4},{"X":5,"Y":2},{"X":5,"Y":3},{"Y":4,"X":5},{"X":5,"Y":5},{"X":6,"Y":5}],"EliminatedBy":"","ID":"gs_FhYrcM9GcjJyVJ6FXrgBXPxK","EliminatedOnTurn":0},{"EliminatedBy":"","Body":[{"X":4,"Y":10},{"Y":9,"X":4},{"Y":9,"X":5},{"X":5,"Y":8},{"X":4,"Y":8},{"X":4,"Y":7},{"Y":6,"X":4},{"Y":5,"X":4},{"X":4,"Y":4},{"X":4,"Y":3},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":4,"X":2},{"X":3,"Y":4},{"X":3,"Y":5},{"Y":5,"X":2},{"Y":5,"X":1},{"X":0,"Y":5},{"X":0,"Y":6},{"X":0,"Y":7},{"Y":8,"X":0},{"X":1,"Y":8},{"X":2,"Y":8},{"X":2,"Y":9},{"X":2,"Y":10}],"ID":"you","Health":80,"EliminatedCause":"","EliminatedOnTurn":0}],"Food":[{"X":8,"Y":9}]}`),
			okMoves:     []generator.Direction{generator.DirectionLeft},
		},
		{
			explanation: "don't panic, we can still get out",
			state:       []byte(`{"Width":11,"Hazards":null,"Height":11,"Turn":157,"Snakes":[{"Body":[{"X":1,"Y":6},{"X":1,"Y":7},{"X":2,"Y":7},{"X":3,"Y":7},{"X":4,"Y":7},{"X":5,"Y":7},{"Y":6,"X":5},{"X":5,"Y":5},{"Y":5,"X":6},{"X":6,"Y":4},{"X":7,"Y":4},{"Y":3,"X":7},{"Y":3,"X":6},{"X":5,"Y":3},{"X":4,"Y":3},{"X":3,"Y":3},{"X":3,"Y":4},{"X":2,"Y":4},{"Y":4,"X":1},{"Y":5,"X":1},{"X":0,"Y":5},{"X":0,"Y":4},{"X":0,"Y":3},{"Y":2,"X":0},{"X":0,"Y":1},{"X":1,"Y":1}],"EliminatedCause":"","EliminatedOnTurn":0,"ID":"you","Health":89,"EliminatedBy":""},{"EliminatedBy":"","EliminatedCause":"","ID":"gs_pdW6CkXr7xmdvgrGrKJpB6X9","Body":[{"X":3,"Y":10},{"Y":10,"X":2},{"Y":9,"X":2},{"Y":8,"X":2},{"Y":8,"X":3},{"X":4,"Y":8},{"X":5,"Y":8},{"Y":8,"X":6},{"Y":7,"X":6},{"Y":6,"X":6},{"X":7,"Y":6},{"Y":5,"X":7},{"Y":5,"X":8},{"Y":6,"X":8},{"X":8,"Y":7},{"X":8,"Y":8}],"EliminatedOnTurn":0,"Health":91}],"Food":[{"Y":2,"X":2}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "don't chase treat into imminent death",
			state:       []byte(`{"Turn":27,"Hazards":null,"Width":11,"Food":[{"X":9,"Y":5},{"Y":5,"X":7}],"Height":11,"Snakes":[{"EliminatedCause":"","EliminatedBy":"","Health":91,"EliminatedOnTurn":0,"Body":[{"Y":4,"X":7},{"X":6,"Y":4},{"X":6,"Y":5},{"X":5,"Y":5},{"Y":5,"X":4},{"X":4,"Y":4},{"X":4,"Y":3}],"ID":"you"},{"Body":[{"X":9,"Y":6},{"X":8,"Y":6},{"Y":6,"X":7},{"X":7,"Y":7},{"Y":8,"X":7},{"X":7,"Y":9},{"X":7,"Y":10},{"Y":10,"X":7}],"EliminatedBy":"","EliminatedCause":"","ID":"gs_Fw7fv3pKbVS9dCQ7DT7rrkC6","EliminatedOnTurn":0,"Health":100}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight, generator.DirectionDown},
		},
		{
			explanation: "check for longest path even if no snack available",
			state:       []byte(`{"Turn":385,"Height":11,"Snakes":[{"Body":[{"X":8,"Y":3},{"X":9,"Y":3},{"X":9,"Y":4},{"Y":5,"X":9},{"Y":5,"X":10},{"X":10,"Y":6},{"Y":7,"X":10},{"Y":8,"X":10},{"Y":9,"X":10},{"X":9,"Y":9},{"Y":8,"X":9},{"Y":8,"X":8},{"X":7,"Y":8},{"Y":8,"X":6},{"Y":8,"X":5},{"Y":8,"X":4},{"Y":8,"X":3},{"X":2,"Y":8},{"X":2,"Y":7},{"X":2,"Y":6},{"X":3,"Y":6},{"Y":6,"X":4},{"X":5,"Y":6},{"X":6,"Y":6},{"Y":6,"X":7},{"X":8,"Y":6},{"X":8,"Y":5},{"Y":5,"X":7},{"X":6,"Y":5},{"Y":5,"X":5},{"X":4,"Y":5},{"X":3,"Y":5},{"X":2,"Y":5},{"Y":4,"X":2},{"X":2,"Y":3},{"X":2,"Y":2},{"X":3,"Y":2},{"Y":2,"X":4},{"Y":2,"X":5},{"X":6,"Y":2},{"X":6,"Y":3},{"Y":3,"X":7}],"EliminatedBy":"","EliminatedCause":"","ID":"gs_yYVhSD4h8jC7yJMDtw4YjX88","EliminatedOnTurn":0,"Health":83},{"Body":[{"X":1,"Y":8},{"X":0,"Y":8},{"Y":7,"X":0},{"X":0,"Y":6},{"Y":5,"X":0},{"Y":4,"X":0},{"X":0,"Y":3},{"Y":2,"X":0},{"Y":1,"X":0},{"Y":1,"X":1},{"X":2,"Y":1},{"Y":1,"X":3},{"Y":1,"X":4},{"Y":1,"X":5},{"X":6,"Y":1},{"X":7,"Y":1},{"Y":1,"X":8},{"X":9,"Y":1},{"X":9,"Y":2},{"X":10,"Y":2},{"X":10,"Y":1},{"X":10,"Y":0},{"Y":0,"X":9},{"X":9,"Y":0}],"EliminatedBy":"","EliminatedCause":"","ID":"you","Health":100,"EliminatedOnTurn":0}],"Food":[{"Y":4,"X":5},{"Y":4,"X":6},{"Y":4,"X":4}],"Hazards":null,"Width":11}`),
			okMoves:     []generator.Direction{generator.DirectionUp},
		},
		{
			explanation: "don't go into corner, despite being trapped just in case",
			state:       []byte(`{"Turn":121,"Height":11,"Width":11,"Food":[{"X":10,"Y":0}],"Snakes":[{"ID":"you","Body":[{"X":10,"Y":2},{"X":9,"Y":2},{"X":8,"Y":2},{"X":8,"Y":1},{"X":8,"Y":0},{"X":7,"Y":0},{"X":6,"Y":0},{"X":5,"Y":0},{"X":4,"Y":0},{"X":3,"Y":0},{"X":2,"Y":0},{"X":1,"Y":0},{"X":0,"Y":0},{"X":0,"Y":1},{"X":1,"Y":1},{"X":1,"Y":2},{"X":1,"Y":3},{"X":1,"Y":4},{"X":1,"Y":5},{"X":1,"Y":6},{"X":1,"Y":7},{"X":1,"Y":8},{"X":1,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"3066e773-5744-446b-955c-e1cc2f994c52","Body":[{"X":6,"Y":4},{"X":6,"Y":3},{"X":6,"Y":2},{"X":7,"Y":2},{"X":7,"Y":3},{"X":8,"Y":3},{"X":8,"Y":4},{"X":8,"Y":5},{"X":8,"Y":6},{"X":7,"Y":6},{"X":6,"Y":6},{"X":5,"Y":6}],"Health":96,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionUp},
		},
		{
			explanation: "don't get todo logic",
			state:       []byte(`{"Turn":376,"Height":12,"Width":11,"Food":[{"X":5,"Y":0},{"X":3,"Y":1},{"X":9,"Y":2},{"X":9,"Y":9}],"Snakes":[{"ID":"you2","Body":[{"X":7,"Y":4},{"X":6,"Y":4},{"X":5,"Y":4},{"X":5,"Y":5},{"X":5,"Y":6},{"X":5,"Y":7},{"X":5,"Y":8},{"X":4,"Y":8},{"X":4,"Y":7},{"X":4,"Y":6},{"X":4,"Y":5},{"X":3,"Y":5},{"X":3,"Y":6},{"X":2,"Y":6},{"X":2,"Y":7},{"X":2,"Y":8},{"X":1,"Y":8},{"X":1,"Y":9},{"X":1,"Y":10},{"X":2,"Y":10},{"X":3,"Y":10},{"X":4,"Y":10},{"X":5,"Y":10},{"X":6,"Y":10},{"X":7,"Y":10},{"X":8,"Y":10},{"X":8,"Y":9},{"X":8,"Y":8},{"X":7,"Y":8},{"X":6,"Y":8},{"X":6,"Y":7},{"X":6,"Y":6},{"X":7,"Y":6},{"X":8,"Y":6},{"X":8,"Y":5},{"X":8,"Y":4}],"Health":82,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"you","Body":[{"X":1,"Y":2},{"X":1,"Y":3},{"X":2,"Y":3},{"X":3,"Y":3},{"X":3,"Y":2},{"X":2,"Y":2},{"X":2,"Y":1},{"X":1,"Y":1},{"X":1,"Y":0},{"X":0,"Y":0},{"X":0,"Y":1},{"X":0,"Y":2},{"X":0,"Y":3},{"X":0,"Y":4},{"X":0,"Y":5},{"X":0,"Y":6},{"X":0,"Y":7},{"X":1,"Y":7},{"X":1,"Y":6},{"X":1,"Y":5},{"X":2,"Y":5},{"X":2,"Y":4},{"X":3,"Y":4},{"X":4,"Y":4},{"X":4,"Y":3},{"X":5,"Y":3},{"X":6,"Y":3},{"X":6,"Y":2},{"X":6,"Y":1},{"X":6,"Y":0},{"X":7,"Y":0},{"X":7,"Y":1},{"X":8,"Y":1},{"X":8,"Y":0}],"Health":96,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionUp},
		},
		{
			explanation: "definitely don't panic. escape is possible",
			state:       []byte(`{"Hazards":null,"Food":[{"Y":9,"X":4},{"Y":10,"X":10},{"X":3,"Y":7}],"Width":11,"Snakes":[{"ID":"gs_Mryy8QHFPhDq3cxWJJpRFkg4","EliminatedCause":"","EliminatedOnTurn":0,"Health":78,"Body":[{"Y":6,"X":3},{"X":3,"Y":5},{"X":3,"Y":4},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":3,"X":1},{"X":0,"Y":3},{"Y":4,"X":0},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"X":0,"Y":8},{"Y":8,"X":1},{"X":2,"Y":8},{"Y":7,"X":2},{"Y":7,"X":1}],"EliminatedBy":""},{"Body":[{"X":9,"Y":2},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"Y":1,"X":5},{"X":5,"Y":2},{"Y":3,"X":5},{"X":5,"Y":4},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":4,"X":9},{"X":10,"Y":4},{"Y":5,"X":10},{"Y":5,"X":9},{"Y":5,"X":8},{"Y":5,"X":7},{"X":7,"Y":6}],"EliminatedOnTurn":0,"EliminatedCause":"","Health":97,"EliminatedBy":"","ID":"you"}],"Height":11,"Turn":163}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
		{
			explanation: "don't timeout for pete's sake",
			state:       []byte(`{"Turn":394,"Height":12,"Width":11,"Food":[{"X":3,"Y":6},{"X":4,"Y":10},{"X":2,"Y":6},{"X":3,"Y":4},{"X":9,"Y":8},{"X":4,"Y":8},{"X":0,"Y":5},{"X":6,"Y":0}],"Snakes":[{"ID":"you2","Body":[{"X":4,"Y":3},{"X":4,"Y":2},{"X":5,"Y":2},{"X":5,"Y":3},{"X":6,"Y":3},{"X":7,"Y":3},{"X":7,"Y":2},{"X":6,"Y":2},{"X":6,"Y":1},{"X":7,"Y":1},{"X":7,"Y":0},{"X":8,"Y":0},{"X":8,"Y":1},{"X":8,"Y":2},{"X":9,"Y":2},{"X":9,"Y":1},{"X":10,"Y":1},{"X":10,"Y":2},{"X":10,"Y":3},{"X":9,"Y":3},{"X":8,"Y":3},{"X":8,"Y":4},{"X":7,"Y":4},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"X":5,"Y":4},{"X":4,"Y":4},{"X":4,"Y":5},{"X":3,"Y":5},{"X":2,"Y":5},{"X":1,"Y":5},{"X":1,"Y":4},{"X":2,"Y":4},{"X":2,"Y":3},{"X":2,"Y":2},{"X":3,"Y":2},{"X":3,"Y":3}],"Health":65,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"you","Body":[{"X":1,"Y":6},{"X":1,"Y":7},{"X":0,"Y":7},{"X":0,"Y":8},{"X":0,"Y":9},{"X":1,"Y":9},{"X":1,"Y":10},{"X":1,"Y":11},{"X":2,"Y":11},{"X":3,"Y":11},{"X":4,"Y":11},{"X":5,"Y":11},{"X":6,"Y":11},{"X":7,"Y":11},{"X":8,"Y":11},{"X":9,"Y":11},{"X":10,"Y":11},{"X":10,"Y":10},{"X":10,"Y":9},{"X":10,"Y":8},{"X":10,"Y":7},{"X":10,"Y":6},{"X":9,"Y":6},{"X":9,"Y":7},{"X":8,"Y":7},{"X":8,"Y":6},{"X":7,"Y":6},{"X":6,"Y":6},{"X":5,"Y":6},{"X":4,"Y":6},{"X":4,"Y":7},{"X":3,"Y":7},{"X":3,"Y":8},{"X":3,"Y":9},{"X":2,"Y":9},{"X":2,"Y":8},{"X":2,"Y":8}],"Health":100,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionLeft},
		},
		{
			explanation: "plenty of space, but only if you go the right way",
			state:       []byte(`{"Turn":215,"Height":12,"Width":11,"Food":[{"X":10,"Y":8}],"Snakes":[{"ID":"c8a230f0-9e4d-442d-8f60-c63d80ba0ab9","Body":[{"X":3,"Y":5},{"X":2,"Y":5},{"X":2,"Y":6},{"X":2,"Y":7},{"X":2,"Y":8},{"X":2,"Y":9},{"X":2,"Y":10},{"X":2,"Y":11},{"X":3,"Y":11},{"X":4,"Y":11},{"X":5,"Y":11},{"X":6,"Y":11},{"X":6,"Y":10},{"X":6,"Y":9},{"X":6,"Y":8},{"X":6,"Y":7},{"X":6,"Y":6},{"X":7,"Y":6},{"X":8,"Y":6},{"X":9,"Y":6},{"X":9,"Y":5},{"X":8,"Y":5},{"X":7,"Y":5},{"X":6,"Y":5},{"X":5,"Y":5},{"X":5,"Y":4},{"X":4,"Y":4},{"X":4,"Y":5}],"Health":93,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"you","Body":[{"X":0,"Y":4},{"X":1,"Y":4},{"X":2,"Y":4},{"X":2,"Y":3},{"X":2,"Y":2},{"X":2,"Y":1},{"X":3,"Y":1},{"X":3,"Y":2},{"X":3,"Y":3},{"X":4,"Y":3},{"X":4,"Y":2},{"X":4,"Y":1},{"X":4,"Y":0},{"X":5,"Y":0},{"X":5,"Y":1},{"X":6,"Y":1},{"X":6,"Y":2},{"X":7,"Y":2},{"X":7,"Y":1},{"X":8,"Y":1}],"Health":94,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
		{
			explanation: "it's a trap",
			state:       []byte(`{"Width":11,"Height":11,"Food":[{"Y":10,"X":9},{"Y":5,"X":5}],"Turn":198,"Snakes":[{"ID":"gs_RBHx6yBmt3FvYxtr8CxCgCXb","EliminatedCause":"","Health":51,"EliminatedBy":"","EliminatedOnTurn":0,"Body":[{"Y":9,"X":5},{"X":4,"Y":9},{"X":3,"Y":9},{"X":3,"Y":8},{"Y":8,"X":4},{"Y":8,"X":5},{"Y":8,"X":6},{"Y":8,"X":7},{"X":8,"Y":8},{"X":9,"Y":8},{"Y":8,"X":10},{"X":10,"Y":7},{"X":9,"Y":7},{"Y":7,"X":8},{"Y":7,"X":7},{"Y":7,"X":6},{"X":5,"Y":7}]},{"EliminatedOnTurn":0,"Body":[{"Y":3,"X":5},{"X":4,"Y":3},{"Y":4,"X":4},{"Y":4,"X":3},{"X":3,"Y":3},{"X":2,"Y":3},{"Y":2,"X":2},{"Y":1,"X":2},{"X":1,"Y":1},{"Y":1,"X":0},{"Y":2,"X":0},{"Y":2,"X":1},{"Y":3,"X":1},{"Y":3,"X":0},{"X":0,"Y":4},{"X":1,"Y":4},{"Y":4,"X":2},{"X":2,"Y":4}],"EliminatedCause":"","Health":100,"ID":"you","EliminatedBy":""},{"EliminatedOnTurn":0,"EliminatedBy":"","Health":97,"EliminatedCause":"","ID":"gs_Pq4RKFjCVJhdJYVdHSwRSSrH","Body":[{"Y":6,"X":6},{"X":7,"Y":6},{"Y":5,"X":7},{"X":6,"Y":5},{"Y":4,"X":6},{"X":6,"Y":3},{"X":6,"Y":2},{"Y":1,"X":6},{"X":5,"Y":1},{"Y":1,"X":4},{"X":4,"Y":0},{"Y":0,"X":5},{"X":6,"Y":0},{"X":7,"Y":0},{"X":8,"Y":0},{"X":8,"Y":1},{"Y":2,"X":8},{"X":8,"Y":3},{"X":8,"Y":4},{"X":8,"Y":5},{"Y":5,"X":9},{"X":9,"Y":4},{"Y":3,"X":9},{"X":9,"Y":2},{"Y":2,"X":10}]}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionDown},
		},
		{
			explanation: "is your opponent hungry",
			state:       []byte(`{"Snakes":[{"EliminatedBy":"","Body":[{"Y":1,"X":0},{"Y":2,"X":0},{"Y":3,"X":0},{"Y":4,"X":0},{"Y":4,"X":1},{"X":2,"Y":4},{"X":2,"Y":5},{"X":1,"Y":5},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"Y":7,"X":1},{"X":1,"Y":8},{"Y":8,"X":2}],"Health":92,"EliminatedOnTurn":0,"ID":"gs_MdfdWg8J4tWG9vk9DFfrWXXK","EliminatedCause":""},{"Body":[{"Y":3,"X":6},{"X":5,"Y":3},{"X":5,"Y":2},{"X":6,"Y":2},{"Y":2,"X":7},{"Y":1,"X":7},{"X":6,"Y":1},{"Y":1,"X":5},{"X":5,"Y":0},{"X":4,"Y":0},{"X":3,"Y":0}],"Health":88,"EliminatedOnTurn":0,"EliminatedBy":"","ID":"gs_Y9JHCFXjQHk7WD9TSwYMSJS3","EliminatedCause":""},{"EliminatedBy":"","EliminatedCause":"","EliminatedOnTurn":0,"ID":"you","Health":95,"Body":[{"X":0,"Y":9},{"Y":9,"X":1},{"Y":9,"X":2},{"X":3,"Y":9},{"Y":9,"X":4},{"X":4,"Y":8},{"Y":8,"X":5},{"X":6,"Y":8},{"Y":8,"X":7},{"X":7,"Y":7},{"Y":6,"X":7},{"X":7,"Y":5}]}],"Food":[{"Y":0,"X":0},{"Y":1,"X":8}],"Turn":105,"Height":11,"Width":11,"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionUp},
		},
		{
			explanation: "smooth check",
			state:       []byte(`{"Turn":161,"Height":11,"Width":11,"Food":[{"X":3,"Y":9},{"X":3,"Y":1},{"X":1,"Y":10},{"X":2,"Y":8}],"Snakes":[{"ID":"0d3db8a8-eec2-40a8-9944-782d245a14c1","Body":[{"X":0,"Y":8},{"X":0,"Y":7},{"X":0,"Y":6},{"X":1,"Y":6},{"X":1,"Y":7},{"X":2,"Y":7},{"X":3,"Y":7},{"X":4,"Y":7},{"X":5,"Y":7},{"X":6,"Y":7},{"X":7,"Y":7},{"X":8,"Y":7},{"X":9,"Y":7},{"X":9,"Y":6},{"X":10,"Y":6},{"X":10,"Y":7},{"X":10,"Y":8}],"Health":96,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"you","Body":[{"X":1,"Y":1},{"X":1,"Y":2},{"X":0,"Y":2},{"X":0,"Y":3},{"X":0,"Y":4},{"X":0,"Y":5},{"X":1,"Y":5},{"X":2,"Y":5},{"X":3,"Y":5},{"X":4,"Y":5},{"X":5,"Y":5},{"X":5,"Y":4},{"X":5,"Y":3},{"X":6,"Y":3}],"Health":98,"EliminatedCause":"","EliminatedOnTurn":0,"EliminatedBy":""}],"Hazards":null}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "obvious doom",
			state:       []byte(`{"Height":11,"Hazards":null,"Food":[{"X":2,"Y":8}],"Width":11,"Turn":115,"Snakes":[{"ID":"gs_crKh3x7CkwKjr9SJR7CmdR7c","Health":100,"Body":[{"Y":1,"X":0},{"Y":2,"X":0},{"Y":3,"X":0},{"X":0,"Y":4},{"Y":5,"X":0},{"Y":5,"X":1},{"Y":5,"X":2},{"X":2,"Y":6},{"X":3,"Y":6},{"X":3,"Y":6}],"EliminatedOnTurn":0,"EliminatedBy":"","EliminatedCause":""},{"Health":95,"ID":"gs_Hr4bGmRYBBCFSv4mK9F9vHK8","EliminatedBy":"","Body":[{"X":3,"Y":4},{"Y":4,"X":4},{"Y":4,"X":5},{"X":6,"Y":4},{"Y":3,"X":6},{"Y":3,"X":5},{"X":4,"Y":3},{"X":3,"Y":3},{"X":3,"Y":2},{"X":4,"Y":2},{"X":5,"Y":2},{"Y":2,"X":6},{"Y":2,"X":7},{"Y":1,"X":7},{"X":6,"Y":1}],"EliminatedCause":"","EliminatedOnTurn":0},{"EliminatedBy":"","EliminatedCause":"","Body":[{"X":5,"Y":6},{"Y":6,"X":6},{"Y":6,"X":7},{"Y":5,"X":7},{"X":8,"Y":5},{"X":8,"Y":4},{"X":8,"Y":3},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0}],"Health":71,"EliminatedOnTurn":0,"ID":"you"}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
		{
			explanation: "don't panic",
			state:       []byte(`{"Food":[{"Y":8,"X":3},{"X":0,"Y":6},{"X":9,"Y":2}],"Turn":170,"Height":11,"Hazards":null,"Width":11,"Snakes":[{"ID":"you","EliminatedCause":"","Health":95,"Body":[{"X":8,"Y":10},{"X":7,"Y":10},{"Y":9,"X":7},{"X":7,"Y":8},{"Y":7,"X":7},{"X":7,"Y":6},{"Y":6,"X":6},{"Y":6,"X":5},{"Y":6,"X":4},{"X":4,"Y":7},{"Y":8,"X":4},{"X":4,"Y":9},{"Y":10,"X":4},{"X":3,"Y":10}],"EliminatedOnTurn":0,"EliminatedBy":""},{"ID":"gs_9tYCpg9rmbCdv9QrjRQ74md4","EliminatedBy":"","EliminatedOnTurn":0,"EliminatedCause":"","Body":[{"Y":7,"X":1},{"Y":6,"X":1},{"Y":6,"X":2},{"X":3,"Y":6},{"X":3,"Y":5},{"Y":5,"X":2},{"X":1,"Y":5},{"X":1,"Y":4},{"X":0,"Y":4},{"X":0,"Y":3},{"Y":3,"X":1},{"X":1,"Y":2},{"Y":1,"X":1},{"Y":1,"X":2}],"Health":87},{"ID":"gs_vDjhQKPPKWQRRwVW8jpmTgG7","EliminatedBy":"","Body":[{"X":10,"Y":4},{"Y":5,"X":10},{"Y":6,"X":10},{"X":9,"Y":6},{"X":9,"Y":5},{"X":8,"Y":5},{"Y":5,"X":7},{"X":6,"Y":5},{"Y":5,"X":5},{"Y":4,"X":5},{"Y":3,"X":5},{"Y":2,"X":5}],"Health":92,"EliminatedCause":"","EliminatedOnTurn":0}]}`),
			okMoves:     []generator.Direction{generator.DirectionRight},
		},
	}
)

func TestMove(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	youID := "you"

	for _, test := range tests {

		if test.explanation != "avoid an awkward corner when i'm shorter going for snack" {
			continue
		}

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
		if reason == "yeet todo logic" {
			generator.PrintMap(s)
			t.Log("got todo logic")
			t.Fail()
			continue
		}

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

		generator.PrintMap(s)

		t.Logf("%s FAILED: got %s because %s. ok moves: %+v", test.explanation, move.String(), reason, test.okMoves)
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

// should have chased tail to freedom
// https://play.battlesnake.com/g/bf7ae65f-fc87-425f-888d-32f8d4de1e80/
// turn 163

// should try and kill in corner, plus needs to know about length after eating snack
// https://play.battlesnake.com/g/bfc75d0a-c92d-4d10-a75c-e73d2fdf8bd4
// turn 90

// should adjust weightings according to heads of snakes that are longer than you
// https://play.battlesnake.com/g/d85697ea-d661-4e3a-9bf6-66264a20d4be/
// turn 25
