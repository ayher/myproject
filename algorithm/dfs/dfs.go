package dfs

import (
	twoTree "myproject/algorithm/two-tree"
	"myproject/public/fmt"
)

func DFS(node *twoTree.Tree)  {
	if node==nil{
		return
	}else{
		fmt.Println(node.Value)
		DFS(node.L)
		DFS(node.R)
	}
}
