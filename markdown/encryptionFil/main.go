package main

import (
	"myproject/markdown/encryptionFil/src"
	"myproject/public/fmt"
)


func main()  {
	pri:="dff93638b5b0a53cda4d87873dbc98104185b952c487cea09be7c1b7bd839220"
	fmt.Println("private key:",pri)
	var err error
	err=src.Save("RSA PRIVATE KEY",pri,"./key/private.pem")
	if err!=nil{
		fmt.Error(err)
	}
	_,c,err:=src.Read("./key/private.pem")
	if err!=nil{
		fmt.Error(err)
	}
	fmt.Println("result:",c)
}



