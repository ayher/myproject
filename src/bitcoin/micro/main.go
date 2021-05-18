package main

import (

	"myproject/public/micro"
	"myproject/src/bitcoin/micro/server"

)

func main()  {
	micro:=micro.Micro{
		"btc",
	}
	micro.Register(&server.Btc{})
}