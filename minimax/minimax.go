package minimax

import (
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
func Search(player int, state *rules.BoardState, ruleset rules.Ruleset, depth int, alpha, beta float64) (*rules.BoardState, float64) {

	// fmt.Println("player", player)

	finishedCheck := GameFinished(state, player)
	if finishedCheck != 0 {
		return state, finishedCheck
	}

	// generator.PrintMap(state)
	if depth == 0 {
		// TODO: heuristic value
		control := HeuristicAnalysis(state, player)
		// fmt.Println("hit bottom", control, player)
		// generator.PrintMap(state)
		return state, control

	}

	// bestMoveScore := math.Inf(-1)
	var bestMove *rules.BoardState
	// maxValue := rules.DirectionUnknown
	safeMoves := generator.AllMovesForSnake(state, player)
	for move, isOk := range safeMoves {
		if !isOk {
			continue
		}

		snakeMove := rules.SnakeMoveIndex{
			Index: player,
			Move:  rules.Direction(move),
		}

		nextState, err := ruleset.ApplySingleMove(state, snakeMove)
		if err != nil {
			print(err)
			print(state)
			panic(err)
		}

		_, tmpScore := Search((player+1)%2, nextState, ruleset, depth-1, -beta, -alpha)
		tmpScore = -tmpScore

		// if score > bestMoveScore {
		// 	bestMoveScore = score
		// 	bestMove = potentialMove
		// }

		if tmpScore > alpha {
			// fmt.Println("found new alpha", tmpScore, "player", player)
			alpha = tmpScore
			bestMove = nextState
			if beta <= alpha {
				// fmt.Println("pruning", player, beta, alpha)
				return bestMove, beta
			}
		}

		if bestMove == nil {
			bestMove = nextState
		}

	}

	if bestMove == nil {
		bestMove = state
		// fmt.Println("no good moves found")
	}

	return bestMove, alpha

}

// Node represents an element in the decision tree
type Node struct {
	// Score is available when supplied by an evaluation function or when calculated
	Score      *int
	parent     *Node
	children   []*Node
	isOpponent bool

	Board *rules.BoardState
}

// New returns a new minimax structure
func New() Node {
	n := Node{isOpponent: false}
	return n
}

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
func (node *Node) Evaluate() {
	for _, cn := range node.children {
		if !cn.isTerminal() {
			cn.Evaluate()
		}

		if cn.parent.Score == nil {
			cn.parent.Score = cn.Score
		} else if cn.isOpponent && *cn.Score > *cn.parent.Score {
			cn.parent.Score = cn.Score
		} else if !cn.isOpponent && *cn.Score < *cn.parent.Score {
			cn.parent.Score = cn.Score
		}
	}
}

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

// AddTerminal adds a terminal node (or leaf node).  These nodes
// should contain a score and no children
func (node *Node) AddTerminal(score int, board *rules.BoardState) *Node {
	return node.add(&score, board)
}

// Add a new node to structure, this node should have children and
// an unknown score
func (node *Node) Add(board *rules.BoardState) *Node {
	return node.add(nil, board)
}

func (node *Node) add(score *int, board *rules.BoardState) *Node {
	childNode := Node{parent: node, Score: score, Board: board}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)
	return &childNode
}

func (node *Node) isTerminal() bool {
	return len(node.children) == 0
}
