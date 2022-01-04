package rules

type RulesetError string

func (err RulesetError) Error() string { return string(err) }

const (
	MoveUp    = "up"
	MoveDown  = "down"
	MoveRight = "right"
	MoveLeft  = "left"

	BoardSizeSmall  = 7
	BoardSizeMedium = 11
	BoardSizeLarge  = 19

	SnakeMaxHealth = 100
	SnakeStartSize = 3

	// bvanvugt - TODO: Just return formatted strings instead of codes?
	NotEliminated                   = ""
	EliminatedByCollision           = "snake-collision"
	EliminatedBySelfCollision       = "snake-self-collision"
	EliminatedByOutOfHealth         = "out-of-health"
	EliminatedByHeadToHeadCollision = "head-collision"
	EliminatedByOutOfBounds         = "wall-collision"

	// TODO - Error consts
	ErrorTooManySnakes   = RulesetError("too many snakes for fixed start positions")
	ErrorNoRoomForSnake  = RulesetError("not enough space to place snake")
	ErrorNoRoomForFood   = RulesetError("not enough space to place food")
	ErrorNoMoveFound     = RulesetError("move not provided for snake")
	ErrorZeroLengthSnake = RulesetError("snake is length zero")
)

type Point struct {
	X byte
	Y byte
}

type Direction uint8

func (d Direction) String() string {
	return [...]string{"left", "right", "up", "down", "unknown"}[d]
}

const (
	// Not starting with Unknown even though that would help with empty moves.
	// An empty move now defaults to left (ie 0). This is fine since only meant
	// for internal consumption and i won't pass empty moves.
	DirectionLeft Direction = iota
	DirectionRight
	DirectionUp
	DirectionDown
	DirectionUnknown

	directionCount = 4
)

type Snake struct {
	ID               string
	Body             []Point
	Health           byte
	EliminatedCause  string
	EliminatedOnTurn int32
	EliminatedBy     string
}

type SnakeMove struct {
	ID   string
	Move Direction
}

type SnakeMoveIndex struct {
	Index int
	Move  Direction
}

type Ruleset interface {
	Name() string
	ModifyInitialBoardState(initialState *BoardState) (*BoardState, error)
	CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error)
	ApplySingleMove(prevState *BoardState, move SnakeMoveIndex) (*BoardState, error)
	IsGameOver(state *BoardState) (bool, error)
}
