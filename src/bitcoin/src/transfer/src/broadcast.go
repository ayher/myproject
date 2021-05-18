package src

import (
	"errors"
	"myproject/public/fmt"
	"myproject/public/request"
)

func (self *BtcTransfer)Broadcast() (string,error) {

	data,err:=request.Http.Curl("-d","tx_hex="+self.HexSignedTx,"https://chain.so/api/v2/send_tx/btctest")
	if err!=nil{
		return "",err
	}
	fmt.Println(data)
	status:=data["status"]
	if status=="success"{
		d:=data["data"].(map[string]interface{})
		txid:=d["txid"].(string)
		return txid,nil
	}
	return "",errors.New("err")
}