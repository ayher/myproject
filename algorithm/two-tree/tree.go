package two_tree

import "myproject/public/fmt"

type Tree struct {
	Value int
	L,R *Tree
}

var i=-1

func Create(value []int) *Tree{
	i++
	if value[i]==0{
		return nil
	}else{
		node:=&Tree{}
		node.Value=value[i]
		node.L=Create(value)
		node.R=Create(value)
		return node
	}
}

func Preorder(head *Tree){
	if head==nil{
		return
	}else{
		fmt.Println(head.Value)
		Preorder(head.L)
		Preorder(head.R)
	}
}