package main

import (
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	bitcoinsrc "myproject/bitcoin/src"
	"myproject/public/fmt"
	"os"
	"runtime"
	"strings"
)

func main()  {
	pri,pub:=bitcoinsrc.ImportPrikey("dff93638b5b0a53cda4d87873dbc98104185b952c487cea09be7c1b7bd839220")
	add:=bitcoinsrc.PubkeyToAddress(pub)

	fmt.Println("private key:",pri)
	fmt.Println("public key:",pub)
	fmt.Println("address:",add)

	var err error
	err=Save("RSA PRIVATE KEY",pri,"./key/pricate.pem")
	if err!=nil{
		fmt.Error(err)
	}
	_,c,err:=Read("./key/pricate.pem")
	if err!=nil{
		fmt.Error(err)
	}
	fmt.Println("result:",c)
}

func Save(filetype,content,path string) error {
	_, fileStr, _, _ := runtime.Caller(1)
	w:=strings.LastIndex(fileStr,"/")
	os.Chdir(fileStr[:w])

	b,_:=hex.DecodeString(content)
	block := &pem.Block{
		Type:  filetype,
		Bytes: b,
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func Read(path string) (string,string,error) {
	_, fileStr, _, _ := runtime.Caller(1)
	w:=strings.LastIndex(fileStr,"/")
	os.Chdir(fileStr[:w])

	file, err := os.Open(path)
	if err != nil {
		return "","",err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)

	block,_:=pem.Decode(content)
	return block.Type,hex.EncodeToString(block.Bytes),nil
}
