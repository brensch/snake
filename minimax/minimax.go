package minimax

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/brensch/snake/generator"
	"github.com/brensch/snake/rules"
)

// const (
// 	MAXPLAYERINDEX = 0
// 	MINPLAYERINDEX = 1
// )

// Search requires that the maximising player is at index 0, min at 1
// TODO: change the boardstate struct to have explicit max and min player fields
// no longer using ismaxplayer, just swap alpha and beta at every layer
// func Search(player int, state *rules.BoardState, ruleset rules.Ruleset, depth int, alpha, beta float64) (*rules.BoardState, float64) {

// 	// fmt.Println("player", player)

// 	finishedCheck := GameFinished(state, player)
// 	if finishedCheck != 0 {
// 		return state, finishedCheck
// 	}

// 	// generator.PrintMap(state)
// 	if depth == 0 {
// 		// TODO: heuristic value
// 		control := HeuristicAnalysis(state, player)
// 		// fmt.Println("hit bottom", control, player)
// 		// generator.PrintMap(state)
// 		return state, control

// 	}

// 	// bestMoveScore := math.Inf(-1)
// 	var bestMove *rules.BoardState
// 	// maxValue := rules.DirectionUnknown
// 	safeMoves := generator.AllMovesForSnake(state, player)
// 	for move, isOk := range safeMoves {
// 		if !isOk {
// 			continue
// 		}

// 		snakeMove := rules.SnakeMoveIndex{
// 			Index: player,
// 			Move:  rules.Direction(move),
// 		}

// 		nextState, err := ruleset.ApplySingleMove(state, snakeMove)
// 		if err != nil {
// 			print(err)
// 			print(state)
// 			panic(err)
// 		}

// 		_, tmpScore := Search((player+1)%2, nextState, ruleset, depth-1, -beta, -alpha)
// 		tmpScore = -tmpScore

// 		// if score > bestMoveScore {
// 		// 	bestMoveScore = score
// 		// 	bestMove = potentialMove
// 		// }

// 		if tmpScore > alpha {
// 			// fmt.Println("found new alpha", tmpScore, "player", player)
// 			alpha = tmpScore
// 			bestMove = nextState
// 			if beta <= alpha {
// 				// fmt.Println("pruning", player, beta, alpha)
// 				return bestMove, beta
// 			}
// 		}

// 		if bestMove == nil {
// 			bestMove = nextState
// 		}

// 	}

// 	if bestMove == nil {
// 		bestMove = state
// 		// fmt.Println("no good moves found")
// 	}

// 	return bestMove, alpha

// }

func (n *Node) PropagateScore(parent *Node, score float64) {

	if n.Score == nil {
		n.Score = &score
		// if score == 0.39669421487603307 {
		// 	fmt.Println(score)
		// }
		// if score == 0.39669421487603307 {
		// 	fmt.Println(score)
		// }

	}

	if parent == nil {
		n.Score = &score
		return
	}

	if !n.IsMaximising {
		if parent.Alpha <= score {
			parent.Alpha = score
			// parent.Score = &score
			// if score == 0.39669421487603307 {
			// 	fmt.Println("alpha", score)
			// }
		}
		return
	}

	if parent.Beta >= score {
		parent.Beta = score
		// parent.Score = &score
		// if score == 0.39669421487603307 {
		// 	fmt.Println("beta", score)
		// }
	}
}

func (n *Node) FindBestChild() *Node {
	if n.Score == nil {
		fmt.Println("found no score on node", n.Alpha, n.Beta)
		return nil
	}

	if len(n.Children) == 0 {
		fmt.Println("reached bottom of tree. ")

		return nil
	}
	var matchingChildren []*Node
	for _, child := range n.Children {

		if child.Score != nil && *child.Score == *n.Score {
			matchingChildren = append(matchingChildren, child)
		}
	}

	return matchingChildren[0]

	if len(matchingChildren) == 1 {
		return matchingChildren[0]

	}

	if len(matchingChildren) > 1 {
		fmt.Println("got more than one matching child")
		for _, matchingChild := range matchingChildren {
			matchingChild.Print()
			generator.PrintMap(matchingChild.State)
		}
		return matchingChildren[0]
	}

	fmt.Println("didn't find matching node", n.Alpha, n.Beta, *n.Score, len(n.Children))
	for _, child := range n.Children {
		fmt.Println("child:", child.Alpha, child.Beta)
		// generator.PrintMap(child.State)
	}

	fmt.Println("what does this mean")

	return nil
}

// Too damn slow.
func (n *Node) CopyNode() *Node {

	if n == nil {
		return nil
	}

	contents, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}

	var copiedNode *Node
	err = json.Unmarshal(contents, &copiedNode)
	if err != nil {
		panic(err)
	}

	return copiedNode
}

