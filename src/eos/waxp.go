package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/token"
	"runtime/debug"
	"strconv"
	"strings"
)

/*
https://www.alohaeos.com/zh/tools/testnets
网站:
	waxsweden.org/testnet/
	Monitor / Explorer:
	wax-test.bloks.io
电报:
	waxtestnet
Aloha EOS API端点:
	https://api.waxtest.alohaeos.com
Aloha EOS P2P端点:
	peer.waxtest.alohaeos.com:9876


账户：
{
    "msg": "succeeded",
    "keys": {
        "active_key": {
            "public": "EOS7vSuEuBZDUnBJgU1rcN1MCJzbitVomWgjRcPc9gJFqF81FUG6W",
            "private": "5KPgrnorzh7yzT47LUyBUMrzr4bHWJ6pfgEnVQGxfMHzvmSCRvE"
        },
        "owner_key": {
            "public": "EOS7DLExY53ZEqix5akgmN6CgYRG3S1DinUwE6njHY2Vf6nwfpxcz",
            "private": "5JbADJuFVJJKh9GNxvUJuoNDB4Bq48PBmSvW1dFbznp1jCkvM1Q"
        }
    },
    "account": "karlmarxbest"
}

{
    "msg": "succeeded",
    "keys": {
        "active_key": {
            "public": "EOS8kRXUisC93sLdmVLKp5kkmo9EWmLAEzVwdMpq7aH8JkpsVqRMV",
            "private": "5JznyQPM8vsqipmD5BmZGtcJ1CW2g59q8tjtMeiVMw4gh8CgJjH"
        },
        "owner_key": {
            "public": "EOS6mAmfehCfpest5xUraibrZsVeygD4TCTN9HpSa4k39f5egMLtB",
            "private": "5K4uYoCXQzvCZCvfrdqqGQXUozrf87WfonJcH1Zcaqo4r5aML4w"
        }
    },
    "account": "karlmarxtest"
}

*/



//func main() {
//	//test()
//	//return
//	//api := getApi("https://eos.getblock.io")
//	api := eos.New("https://api.waxtest.alohaeos.com")
//
//	// 将私钥保存在本地
//	keyBag := &eos.KeyBag{}
//	err := keyBag.ImportPrivateKey(readPrivateKeyWxp())
//	if err != nil {
//		panic(fmt.Errorf("import private key: %s", err))
//	}
//	api.SetSigner(keyBag)
//
//	from := eos.AccountName("karlmarxbest")
//	to := eos.AccountName("karlmarxtest")
//	var WAXSymbol = eos.Symbol{Precision: 8, Symbol: "WAX"}
//	quantity, err := eos.NewFixedSymbolAssetFromString(WAXSymbol,strconv.FormatFloat(1, 'g', -1, 64)+" WAX")
//	memo := ""
//
//	if err != nil {
//		panic(fmt.Errorf("invalid quantity: %s", err))
//	}
//
//	txOpts := &eos.TxOptions{}
//	if err := txOpts.FillFromChain(api); err != nil {
//		panic(fmt.Errorf("filling tx opts: %s", err))
//	}
//	// 构建交易
//	tx := eos.NewTransaction([]*eos.Action{token.NewTransfer(from, to, quantity, memo)}, txOpts)
//
//	// 签名
//	signedTx, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
//	if err != nil {
//		panic(fmt.Errorf("sign transaction: %s", err))
//	}
//
//	content, err := json.MarshalIndent(signedTx, "", "  ")
//	if err != nil {
//		panic(fmt.Errorf("json marshalling transaction: %s", err))
//	}
//
//	fmt.Println(string(content))
//
//	// 广播
//	response, err := api.PushTransaction(packedTx)
//	if err != nil {
//		panic(fmt.Errorf("push transaction: %s", err))
//	}
//
//	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
//}
//
//func readPrivateKeyWxp() string {
//	privateKey := "5KPgrnorzh7yzT47LUyBUMrzr4bHWJ6pfgEnVQGxfMHzvmSCRvE"
//
//	return privateKey
//}

