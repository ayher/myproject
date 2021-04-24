package main

import (
	"fmt"
	"myproject/src/bitcoin/src"
)

//8353c1e67ba0dc3061dd1220a5ebbc5ad8cbc5fdb2c7f4d428696ab0a9ea2176
//0469a1c7cd7bc162b98c98533b024d5a0c6fdb55ffa508e5a9f22f9a659409c850fccc2ec15549cfd2583c61247dd71123ec9becc180deecd6d5bd97e7f2f7dfda
//14EHU1tLH2jFJCTq47qWSkFLD6172k2FdY

func main()  {
	pri,pub:=src.NewKey()
	add:=src.PubkeyToAddress(pub)
	fmt.Println(pri)
	fmt.Println(pub)
	fmt.Println(add)
}