// Cancelling context is handled outside this functions
func (n *Node) DeepeningSearch(ctx context.Context, ruleset rules.Ruleset) *rules.BoardState {

	// i think there's a better way to do transposition tables.
	// currently heuristic scores are taking about 6000ns so if we can replace that with a 40ns lookup
	// our iterative deepening may actually do as i expect.
	previousHeuristicScores := make(map[uint64]float64)

	// 8 seems like a good tradeoff - if we don't finish at least one it's very bad
	depth := 8
	// fmt.Println("---------------------------------------- start depth", depth)

	// fmt.Println("best score was", *bestChild.Score)

	// generator.PrintMap(n.FindBestChild().State)
	// fmt.Println("---------------------------------------- end depth", depth)

	var bestState *rules.BoardState
	var bestScore float64

	for {
		deepestDepth, err := n.Search(ctx, depth, depth, ruleset, nil, previousHeuristicScores)
		if err != nil {
			break
		}
		bestChild := n.FindBestChild()
		bestState = bestChild.State
		bestScore = *bestChild.Score

		if deepestDepth > 0 {
			break
		}

		// set score to nil since i only set the child score to nil when retracing
		// so this would be missed on the first node.
		n.Score = nil
		n.Children = nil
		n.Alpha = math.Inf(-1)
		n.Beta = math.Inf(1)

		// get best move before each round, only
		depth++
		// fmt.Println("---------------------------------------- start depth", depth)
		// deepestDepth, err := n.Search(ctx, depth, depth, ruleset, nil, previousHeuristicScores)
		// // fmt.Printf("checking depth %d, deepest: %d\n", depth, deepestDepth)
		// if err != nil {
		// 	break
		// }
		// // if ctx.Err() != nil {
		// // 	fmt.Println("error not picked up in search")
		// // 	break
		// // }
		// // this means we fully explored the tree and should give up.
		// if deepestDepth > 0 {
		// 	// fmt.Println("incomplete search", deepestDepth)
		// 	// fmt.Println("---------------------------------------- end depth", depth)
		// 	break

		// }
		// bestChild = n.FindBestChild()
		// bestState = *bestChild.State
		// fmt.Println("best score was", *bestChild.Score)

		// generator.PrintMap(n.FindBestChild().State)
		// fmt.Println("deepest depth", deepestDepth)

		// n.ExploreBestPath()

		// fmt.Println("---------------------------------------- end depth", depth)

	}

	// if we have a best score of -1, then our best child is garbage. do an reassessment with no
	// lookahead and just submit that.
	if bestState != nil && bestScore != -1 {
		return bestState
	}

	fmt.Println("emergency")

	n.Score = nil
	n.Children = nil
	n.Alpha = math.Inf(-1)
	n.Beta = math.Inf(1)
	// if we didn't finish a single loop in the given time, do an emergency check of depth 0 and return
	// the best immediate child.
	n.Search(context.Background(), 0, 0, ruleset, nil, previousHeuristicScores)
	return n.FindBestChild().State

	// fmt.Println("got to depth", depth-1)
	// return bestState

}

func (n *Node) Search(ctx context.Context, depth, deepestDepth int, ruleset rules.Ruleset, parent *Node, hashTable map[uint64]float64) (int, error) {

	// fmt.Println("player", player)

	if ctx.Err() != nil {
		// fmt.Println("context finished")
		return depth, ctx.Err()
	}

	// if we have reached a new depth, set it for returning to all subsequent recursive calls
	if depth < deepestDepth {
		deepestDepth = depth
	}

	finishedScore := GameFinished(n.State, n.IsMaximising)
	if finishedScore != 0 {
		// fmt.Println("got a finished game with max", n.IsMaximising)
		// gamebytes, _ := json.Marshal(n.State)
		// fmt.Println(string(gamebytes))
		// generator.PrintMap(n.State)
		// fmt.Println("got score", finishedScore)
		n.PropagateScore(parent, finishedScore)
		// n.Score = &finishedScore

		return deepestDepth, nil
	}

	if depth == 0 {

		// check if we've seen the answer to this move before
		// TODO: do proper transposition with alpha and beta and shizzzzzzzzzz
		// hash := Hash(n.State)
		// control, ok := hashTable[hash]
		// if !ok {
		control := HeuristicAnalysis(n.State)
		// 	hashTable[hash] = control
		// }

		// fmt.Println("got heuristic at depth 0", control)
		// generator.PrintMap(n.State)
		// ShortestPathsBreadthPrint(n.State)

		n.PropagateScore(parent, control)

		// if control == 0.628099173553719 {
		// 	fmt.Println("sup")
		// 	n.Print()
		// 	generator.PrintMap(n.State)
		// 	ShortestPathsBreadthPrint(n.State)
		// 	panic("yoo")
		// }

		// n.Score = &control
		// fmt.Println("hit bottom", control, n.player)

		return deepestDepth, nil
	}

	// bestMoveScore := math.Inf(-1)
	// var bestMove *rules.BoardState
	// maxValue := rules.DirectionUnknown
	playerIndex := 0
	if !n.IsMaximising {
		playerIndex = 1
	}

	safeMoves, count := generator.AllMovesForSnake(n.State, playerIndex)
	n.Children = make([]*Node, count)
	moveNumber := 0
	for move, isOk := range safeMoves {
		if !isOk {
			continue
		}

		snakeMove := rules.SnakeMoveIndex{
			Index: playerIndex,
			Move:  rules.Direction(move),
		}

		nextState, err := ruleset.ApplySingleMove(n.State, snakeMove)
		if err != nil {
			panic(err)
		}

		childNode := &Node{
			Alpha:        n.Alpha,
			Beta:         n.Beta,
			IsMaximising: !n.IsMaximising,
			State:        nextState,
		}

		n.Children[moveNumber] = childNode
		moveNumber++

		deepestDepth, err = childNode.Search(ctx, depth-1, deepestDepth, ruleset, n, hashTable)
		if err != nil {
			return deepestDepth, err
		}

		if n.Alpha >= n.Beta {
			// fmt.Println("pruning")
			return deepestDepth, err
		}

	}

	// once all moves evaluated, figure out score from alpha and beta
	score := n.Alpha
	if !n.IsMaximising {
		score = n.Beta
	}

	n.PropagateScore(parent, score)

	return deepestDepth, nil
}

