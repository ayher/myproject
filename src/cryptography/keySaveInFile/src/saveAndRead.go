package src

import (
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

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