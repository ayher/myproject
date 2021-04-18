package main

import (
	bitcoinsrc "myproject/bitcoin/src"
	"myproject/public/fmt"
	"myproject/cryptography/keySaveInFile/src"
)

func main()  {
	pri,pub:=bitcoinsrc.ImportPrikey("dff93638b5b0a53cda4d87873dbc98104185b952c487cea09be7c1b7bd839220")
	add:=bitcoinsrc.PubkeyToAddress(pub)

	fmt.Println("private key:",pri)
	fmt.Println("public key:",pub)
	fmt.Println("address:",add)

	var err error
	err=src.Save("RSA PRIVATE KEY",pri,"./key/pricate.pem")
	if err!=nil{
		fmt.Error(err)
	}
	_,c,err:=src.Read("./key/pricate.pem")
	if err!=nil{
		fmt.Error(err)
	}
	fmt.Println("result:",c)
}


