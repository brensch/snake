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

	}

	if parent == nil {
		n.Score = &score
		return
	}

	if !n.IsMaximising {
		if parent.Alpha <= score {
			parent.Alpha = score
			n.Score = &score
		}
		return
	}

	if parent.Beta >= score {
		parent.Beta = score
		n.Score = &score
	}
}

func (n *Node) FindBestChild() *Node {
	if n.Score == nil {
		fmt.Println("found no score on node", n.Alpha, n.Beta)
		return nil
	}

	if len(n.Children) == 0 {
		fmt.Println("reached bottom of tree")
		return nil
	}
	for _, child := range n.Children {

		if child.Score != nil && *child.Score == *n.Score {
			return child
		}
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
func (n *Node) DeepeningSearch(ctx context.Context, ruleset rules.Ruleset) rules.BoardState {

	// start with a modest check
	// startingNode :=

	depth := 10
	// fmt.Println("---------------------------------------- start depth", depth)
	n.Search(ctx, depth, depth, ruleset, nil)
	bestChild := n.FindBestChild()
	bestState := *bestChild.State
	// fmt.Println("best score was", *bestChild.Score)

	// generator.PrintMap(n.FindBestChild().State)
	// fmt.Println("---------------------------------------- end depth", depth)

	for {
		// set score to nil since i only set the child score to nil when retracing
		// so this would be missed on the first node.
		n.Score = nil
		// remove all work as a test
		n.Children = nil
		n.Alpha = math.Inf(-1)
		n.Beta = math.Inf(1)

		// get best move before each round, only
		depth++
		// fmt.Println("---------------------------------------- start depth", depth)
		deepestDepth, err := n.Search(ctx, depth, depth, ruleset, nil)
		// fmt.Printf("checking depth %d, deepest: %d\n", depth, deepestDepth)
		if err != nil {
			// fmt.Println("---------------------------------------- end depth", depth)

			break
		}
		// if ctx.Err() != nil {
		// 	fmt.Println("error not picked up in search")
		// 	break
		// }
		// this means we fully explored the tree and should give up.
		if deepestDepth > 0 {
			// fmt.Println("---------------------------------------- end depth", depth)
			break

		}
		bestChild = n.FindBestChild()
		bestState = *bestChild.State
		// fmt.Println("best score was", *bestChild.Score)

		// generator.PrintMap(n.FindBestChild().State)
		// fmt.Println("deepest depth", deepestDepth)

		// n.ExploreBestPath()

		// fmt.Println("---------------------------------------- end depth", depth)

	}
	fmt.Println("got to depth", depth-1)
	return bestState

}

func (n *Node) Search(ctx context.Context, depth, deepestDepth int, ruleset rules.Ruleset, parent *Node) (int, error) {

	// fmt.Println("player", player)

	if ctx.Err() != nil {
		// fmt.Println("context finished")
		return depth, ctx.Err()
	}

	// if we have reached a new depth, set it for returning to all subsequent recursive calls
	if depth < deepestDepth {
		deepestDepth = depth
	}

	// // we've been here before. reset order children, then reset score
	// if len(n.Children) > 0 {

	// 	for _, child := range n.Children {
	// 		child.Alpha = math.Inf(-1)
	// 		child.Beta = math.Inf(1)
	// 		child.Score = nil
	// 		var err error
	// 		deepestDepth, err = child.Search(ctx, depth-1, deepestDepth, ruleset, n)
	// 		if err != nil {
	// 			return deepestDepth, err
	// 		}

	// 		if n.Alpha >= n.Beta {
	// 			return deepestDepth, err
	// 		}
	// 	}

	// 	// once all moves evaluated, figure out score from alpha and beta
	// 	score := n.Alpha
	// 	if !n.IsMaximising {
	// 		score = n.Beta
	// 	}

	// 	n.PropagateScore(parent, score)
	// 	return deepestDepth, nil

	// }

	finishedScore := GameFinished(n.State)
	if finishedScore != 0 {
		// fmt.Println("got score", finishedScore)
		n.PropagateScore(parent, finishedScore)
		// n.Score = &finishedScore

		return deepestDepth, nil
	}

	// check if we've been here before to fastforward
	// if n.Score != nil && len(n.Children) > 0 {

	// }

	// generator.PrintMap(n.State)
	if depth == 0 {

		control := HeuristicAnalysis(n.State)
		// fmt.Println("----------")
		// fmt.Println("got score", control)
		// generator.PrintMap(n.State)
		// fmt.Println("----------")

		n.PropagateScore(parent, control)

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
			print(err)
			print(n.State)
			panic(err)
		}

		// fmt.Println("made new node")
		childNode := &Node{
			Alpha:        n.Alpha,
			Beta:         n.Beta,
			IsMaximising: !n.IsMaximising,
			State:        nextState,
		}

		// fmt.Println("child", moveNumber)

		n.Children[moveNumber] = childNode
		moveNumber++

		deepestDepth, err := childNode.Search(ctx, depth-1, deepestDepth, ruleset, n)
		if err != nil {
			return deepestDepth, err
		}

		if n.Alpha >= n.Beta {
			fmt.Println("pruning")
			return deepestDepth, err
		}

	}

	// if moveNumber != count {
	// 	fmt.Println("crazy")
	// }

	// go through each and do a depth 0 search so we can sort them.
	// this improves alpha beta pruning.

	// if depth > 1 {

	// 	for _, child := range n.Children {
	// 		_, err := child.Search(ctx, 0, 0, ruleset, n)
	// 		if err != nil {
	// 			return deepestDepth, err
	// 		}
	// 	}

	// 	sort.Sort(n.Children)

	// }

	// and now do the search
	// for number, child := range n.Children {
	// 	fmt.Println("child", number)

	// 	deepestDepth, err := child.Search(ctx, depth-1, deepestDepth, ruleset, n)
	// 	if err != nil {
	// 		return deepestDepth, err
	// 	}

	// 	if n.Alpha >= n.Beta {
	// 		return deepestDepth, err
	// 	}
	// }

	// once all moves evaluated, figure out score from alpha and beta
	score := n.Alpha
	if !n.IsMaximising {
		score = n.Beta
	}

	n.PropagateScore(parent, score)

	// and sort so future iterations benefit
	// sort.Sort(n.Children)

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

// // New returns a new minimax structure
// func New() Node {
// 	n := Node{isOpponent: false}
// 	return n
// }

// // GetBestChildNode returns the first child node with the matching score
// func (node *Node) GetBestChildNode() *Node {
// 	for _, cn := range node.children {
// 		if cn.Score == node.Score {
// 			return cn
// 		}
// 	}

// 	return nil
// }

// Evaluate runs through the tree and caculates the score from the terminal nodes
// all the the way up to the root node
// func (node *Node) Evaluate() {
// 	for _, cn := range node.children {
// 		if !cn.isTerminal() {
// 			cn.Evaluate()
// 		}

// 		if cn.parent.Score == nil {
// 			cn.parent.Score = cn.Score
// 		} else if cn.isOpponent && *cn.Score > *cn.parent.Score {
// 			cn.parent.Score = cn.Score
// 		} else if !cn.isOpponent && *cn.Score < *cn.parent.Score {
// 			cn.parent.Score = cn.Score
// 		}
// 	}
// }

// // Print the node for debugging purposes
// func (node *Node) Print(level int) {
// 	var padding = ""
// 	for j := 0; j < level; j++ {
// 		padding += " "
// 	}

// 	var s = ""
// 	if node.Score != nil {
// 		s = strconv.Itoa(*node.Score)
// 	}

// 	fmt.Println(padding, node.isOpponent, node.Board, "["+s+"]")

// 	for _, cn := range node.children {
// 		level += 2
// 		cn.Print(level)
// 		level -= 2
// 	}
// }

// // AddTerminal adds a terminal node (or leaf node).  These nodes
// // should contain a score and no children
// func (node *Node) AddTerminal(score int, board *rules.BoardState) *Node {
// 	return node.add(&score, board)
// }

// // Add a new node to structure, this node should have children and
// // an unknown score
// func (node *Node) Add(board *rules.BoardState) *Node {
// 	return node.add(nil, board)
// }

// func (node *Node) add(score *int, state *rules.BoardState) *Node {
// 	childNode := Node{parent: node, Score: score, State: state}

// 	childNode.player = (node.player + 1) % 2
// 	node.children = append(node.children, &childNode)
// 	return &childNode
// }

// func (node *Node) isTerminal() bool {
// 	return len(node.children) == 0
// }
