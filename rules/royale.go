package rules

import (
	"errors"
	"fmt"
	"math/rand"
)

type RoyaleRuleset struct {
	StandardRuleset

	Seed int64

	ShrinkEveryNTurns int32
}

func (r *RoyaleRuleset) Name() string { return "royale" }

func (r *RoyaleRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	if r.StandardRuleset.HazardDamagePerTurn < 1 {
		return nil, errors.New("royale damage per turn must be greater than zero")
	}

	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	// Royale's only job is now to populate the hazards for next turn - StandardRuleset takes care of applying hazard damage.
	err = r.populateHazards(nextBoardState, prevState.Turn+1)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func (r *RoyaleRuleset) ApplySingleMove(prevState *BoardState, move SnakeMoveIndex) (*BoardState, error) {

	return nil, fmt.Errorf("not implemented for royale")

}

func (r *RoyaleRuleset) populateHazards(b *BoardState, turn int32) error {
	b.Hazards = []Point{}

	if r.ShrinkEveryNTurns < 1 {
		return errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < r.ShrinkEveryNTurns {
		return nil
	}

	randGenerator := rand.New(rand.NewSource(r.Seed))

	numShrinks := turn / r.ShrinkEveryNTurns
	minX, maxX := byte(0), b.Width-1
	minY, maxY := byte(0), b.Height-1
	for i := int32(0); i < numShrinks; i++ {
		switch randGenerator.Intn(4) {
		case 0:
			minX += 1
		case 1:
			maxX -= 1
		case 2:
			minY += 1
		case 3:
			maxY -= 1
		}
	}

	for x := byte(0); x < b.Width; x++ {
		for y := byte(0); y < b.Height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				b.Hazards = append(b.Hazards, Point{x, y})
			}
		}
	}

	return nil
}
