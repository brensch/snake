package minimax

import "github.com/brensch/snake/generator"

func (n *Node) ExploreBestPath() {

	nextChild := n
	for {
		// fmt.Println("got score of", *nextChild.Score)
		temp := nextChild.FindBestChild()
		if temp == nil {
			break
		}
		nextChild = temp
		nextChild.Print()
		generator.PrintMap(nextChild.State)
	}
}