func main()  {
	data:=&Abbc{
		Url:	"https://api.waxtest.alohaeos.com",
		Coin: "WAX",
		From:   "karlmarxbest",
		To:     "karlmarxtest",
		Value:  10,
		Memo:	"123456",
		Precision: 8,
	}
	abbc, err := NewAbbc(data)
	if err!=nil {
		fmt.Println(err)
		return
	}

	summary,extInfo,_,err:=abbc.CalcSummary()
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(extInfo)

	signnature,err:=abbc.Sign(summary,"5KPgrnorzh7yzT47LUyBUMrzr4bHWJ6pfgEnVQGxfMHzvmSCRvE")
	if err!=nil {
		fmt.Println(err)
		return
	}

	abbc.WithSignature(signnature)

	tx,err:=abbc.Broadcast()
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
}

type Abbc struct {
	Coin     string
	Url      string
	From     string
	To       string
	Value    float64
	Memo 	 string
	Tx       *eos.Transaction
	Stx 	 *eos.SignedTransaction
	Precision uint8
}

func NewAbbc(e *Abbc) (*Abbc, error) {
	return e, nil
}

func (self *Abbc) NeedSign() bool {
	return false
}

func (self *Abbc) WithSignature (sign string) (extInfo string, ret error){
	signnature,err:=ecc.NewSignature(sign)
	if err!=nil{
		fmt.Println(err)
		return "",err
	}
	self.Stx.Signatures = append(self.Stx.Signatures, signnature)
	return "",nil
}

func (self *Abbc) CalcSummary() (summary string, extInfo string, nextSignAddress string, ret error) {
	defer func() {
		if errs := recover(); errs != nil {
			ret = fmt.Errorf("%+v, %s", errs, string(debug.Stack()))
		}
	}()

	api := eos.New(self.Url)
	from := eos.AccountName(self.From)
	to := eos.AccountName(self.To)

	var ABBCSymbol = eos.Symbol{Precision: self.Precision, Symbol: strings.ToTitle(self.Coin)}
	quantity, err := eos.NewFixedSymbolAssetFromString(ABBCSymbol,strconv.FormatFloat(1, 'g', -1, 64)+" "+strings.ToTitle(self.Coin))
	memo := self.Memo

	if err != nil {
		errs:=fmt.Errorf("invalid quantity: %s", err)
		fmt.Println(errs)
		return "","","",errs
	}

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(api); err != nil {
		errs:=fmt.Errorf("filling tx opts: %s", err)
		fmt.Println(errs)
		return "","","",errs
	}
	// 构建交易
	tx := eos.NewTransaction([]*eos.Action{token.NewTransfer(from, to, quantity, memo)}, txOpts)

	stx := eos.NewSignedTransaction(tx)

	txdata, cfd, err := stx.PackedTransactionAndCFD()
	if err != nil {
		fmt.Println(err)
		return "","","",err
	}

	sigDigest := eos.SigDigest(txOpts.ChainID, txdata, cfd)

	self.Tx=tx
	extInfo, err = self.toString()
	if err != nil {
		fmt.Println(err)
		return "", "", "", err
	}

	self.Stx=stx
	return hex.EncodeToString(sigDigest),extInfo,self.From,nil
}

func (self *Abbc) toString() (data string, ret error) {
	txBytes, err := json.Marshal(self.Tx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func (self *Abbc) Broadcast() (txid string, ret error) {
	defer func() {
		if errs := recover(); errs != nil {
			ret = fmt.Errorf("%+v, %s", errs, string(debug.Stack()))
		}
	}()

	api := eos.New(self.Url)

	packedTx, err := self.Stx.Pack(eos.CompressionNone)
	// 广播
	response, err := api.PushTransaction(packedTx)
	if err != nil {
		errs:=fmt.Errorf("push transaction: %s", err)
		fmt.Println(errs)
		return "",errs
	}

	return hex.EncodeToString(response.Processed.ID),nil
}

func (self *Abbc) Sign(summary, secret string) (string, error) {
	api:=&eos.API{}

	// 将私钥保存在本地
	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(secret)
	if err != nil {
		errs:=fmt.Errorf("import private key: %s", err)
		fmt.Println(errs)
		return "",errs
	}
	api.SetSigner(keyBag)

	summaryByte,err:=hex.DecodeString(summary)
	if err!=nil{
		fmt.Println(err)
		return "",err
	}
	requiredKeys, err := api.Signer.AvailableKeys()
	if err!=nil{
		fmt.Println(err)
		return "",err
	}
	signnature,err:=keyBag.SignDigest(summaryByte,requiredKeys[0])

	return signnature.String(),nil
}