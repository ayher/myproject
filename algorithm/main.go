package main

import (
	"myproject/algorithm/dfs"
	two_tree "myproject/algorithm/two-tree"
)

func main()  {
	h:=two_tree.Create([]int{1,2,4,8,0,0,9,0,0,5,10,0,0,11,0,0,3,6,12,0,0,0,7,0,0})
	dfs.DFS(h)
}