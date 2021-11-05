package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeft(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: DirectionLeft},
		{ID: "bottomRight", Move: DirectionLeft},
		{ID: "topLeft", Move: DirectionLeft},
		{ID: "topRight", Move: DirectionLeft},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{10, 0}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{9, 0}}},
		{ID: "topLeft", Health: 10, Body: []Point{{10, 10}}},
		{ID: "topRight", Health: 10, Body: []Point{{9, 10}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
	}
}

func TestRight(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: DirectionRight},
		{ID: "bottomRight", Move: DirectionRight},
		{ID: "topLeft", Move: DirectionRight},
		{ID: "topRight", Move: DirectionRight},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{1, 0}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{0, 0}}},
		{ID: "topLeft", Health: 10, Body: []Point{{1, 10}}},
		{ID: "topRight", Health: 10, Body: []Point{{0, 10}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
	}
}

func TestUp(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: DirectionUp},
		{ID: "bottomRight", Move: DirectionUp},
		{ID: "topLeft", Move: DirectionUp},
		{ID: "topRight", Move: DirectionUp},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{0, 1}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{10, 1}}},
		{ID: "topLeft", Health: 10, Body: []Point{{0, 0}}},
		{ID: "topRight", Health: 10, Body: []Point{{10, 0}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

func TestDown(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: DirectionDown},
		{ID: "bottomRight", Move: DirectionDown},
		{ID: "topLeft", Move: DirectionDown},
		{ID: "topRight", Move: DirectionDown},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{0, 10}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{10, 10}}},
		{ID: "topLeft", Health: 10, Body: []Point{{0, 9}}},
		{ID: "topRight", Health: 10, Body: []Point{{10, 9}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

func TestEdgeCrossingCollision(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "right", Health: 10, Body: []Point{{0, 5}}},
			{ID: "rightEdge", Health: 10, Body: []Point{
				{10, 1},
				{10, 2},
				{10, 3},
				{10, 4},
				{10, 5},
				{10, 6},
			}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "right", Move: DirectionLeft},
		{ID: "rightEdge", Move: DirectionDown},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "right", Health: 0, Body: []Point{{10, 5}}, EliminatedCause: EliminatedByCollision, EliminatedBy: "rightEdge"},
		{ID: "rightEdge", Health: 10, Body: []Point{
			{10, 0},
			{10, 1},
			{10, 2},
			{10, 3},
			{10, 4},
			{10, 5},
		}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

func TestEdgeCrossingEating(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "left", Health: 10, Body: []Point{{0, 5}, {1, 5}}},
			{ID: "other", Health: 10, Body: []Point{{5, 5}}},
		},
		Food: []Point{
			{10, 5},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "left", Move: DirectionLeft},
		{ID: "other", Move: DirectionLeft},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "left", Health: 100, Body: []Point{{10, 5}, {0, 5}, {0, 5}}},
		{ID: "other", Health: 9, Body: []Point{{4, 5}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].Health, snake.Health, snake.ID)

	}
}
