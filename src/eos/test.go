
package main

import (
	"encoding/json"
	"fmt"
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
	"encoding/hex"
	"net/http"
	"net"
	"time"
)

func getApi(baseURL string) *eos.API {
	var head=map[string][]string{
		"x-api-key":{"a62672a2-edcd-4d8e-aa35-d989e1df02cb"},
		"Content-Type":{"application/json"},
	}
	api := &eos.API{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true, // default behavior, because of `nodeos`'s lack of support for Keep alives.
			},
		},
		BaseURL:  baseURL,
		Compress: eos.CompressionZlib,
		Header:   head,
	}

	return api
}

//5JYi7wRZdgWF8oJABWtfHhS2LGNLJMJMnT3ChmXV7mpX92nNTyX EOS6jKuwomHbSpMqw7aJkQmn1XNvjSVP949SvBbv5pgg3nfCB2tsF
//1
//EOS6mFiHGb3aHS8kJXyC4JicP7FbESUHSFAJ7P2W5fp8wR7ZE3y3W

//func test()  {
//
//	keyBag := &eos.KeyBag{}
//	err := keyBag.ImportPrivateKey("5JYi7wRZdgWF8oJABWtfHhS2LGNLJMJMnT3ChmXV7mpX92nNTyX")
//	if err!=nil{
//		fmt.Println(err)
//	}
//	fmt.Println(len(keyBag.Keys))
//	fmt.Println(keyBag.Keys[0].String(),keyBag.Keys[0].PublicKey().String())
//
//	keyBag.Keys[0]
//
//}

func mainn() {
	//test()
	//return
	//api := getApi("https://eos.getblock.io")
	api := eos.New("https://jungle3.cryptolions.io")

	// 将私钥保存在本地
	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(readPrivateKey())
	if err != nil {
		panic(fmt.Errorf("import private key: %s", err))
	}
	api.SetSigner(keyBag)

	from := eos.AccountName("karlmarxbest")
	to := eos.AccountName("karlmarxtest")
	quantity, err := eos.NewEOSAssetFromString("1 EOS")
	memo := ""

	if err != nil {
		panic(fmt.Errorf("invalid quantity: %s", err))
	}

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(api); err != nil {
		panic(fmt.Errorf("filling tx opts: %s", err))
	}
	// 构建交易
	tx := eos.NewTransaction([]*eos.Action{token.NewTransfer(from, to, quantity, memo)}, txOpts)

	// 签名
	signedTx, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		panic(fmt.Errorf("sign transaction: %s", err))
	}

	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		panic(fmt.Errorf("json marshalling transaction: %s", err))
	}

	fmt.Println(string(content))

	// 广播
	response, err := api.PushTransaction(packedTx)
	if err != nil {
		panic(fmt.Errorf("push transaction: %s", err))
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func readPrivateKey() string {
	privateKey := "5JETSXKfw8bX9FbXhczqwuRKffveSQsiBtxkcNNTMLTPdUZbGwd"

	return privateKey
}

func getAPIURL() string {
	//return "https://api.eosflare.io"
	return "https://jungle3.cryptolions.io"
}

