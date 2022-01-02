package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	tronClient "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"math/big"
	"strconv"
)

func main()  {
	t:=&TransferConfig{
		Url:"47.252.3.238:50051",
		FromPri:"906f0f8aba745f8b88323dbf5f51c1d9275a5227411854e7bda1f90927f945b8",
		To: "THpVoBk8rJksPRW1npK4zrEuvquSS36WHG",
		Value: "111",
		Contract: "TNuoKL1ni8aoshfFL1ASca1Gou9RXwAzfn",
		FeeLimit: "5000000",
	}
	t.transfer()
}

type TransferConfig struct {
	Url string
	FromPri string
	To       string        `yaml:"to"`
	Value    string        `yaml:"value"`
	Contract string        `yaml:"contract"`
	FeeLimit string        `yaml:"fee_limit"`
}

func (c *TransferConfig)transfer() error{
	b,_:=hex.DecodeString(c.FromPri)
	pr,pubKey1:=btcec.PrivKeyFromBytes(btcec.S256(),b)
	pubKey:=pubKey1.ToECDSA()
	fromAddress := tronAddress.PubkeyToAddress(*pubKey)
	fmt.Println("Public Key", "pubKey", pubKey)
	fmt.Println("Address", "Address", fromAddress.String())

	// 目的地址
	toAddress, err := tronAddress.Base58ToAddress(c.To)
	if err != nil {
		fmt.Println("Cannot convert Base58 string to address", "to", c.To)
		return nil
	}
	fmt.Println("Read config", "To", toAddress.String())

	// 转的额度
	value, ok := new(big.Int).SetString(c.Value, 10)
	if !ok {
		fmt.Println("Cannot convert string to big int", "value", c.Value)
		return nil
	}
	fmt.Println("Read config", "value", value)

	var contractAddress tronAddress.Address
	var feeLimit int64
	if c.Contract != "" {
		// trc20合约地址
		contractAddress, err = tronAddress.Base58ToAddress(c.Contract)
		if err != nil {
			fmt.Println("Cannot convert Base58 string to address", "to", c.Contract)
			return nil
		}
		fmt.Println("Read config", "CT", toAddress.String())

		feeLimit, err = strconv.ParseInt(c.FeeLimit, 10, 64)
		if err != nil {
			fmt.Println("Cannot convert string to unit64", "fee_limit", c.FeeLimit)
			return nil
		}
		fmt.Println("Read config", "feeLimit", feeLimit)
	}

	// 新client
	client := tronClient.NewGrpcClient(c.Url)
	if err := client.Start(); err != nil {
		fmt.Println("client start error", "error", err)
		return nil
	}

	var tx *api.TransactionExtention
	if c.Contract == "" {
		// TRX原生币转账
		fmt.Println("TRX Transfer---",fromAddress.String(),toAddress.String(),value.Int64())
		tx, err = client.Transfer(fromAddress.String(), toAddress.String(), value.Int64())
		if err != nil {
			fmt.Println("Transfer error", "error", err)
			return nil
		}
	} else {
		// TRC20转账
		fmt.Println("TRC20 Transfer")
		tx, err = client.TRC20Send(fromAddress.String(), toAddress.String(), contractAddress.String(), value, feeLimit)
		if err != nil {
			fmt.Println("TRC20Send error", "error", err)
			return nil
		}
	}

	// 门限签
	digest := tx.Txid
	digestString := hex.EncodeToString(digest)
	fmt.Println("Digest", "Tx digest", digestString)

	sig, err := crypto.Sign(digest, pr.ToECDSA())
	//sig,err:=pr.Sign(digest)
	if err!=nil {
		panic(err)
	}
	// 附上签名
	tx.Transaction.Signature = append(tx.Transaction.Signature, sig)
	fmt.Println("sign success", "signTx", tx)

	//return nil
	// 发起转账
	result, err := client.Broadcast(tx.Transaction)
	if err != nil {
		fmt.Println("Broadcast err", "err", err)
		return nil
	}
	if result.Code != 0 {
		fmt.Println("Broadcast err", "code", result.Code)
		fmt.Println("Broadcast err", "message", string(result.Message[:]))
		return nil
	}

	fmt.Println("Broadcast success", "txid", hex.EncodeToString(tx.Txid[:]))
	fmt.Println("Broadcast success", "explorer", "https://nile.tronscan.org/#/transaction/"+hex.EncodeToString(tx.Txid[:]))
	return nil
}

const (
	trc20TransferMethodSignature = "0xa9059cbb"
)

//func trc20Send(g *tronClient.GrpcClient,from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
//	addrB, err := address.Base58ToAddress(to)
//	if err != nil {
//		return nil, err
//	}
//	ab := common.LeftPadBytes(amount.Bytes(), 32)
//	req := trc20TransferMethodSignature + "0000000000000000000000000000000000000000000000000000000000000000"[len(addrB.Hex())-4:] + addrB.Hex()[4:]
//	req += common.Bytes2Hex(ab)[3:]
//	return g.TRC20Call(g,from, contract, req, false, feeLimit)
//}
//
//func TRC20Call(g *tronClient.GrpcClient,from, contractAddress, data string, constant bool, feeLimit int64) (*api.TransactionExtention, error) {
//	var err error
//	fromDesc := address.HexToAddress("410000000000000000000000000000000000000000")
//	if len(from) > 0 {
//		fromDesc, err = address.Base58ToAddress(from)
//		if err != nil {
//			return nil, err
//		}
//	}
//	contractDesc, err := address.Base58ToAddress(contractAddress)
//	if err != nil {
//		return nil, err
//	}
//	dataBytes, err := common.FromHex(data)
//	if err != nil {
//		return nil, err
//	}
//	ct := &core.TriggerSmartContract{
//		OwnerAddress:    fromDesc.Bytes(),
//		ContractAddress: contractDesc.Bytes(),
//		Data:            dataBytes,
//	}
//	result := &api.TransactionExtention{}
//	if constant {
//		result, err = g.triggerConstantContract(ct)
//
//	} else {
//		result, err = g.triggerContract(ct, feeLimit)
//	}
//	if err != nil {
//		return nil, err
//	}
//	if result.Result.Code > 0 {
//		return result, fmt.Errorf(string(result.Result.Message))
//	}
//	return result, nil
//
//}
