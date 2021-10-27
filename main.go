package main

import (
	"context"
	"net/http"
	"os"

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

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	log.Info("welcome to snake")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func info() PingResponse {
	log.Info("received a ping")
	return PingResponse{
		ApiVersion: "1",
		Author:     "brend",
		Color:      "#ff8645",
		Head:       "replit-mark",
		Tail:       "replit-notmark",
		Version:    "0.2",
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

func move(ctx context.Context, req EngineRequest) EngineResponse {

	state, ruleset, you := req.ToState()

	move, reason := Move(ctx, state, ruleset, you, req.Turn, req.Game.ID)

	return EngineResponse{
		Move:  move.String(),
		Shout: reason,
	}

}