func (n *Node) Print() {
	score := -69.0
	if n.Score != nil {
		score = *n.Score
	}
	fmt.Printf("node stats: alpha %f, beta %f, score: %f ismax: %t\n", n.Alpha, n.Beta, score, n.IsMaximising)
}

// Node represents an element in the decision tree
type Node struct {
	// Score is available when supplied by an evaluation function or when calculated
	Score *float64
	// Parent   *Node
	Children NodeList

	Alpha float64
	Beta  float64

	IsMaximising bool

	State *rules.BoardState
}

type NodeList []*Node

func (nl NodeList) Len() int { return len(nl) }
func (nl NodeList) Less(i, j int) bool {

	if nl[i].Score == nil {
		return true
	}

	if nl[j].Score == nil {
		return false
	}

	return *nl[i].Score > *nl[j].Score
}
func (nl NodeList) Swap(i, j int) { nl[i], nl[j] = nl[j], nl[i] }

func Search2(state *rules.BoardState, depth uint, alpha, beta float64, isMaxer bool, ruleset rules.Ruleset) (*rules.BoardState, float64) {

	finishedScore := GameFinished(state, isMaxer)
	if finishedScore != 0 {
		// fmt.Println("got score", finishedScore)
		// n.Score = &finishedScore

		if isMaxer {
			return state, finishedScore
		}
		return state, -finishedScore
	}
	// ending, heuristic := state.Evaluate()
	// switch ending {
	// case LOSS:
	// 	return state, math.Inf(-1)
	// case TIE:
	// 	return state, 0
	// case WIN:
	// 	return state, math.Inf(1)
	// }

	if depth == 0 {
		heuristic := HeuristicAnalysis(state)
		fmt.Println("got heuristic", heuristic)
		generator.PrintMap(state)
		ShortestPathsBreadthPrint(state)

		if heuristic == 0.12809917355371903 {
			generator.PrintMap(state)
		}
		return state, heuristic
	}
	index := 0
	if !isMaxer {
		index = 1
	}

	// children := state.Children()
	safeMoves, _ := generator.AllMovesForSnake(state, index)
	// children = make([]*Node, count)

	// moveScores := make(moveScores, len(children))
	// for i := range children {
	// 	moveScores[i] = moveScore{i, 0.0}
	// }
	var tmpScore float64
	// if depth > 1 {
	// 	// Pre-sort the possible moves by their score to speed up the pruning
	// 	for i, child := range children {
	// 		// Depth-0 search to force heuristic scoring
	// 		_, tmpScore = Search(child, 0, -beta, -alpha)
	// 		moveScores[i].moveScore = -tmpScore
	// 	}
	// 	sort.Sort(moveScores)
	// }
	var bestChild *rules.BoardState
	for move, isOk := range safeMoves {
		if !isOk {
			continue
		}

		snakeMove := rules.SnakeMoveIndex{
			Index: index,
			Move:  rules.Direction(move),
		}

		child, err := ruleset.ApplySingleMove(state, snakeMove)
		if err != nil {
			panic(err)
		}
		// child := children[moveScore.moveIndex]

		_, tmpScore = Search2(child, depth-1, -beta, -alpha, !isMaxer, ruleset)
		tmpScore = -tmpScore
		if tmpScore > alpha {
			alpha = tmpScore
			bestChild = child
			if beta <= alpha {
				return bestChild, beta
			}
		}
		if bestChild == nil {
			// Take the first child, in case all the children are terrible.
			bestChild = child
		}
	}
	if bestChild == nil {
		// No possible moves, so return the current state.
		bestChild = state
	}
	return bestChild, alpha
}
