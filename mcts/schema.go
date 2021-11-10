package mcts

import (
	"math/rand"
	"time"
)

const (
	youID = "you" // setting this statically on ingest to reduce lookup time
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	r1         = rand.New(randSource)
)

// type Tree struct {
// 	Nodes []Node
// }
