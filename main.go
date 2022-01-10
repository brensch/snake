package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	// for {
	// 	stress()
	// }
	heuristicCatalog := make(map[string]map[uint64]float64)
	s := &server{
		heuristicCatalog: heuristicCatalog,
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove(s))
	http.HandleFunc("/end", HandleEnd)

	log.Info("welcome to snake")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func stress() {

	state := []byte(`{"Hazards":null,"Food":[{"Y":9,"X":4},{"Y":10,"X":10},{"X":3,"Y":7}],"Width":11,"Snakes":[{"ID":"gs_Mryy8QHFPhDq3cxWJJpRFkg4","EliminatedCause":"","EliminatedOnTurn":0,"Health":78,"Body":[{"Y":6,"X":3},{"X":3,"Y":5},{"X":3,"Y":4},{"Y":3,"X":3},{"X":2,"Y":3},{"Y":3,"X":1},{"X":0,"Y":3},{"Y":4,"X":0},{"Y":5,"X":0},{"Y":6,"X":0},{"Y":7,"X":0},{"X":0,"Y":8},{"Y":8,"X":1},{"X":2,"Y":8},{"Y":7,"X":2},{"Y":7,"X":1}],"EliminatedBy":""},{"Body":[{"X":9,"Y":2},{"Y":2,"X":8},{"Y":1,"X":8},{"X":8,"Y":0},{"Y":0,"X":7},{"Y":1,"X":7},{"Y":1,"X":6},{"Y":1,"X":5},{"X":5,"Y":2},{"Y":3,"X":5},{"X":5,"Y":4},{"X":5,"Y":5},{"X":6,"Y":5},{"X":6,"Y":4},{"Y":4,"X":7},{"Y":4,"X":8},{"Y":4,"X":9},{"X":10,"Y":4},{"Y":5,"X":10},{"Y":5,"X":9},{"Y":5,"X":8},{"Y":5,"X":7},{"X":7,"Y":6}],"EliminatedOnTurn":0,"EliminatedCause":"","Health":97,"EliminatedBy":"","ID":"you"}],"Height":11,"Turn":163}`)

	var s *rules.BoardState
	err := json.Unmarshal(state, &s)
	if err != nil {
		fmt.Print("err", err.Error())
	}

	ruleset := &rules.StandardRuleset{
		FoodSpawnChance: 0,
		MinimumFood:     0,
	}
	countCHAN := make(chan int)
	var sendWG, recWG sync.WaitGroup

	totalCount := 0
	recWG.Add(1)
	go func() {
		defer recWG.Done()
		for count := range countCHAN {
			totalCount += count
		}

	}()

	start := time.Now()
	for i := 0; i < 4; i++ {
		sendWG.Add(1)
		go func() {
			defer sendWG.Done()
			counter := 0
			for {
				if time.Since(start) > 500*time.Millisecond {
					break
				}

				generator.SafeMoves(s, ruleset, "you")
				counter++
			}

			countCHAN <- counter
		}()
	}

	sendWG.Wait()
	close(countCHAN)
	recWG.Wait()

	fmt.Println(totalCount)

}

func info() PingResponse {
	log.Info("received a ping")
	return PingResponse{
		ApiVersion: "1",
		Author:     "brend",
		Color:      "#118645",
		// Color:      "#ff8645",
		Head:    "replit-mark",
		Tail:    "replit-notmark",
		Version: "0.2",
	}
}

func start(req EngineRequest) {
	state, ruleset, you := req.ToState()
	var names []string
	for _, snake := range req.Board.Snakes {
		names = append(names, snake.Name)
	}

	log.WithFields(log.Fields{
		"game":    req.Game.ID,
		"you":     you.ID,
		"ruleset": ruleset.Name(),
		"players": names,
		"action":  "start",
		"state":   state,
	}).Info("started game")

}

func end(req EngineRequest) {
	state, ruleset, you := req.ToState()
	victoryStatus := "LOST"
	if len(state.Snakes) == 0 {
		victoryStatus = "DRAW"
	}
	for _, snake := range state.Snakes {
		if snake.ID == you.ID {
			victoryStatus = "WON"
		}
	}

	var names []string
	for _, snake := range req.Board.Snakes {
		names = append(names, snake.Name)
	}

	entry := log.WithFields(log.Fields{
		"game":    req.Game.ID,
		"action":  "end",
		"result":  victoryStatus,
		"state":   state,
		"ruleset": ruleset.Name(),
		"players": names,
	})

	if victoryStatus == "LOST" {
		entry.Error("lost game")
		return
	}

	entry.Warning("won game")

}

func (s *server) move(ctx context.Context, req EngineRequest) EngineResponse {

	state, ruleset, you := req.ToState()

	move, reason := s.Move(ctx, state, ruleset, you, req.Turn, req.Game.ID)

	return EngineResponse{
		Move:  move.String(),
		Shout: reason,
	}

}
