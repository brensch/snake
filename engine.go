package main

import (
	"fmt"

	"github.com/brensch/snake/rules"
)

type EngineCoord struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type EngineSnake struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Health  int32         `json:"health"`
	Body    []EngineCoord `json:"body"`
	Latency string        `json:"latency"`
	Head    EngineCoord   `json:"head"`
	Length  int32         `json:"length"`
	Shout   string        `json:"shout"`
	Squad   string        `json:"squad"`
}

func (s EngineSnake) ToSnake() rules.Snake {
	// if s.Shout != "" {
	// 	fmt.Printf("%s said %s\n", s.Name, s.Shout)
	// }
	return rules.Snake{
		ID:              s.ID,
		Body:            CoordsToPoints(s.Body),
		Health:          s.Health,
		EliminatedCause: "",
		EliminatedBy:    "",
	}
}

func CoordsToPoints(coords []EngineCoord) []rules.Point {
	var points []rules.Point
	for _, coord := range coords {
		points = append(points, rules.Point{X: coord.X, Y: coord.Y})
	}
	return points
}

type EngineBoard struct {
	Height  int32         `json:"height"`
	Width   int32         `json:"width"`
	Food    []EngineCoord `json:"food"`
	Hazards []EngineCoord `json:"hazards"`
	Snakes  []EngineSnake `json:"snakes"`
}

type EngineRuleSetSettings struct {
	HazardDamagePerTurn int32          `json:"hazardDamagePerTurn"`
	FoodSpawnChance     int32          `json:"foodSpawnChance"`
	MinimumFood         int32          `json:"minimumFood"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"`
}

type RoyaleSettings struct {
	ShrinkEveryNTurns int32 `json:"shrinkEveryNTurns"`
}

type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

type EngineRuleSet struct {
	Name     string                `json:"name"`
	Version  string                `json:"version"`
	Settings EngineRuleSetSettings `json:"settings"`
}

type EngineGame struct {
	ID      string        `json:"id"`
	Timeout int32         `json:"timeout"`
	Ruleset EngineRuleSet `json:"ruleset"`
}

type EngineRequest struct {
	Game  EngineGame  `json:"game"`
	Turn  int32       `json:"turn"`
	Board EngineBoard `json:"board"`
	You   EngineSnake `json:"you"`
}

func (r EngineRequest) ToState() (*rules.BoardState, rules.Ruleset, rules.Snake) {
	you := r.You.ToSnake()
	var snakes []rules.Snake

	for _, snake := range r.Board.Snakes {
		snakes = append(snakes, snake.ToSnake())
	}

	// fmt.Printf("%+v\n", r.Game.Ruleset)

	var ruleset rules.Ruleset

	baseRules := &rules.StandardRuleset{
		FoodSpawnChance:     r.Game.Ruleset.Settings.FoodSpawnChance,
		MinimumFood:         r.Game.Ruleset.Settings.MinimumFood,
		HazardDamagePerTurn: r.Game.Ruleset.Settings.HazardDamagePerTurn,
	}

	switch r.Game.Ruleset.Name {
	case "standard":
		ruleset = baseRules
	case "solo":
		ruleset = &rules.SoloRuleset{
			StandardRuleset: *baseRules,
		}
	case "royale":
		ruleset = &rules.RoyaleRuleset{
			StandardRuleset:   *baseRules,
			ShrinkEveryNTurns: r.Game.Ruleset.Settings.RoyaleSettings.ShrinkEveryNTurns,

			// this is unfortunate. game server doesn't expose seed.
			// still should use the royale ruleset to get accurate health calcs
			Seed: 6969,
		}
	case "squad":
		ruleset = &rules.SquadRuleset{
			StandardRuleset: *baseRules,
		}

	}

	// r.Game.Ruleset.Name

	state := &rules.BoardState{
		Turn:    r.Turn,
		Height:  r.Board.Height,
		Width:   r.Board.Width,
		Food:    CoordsToPoints(r.Board.Food),
		Snakes:  snakes,
		Hazards: CoordsToPoints(r.Board.Hazards),
	}

	return state, ruleset, you

}

type EngineResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

type PingResponse struct {
	ApiVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
	Version    string `json:"version"`
}

func GetSnakeFromState(s *rules.BoardState, id string) (rules.Snake, error) {
	for _, snake := range s.Snakes {
		if snake.ID == id {
			return snake, nil
		}
	}

	return rules.Snake{}, fmt.Errorf("snake not found")
}
