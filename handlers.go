package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := info()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	var req EngineRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}

	start(req)
}

func HandleMove(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
		defer cancel()
		start := time.Now()
		var req EngineRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}

		res := s.move(ctx, req)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Printf("ERROR: Failed to encode move response, %s", err)
			return
		}
		fmt.Println("total time", time.Since(start).Milliseconds())
	}
}

func HandleEnd(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var state EngineRequest
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode end json, %s", err)
			return
		}

		s.end(state)

		// Nothing to respond with here
	}
}
