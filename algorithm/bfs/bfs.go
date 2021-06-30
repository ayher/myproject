package bfs

import (
	"myproject/algorithm/queue"
	twoTree "myproject/algorithm/two-tree"
	"myproject/public/fmt"
)

func BFS(node *twoTree.Tree){
	if node==nil{
		return
	}else{
		fmt.Println(node.Value)
		queue.Push(node.L)
		queue.Push(node.R)

		n,err:=queue.Pop()
		if err!=nil{
			return
		}
		BFS(n.(*twoTree.Tree))
	}
}
