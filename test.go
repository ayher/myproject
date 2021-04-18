package main

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
)

func main()  {
	fmt.Println("12312312")
	wif:=""
	decoded := base58.Decode(wif)
	privKeyBytes := decoded[1 : 1+32]
	pri, pub := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	//pri,pub:=src.ImportPrikey(hex.EncodeToString(b))
	privateKey := base58.Encode(pri.Serialize())
	publicKey := base58.Encode(pub.SerializeUncompressed())


	fmt.Println(privateKey,publicKey)
}