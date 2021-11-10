package mcts

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

const (
	ucbC = .5
)

type Node struct {
	State   *rules.BoardState
	MoveSet []rules.SnakeMove

	UntriedMovesets [][]rules.SnakeMove
	Children        []*Node
	Parent          *Node

	Visits         uint64
	OutcomeScore   float64
	SelectionScore float64
}

func NewNode(parent *Node, moveSet []rules.SnakeMove, state *rules.BoardState) *Node {
	moveSets := generator.AllMoveSetsForStateRaw(state)

	return &Node{
		State:           state,
		UntriedMovesets: moveSets,
		Parent:          parent,
		MoveSet:         moveSet,
	}

}

func (n *Node) SelectNextChildToExplore() *Node {
	sort.Sort(bySelectionScore(n.Children))
	return n.Children[0]
}

func (n *Node) MakeRandomUntriedMove() *Node {
	// Select a random move we haven't tried.
	var i int = rand.Intn(len(n.UntriedMovesets))
	var moveSet []rules.SnakeMove = n.UntriedMovesets[i]

	// Remove it from the untried moves.
	n.UntriedMovesets = append(n.UntriedMovesets[:i], n.UntriedMovesets[i+1:]...)

	// // Clone the node's state so we don't alter it.
	// var newState GameState = n.state.Clone()
	// newState.MakeMove(move)

	ruleset := &rules.StandardRuleset{
		FoodSpawnChance: 0,
		MinimumFood:     0,
	}

	nextState, err := ruleset.CreateNextBoardState(n.State, moveSet)
	if err != nil {
		fmt.Println("this shouldn't have happened", err)
		panic(err)
	}

	// Build more of the tree.
	var child *Node = NewNode(n, moveSet, nextState)
	n.Children = append(n.Children, child) // End of children list are the children with lowest selection scores (e.g. no visits).

	// Return a game state that can be used for simulations.
	return child
}

// func (n *Node) Search(depth, maxDepth int) {
// 	fmt.Println("search", depth)

// 	n.HeuristicScore()

// 	if depth >= maxDepth {
// 		return
// 	}

// 	ruleset := &rules.StandardRuleset{
// 		FoodSpawnChance: 0,
// 		MinimumFood:     0,
// 	}

// 	if len(n.Children) == 0 {
// 		n.PopulateChildMoveSets()
// 	}

// 	// don't have policy, just pick at random like a true mcts
// 	childLocation := r1.Intn(len(n.Children))

// 	n.GenerateStateAndScoreChild(ruleset, childLocation)
// 	if n.Score <= -100 {
// 		return
// 	}

// 	// if n.Parent != nil {

// 	// }

// 	n.Children[childLocation].Search(depth+1, maxDepth)

// }

// Generate children without calculating their state
// func (n *Node) PopulateChildMoveSets() {
// 	if n.State == nil {
// 		fmt.Println("got an empty child gen request....")
// 		return
// 	}
// 	potentialMoveSets := generator.AllMoveSetsForStateRaw(n.State)

// 	for _, moveSet := range potentialMoveSets {
// 		// fmt.Println(moveSet)
// 		n.Children = append(n.Children, &Node{
// 			MoveSet: moveSet,
// 			Parent:  n,
// 		})
// 	}
// }

// func (n *Node) GenerateStateAndScoreChild(ruleSet rules.Ruleset, childLocation int) error {
// 	if childLocation >= len(n.Children) {
// 		fmt.Println("got invalid childlocation", childLocation)
// 		return fmt.Errorf("invalid childlocation")
// 	}

// 	child := n.Children[childLocation]
// 	// fmt.Println(child.MoveSet)

// 	childState, err := ruleSet.CreateNextBoardState(n.State, child.MoveSet)
// 	if err != nil {
// 		return err
// 	}

// 	child.State = childState
// 	score := child.HeuristicScore()
// 	child.Score = score
// 	// if child.Score > 0 {
// 	// 	fmt.Println("child", child.Score)
// 	// }
// 	// fmt.Println(child.Score)
// 	parent := child.Parent
// 	for parent != nil {
// 		parent.Score = score
// 		// fmt.Println("parent", parent.Score)
// 		parent = parent.Parent
// 	}

// 	return nil
// }

func (n *Node) UpdateOutcomeScore(outcome float64) {
	if n == nil {
		return
	}

	n.OutcomeScore += outcome
	n.Visits++
	n.Parent.UpdateOutcomeScore(outcome)
	n.ComputeSelectionScore()
}

func (n *Node) ComputeSelectionScore() {
	n.SelectionScore = n.OutcomeScore/float64(n.Visits) + ucbC*math.Sqrt(2*math.Log(float64(n.Parent.getVisits()))/float64(n.Visits))

}

func (n *Node) getVisits() uint64 {
	if n == nil {
		return 0
	}
	return n.Visits
}

func HeuristicScore(state *rules.BoardState, reachedSteps int) float64 {
	totalScore := float64(reachedSteps)

	// if n.Parent != nil {
	// 	// fmt.Println("yo")
	// 	totalScore = n.Parent.Score + 1
	// }

	for _, snake := range state.Snakes {

		if snake.ID == youID && snake.EliminatedCause != "" {
			return 0
		}

		if snake.ID != youID && snake.EliminatedCause != "" {
			totalScore += 20
			if snake.EliminatedBy == youID {
				totalScore += 10
			}
		}

		if snake.ID == youID && snake.Health == 100 {
			totalScore += 50
		}
	}

	return totalScore
}

func WeDied(state *rules.BoardState) bool {

	for _, snake := range state.Snakes {

		if snake.ID == youID && snake.EliminatedCause != "" {
			return true
		}

	}
	return false
}

// func (n *Node) UpperConfidenceBound() {

// }

// bySelectionScore implements sort.Interface to sort *descending* by selection score.
// Example: sort.Sort(bySelectionScore(nodes))
type bySelectionScore []*Node

func (a bySelectionScore) Len() int           { return len(a) }
func (a bySelectionScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bySelectionScore) Less(i, j int) bool { return a[i].SelectionScore > a[j].SelectionScore }

// byVisits implements sort.Interface to sort *descending* by visits.
// Example: sort.Sort(byVisits(nodes))
type byVisits []*Node

func (a byVisits) Len() int           { return len(a) }
func (a byVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byVisits) Less(i, j int) bool { return a[i].Visits > a[j].Visits }
