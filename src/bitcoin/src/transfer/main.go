package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"myproject/public/fmt"
	"myproject/src/bitcoin/src/transfer/src"
)



func main(){

	from1Byte,err:=hex.DecodeString("8353c1e67ba0dc3061dd1220a5ebbc5ad8cbc5fdb2c7f4d428696ab0a9ea2176")
	if err!=nil {
		fmt.Error(err)
	}

	from2Byte,err:=hex.DecodeString("cc9ecfd789b203e441fbfba8bd875461bd5a2eca91fc89b8ee806df4ec9ff649")
	if err!=nil {
		fmt.Error(err)
	}

	from3Byte,err:=hex.DecodeString("c0432ce2f4a189ae2b50c13cd6a2f033d83361b3747f70ac2d0907edbf46aa51")
	if err!=nil {
		fmt.Error(err)
	}

	from:=[]src.Tf{
		{
			Pri: from1Byte,
		},
		{
			Pri: from2Byte,
		},
		{
			Pri: from3Byte,
		},
	}
	to:=[]src.Tf{
		{
			Address: "mw7Ms2JjzU2nyn85ZKsv7ue8VocPmE5aQ8",
			Value: 0,
		},
	}

	btcTransfer:=&src.BtcTransfer{
		FromList: from,
		ToList:   to,
		Net: &chaincfg.TestNet3Params,
	}

	_,e:=btcTransfer.Summary()
	if e!=nil{
		fmt.Println(e)
		return
	}
	//fmt.Println(btcTransfer.RedeemTx)
	btcTransfer.Sign()
	fmt.Println(btcTransfer.HexSignedTx)
	fmt.Println(btcTransfer.Broadcast())
}