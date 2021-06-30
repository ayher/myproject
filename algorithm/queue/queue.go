package queue

import "myproject/public/fmt"

var QUEUE []interface{}
var index=0

func Push(v interface{}){
	QUEUE=append(QUEUE,v )
}

func Pop() (interface{},error) {
	if index> len(QUEUE){
		return 0,fmt.Errorf("err")
	}else{
		index++
		return QUEUE[index-1],nil
	}
}
