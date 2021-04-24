package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"math/big"
	"runtime/debug"
	"github.com/ethereum/go-ethereum/crypto"
)

type BtcConfig struct {
	dbObj                 orm.Ormer
	HotAddress            string                // btc热钱包地址
	HedgeAddress          string                // 对冲钱包地址
	HotAddressBalanceWarn float64               // 热钱包地址小于此值则报警
	HedgeThreshold        float64               // 触发对冲余额值
	HedgeReserve          float64               // 对冲留底数
	NetParam              *chaincfg.Params      // 网络
	TxVersion             int32                 // wire.TxVersion
	RpcClientConfig       *rpcclient.ConnConfig // btc rpc节点链接信息
}

type Btc struct {
	coin      string
	btcConfig *BtcConfig
	rpcClient *rpcclient.Client
	Data      *BtcData
}

type CoinParam struct {
	Coin  string
	Chain string
	From  string
	To    string
	Value float64
	Memo  string
	//一些特殊字段
	ExtEthGasPrice uint64 //eth专用
	ExtTrxTx       string // 通用存储data
}

type BtcData struct { // 这个结构体要和仕学那边程序保持一致
	RedeemTx *wire.MsgTx
	UTXOs    []btcjson.ListUnspentResult
}

func NewBtc(btcConfig *BtcConfig, coinParam CoinParam) (*Btc, error) {
	client, err := rpcclient.New(btcConfig.RpcClientConfig, nil)
	if err != nil {
		return nil, err
	}
	btc := &Btc{
		coin:      coinParam.Coin,
		btcConfig: btcConfig,
		rpcClient: client,
	}
	if coinParam.ExtTrxTx != "" {
		if _, err := btc.fromString(coinParam.ExtTrxTx); err != nil {
			return nil, err
		}
	}
	return btc, nil
}

func (d *Btc) NeedSign() bool {

	// 找到第一个没签名的in
	idx := 0
	inLen := len(d.Data.RedeemTx.TxIn)
	for idx = 0; idx < inLen; idx++ {
		if len(d.Data.RedeemTx.TxIn[idx].SignatureScript) == 0 {
			break
		}
	}

	// 没有全签
	if idx >= inLen {
		return false
	}

	return true
}

func (d *Btc) WithSignature(sign string) (extInfo string, err error) {

	// 附上签名

	// 是否全签名了
	if !d.NeedSign() {
		return "", fmt.Errorf("no need to sign")
	}

	// 找到idx的in，没有签名，从这里开始附
	idx := 0
	inLen := len(d.Data.RedeemTx.TxIn)
	for idx = 0; idx < inLen; idx++ {
		if len(d.Data.RedeemTx.TxIn[idx].SignatureScript) == 0 {
			break
		}
	}

	// 添加签名，sign是RSV三个big.Int按byte顺序拼接的，长度32+32+1
	buffer, err := hex.DecodeString(sign)
	if err != nil {
		return "", err
	}

	R := new(big.Int).SetBytes(buffer[:32])
	S := new(big.Int).SetBytes(buffer[32:64])

	signature := &btcec.Signature{
		R: R,
		S: S,
	}

	// 拿到pkData
	var pkData []byte
	signHashType := txscript.SigHashAll
	compress := false
	publicKeyHex, err := AddressToPublicKey(d.btcConfig.dbObj, d.Data.UTXOs[idx].Address)
	if err != nil {
		return "", err
	}

	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", err
	}

	fromPublicKey, err := btcec.ParsePubKey(publicKey, btcec.S256())
	if err != nil {
		return "", err
	}

	if compress {
		pkData = fromPublicKey.SerializeCompressed()
	} else {
		pkData = fromPublicKey.SerializeUncompressed()
	}

	rawTxInSignature := append(signature.Serialize(), byte(signHashType))
	signatureScript, err := txscript.NewScriptBuilder().AddData(rawTxInSignature).AddData(pkData).Script()

	d.Data.RedeemTx.TxIn[idx].SignatureScript = signatureScript

	return d.toString()
}

func (d *Btc) CalcSummary() (summary string, extInfo string, nextSignAddress string, ret error) {
	defer func() {
		if errs := recover(); errs != nil {
			ret = fmt.Errorf("%+v, %s", errs, string(debug.Stack()))
		}
	}()

	signHashType := txscript.SigHashAll

	// 计算摘要
	// 找到idx的in，没有签名，从这里开始计算摘要
	idx := 0
	for idx = 0; idx < len(d.Data.RedeemTx.TxIn); idx++ {
		if len(d.Data.RedeemTx.TxIn[idx].SignatureScript) == 0 {
			break
		}
	}

	sourcePKScript, err := hex.DecodeString(d.Data.UTXOs[idx].ScriptPubKey)
	if err != nil {
		return "", "", "", err
	}

	hash, err := txscript.CalcSignatureHash(sourcePKScript, signHashType, d.Data.RedeemTx, idx)
	if err != nil {
		return "", "", "", err
	}
	summary = hex.EncodeToString(hash)
	logs.Debug("[Summary][%s]: %s", d.coin, summary)
	extInfo, err = d.toString()
	if err != nil {
		return "", "", "", err
	}

	// 以此地址来签名
	nextSignAddress = d.Data.UTXOs[idx].Address

	return
}

func (d *Btc) Broadcast() (txid string, ret error) {
	defer func() {
		if errs := recover(); errs != nil {
			ret = fmt.Errorf("%+v, %s", errs, string(debug.Stack()))
		}
	}()

	// 总raw输出
	var signedTx bytes.Buffer
	err := d.Data.RedeemTx.Serialize(&signedTx)
	if err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	logs.Debug("[Broadcast][%s]: %s", d.coin, hexSignedTx)

	// 广播raw
	hash, err := d.rpcClient.SendRawTransaction(d.Data.RedeemTx, true)
	if err != nil {
		logs.Error("[Broadcast][%s] err: %+v", d.coin, err)
		return "", err
	}

	return hash.String(), nil

}

func (d *Btc) fromString(data string) (*BtcData, error) {
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}

	btcData := new(BtcData)
	if err := json.Unmarshal(dataBytes, btcData); err != nil {
		return nil, err
	}

	d.Data = btcData

	return d.Data, nil
}

func (d *Btc) toString() (string, error) {

	dataBytes, err := json.Marshal(d.Data)
	if err != nil {
		return "", err
	}

	dataString := hex.EncodeToString(dataBytes)

	return dataString, nil
}

func (self *Btc) Sign(summary, secret string) (string, error) {
	digestHash, err := hex.DecodeString(summary)
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return "", err
	}
	btcPrivKey := btcec.PrivateKey{
		PublicKey: privateKey.PublicKey,
		D:         privateKey.D,
	}
	signature, err := btcPrivKey.Sign(digestHash)
	if err != nil {
		return "", err
	}

	signBytes := make([]byte, 65)
	copy(signBytes[:32], signature.R.Bytes())
	copy(signBytes[32:64], signature.S.Bytes())

	return hex.EncodeToString(signBytes), nil
}


func AddressToPublicKey (dbObj orm.Ormer, address string) (string, error){
	return "",nil
}