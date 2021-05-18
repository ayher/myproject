package src

import (
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"math"
	"myproject/public/fmt"
	"net/http"
	"time"
)

type UnSpent struct {
	Address            string  `json:"address"`
	TotalReceived      int64   `json:"total_received"`
	TotalSent          int64   `json:"total_sent"`
	Balance            int64   `json:"balance"`
	UnconfirmedBalance int64   `json:"unconfirmed_balance"`
	FinalBalance       int64   `json:"final_balance"`
	UTXOs              []*UTXO `json:"txrefs"`
	UnconfirmedUTXOs   []*UTXO `json:"unconfirmed_txrefs"`
}

// 根据api返回地址多个UTXO
func GetUTXOs(address string) ([]*UTXO, error) {

	// https://api.blockcypher.com/v1/btc/test3/addrs/n4TBTn6xNDjNvqoWoaaJWi33dY4sMQ2DqV??unspentOnly=true&includeScript=true&limit=50

	newURL := fmt.Sprintf("https://api.blockcypher.com/v1/btc/test3/addrs/%s??unspentOnly=true&includeScript=true&limit=50", address)

	response, err := http.Get(newURL)
	if err != nil {
		return nil,errors.New("error in GetUTXOs, http.Get:"+err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var resp *UnSpent
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil,errors.New("error in GetUTXOs, json.Unmarshal:"+err.Error())

	}

	// 过滤有效txrefs
	var utxos []*UTXO
	balance := int64(0)
	for _, u := range resp.UTXOs {
		if u.Index >= 0 && u.Spend == false && u.Value>0 {
			utxos = append(utxos, u)
			balance += u.Value
		}
	}

	return utxos, nil
}

func (self *BtcTransfer)Init() error {
	for i,v:=range self.FromList {
		_,fromPublicKey:=btcec.PrivKeyFromBytes(btcec.S256(),v.Pri)
		self.FromList[i].pub=fromPublicKey.SerializeUncompressed()

		fromAddressPubKey, err := btcutil.NewAddressPubKey(fromPublicKey.SerializeUncompressed(), self.Net)
		if err!=nil{
			return err
		}
		self.FromList[i].Address=fromAddressPubKey.EncodeAddress()

		self.FromList[i].Utxos, err = GetUTXOs(fromAddressPubKey.EncodeAddress())
	}

	return nil
}

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

func EstimateFee() (int64, error) {
	info, err := getBlockChainInfo()
	if err != nil {
		return 0, err
	}
	return info.HighFeePerKb, nil
}

func getBlockChainInfo() (*ChainInfo, error) {
	// https://api.blockcypher.com/v1/btc/test3

	newURL := `http://api.blockcypher.com/v1/btc/test3`

	response, err := http.Get(newURL)
	if err != nil {
		return nil, errors.New("error in EstimateFee, http.Get:"+err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var resp *ChainInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.New("error in EstimateFee, json.Unmarshal:"+err.Error())
	}


	return resp, nil
}

type ChainInfo struct {
	Name           string    `json:"name"`
	Height         uint32    `json:"height"`
	Time           time.Time `json:"time"`
	HighFeePerKb   int64     `json:"high_fee_per_kb"`
	MediumFeePerKb int64     `json:"medium_fee_per_kb"`
	LowFeePerKb    int64     `json:"low_fee_per_kb"`
}

func (self *BtcTransfer)Summary() (string,error){
	self.Init()

	// 构建转账数据
	redeemTx, err := NewTx()
	if err != nil {
		return "",err
	}

	// 由UTXO添加inputs
	balanceFrom := int64(0)

	for i,_:=range self.FromList {
		for _, utxo := range self.FromList[i].Utxos {
			// 累加余额
			balanceFrom = balanceFrom + utxo.Value

			utxoHash, err := chainhash.NewHashFromStr(utxo.TxId)
			if err != nil {
				return "",err
			}

			prevOutPoint := wire.NewOutPoint(utxoHash, uint32(utxo.Index))
			txIn := wire.NewTxIn(prevOutPoint, nil, nil)
			redeemTx.AddTxIn(txIn)

			self.FromList[i].TxIndex =append(self.FromList[i].TxIndex, len(redeemTx.TxIn)-1)
		}
	}

	balanceTo := int64(0)
	for i,_:=range self.ToList {
		toAmount:=bctPrecisionConvert(self.ToList[i].Value)
		balanceTo=balanceTo+toAmount
		// 由地址生成address对象
		var toAddress btcutil.Address
		toAddress, err = btcutil.DecodeAddress(self.ToList[i].Address, self.Net)
		if err != nil {
			return "",err
		}
		// 添加to地址到output0
		toAddressP2PKHScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return "",err
		}

		redeemTxOut0 := wire.NewTxOut(toAmount, toAddressP2PKHScript)
		redeemTx.AddTxOut(redeemTxOut0)
	}

	feeRate, err := EstimateFee()
	if err != nil {
		return "",err
	}

	// 这里用的第3种方法
	size := len(redeemTx.TxIn)*180 + 2*34 + 10 + len(redeemTx.TxIn)
	fee := int64(math.Ceil(float64(int64(size) * feeRate)/1000))
	// 检查余额

	fmt.Println(balanceFrom,balanceTo,fee)
	if balanceFrom < balanceTo+fee {
		return "",errors.New("balanceFrom < balanceTo+fee")
	}

	// 把from添加到output, 找零
	_,fromPublicKey:=btcec.PrivKeyFromBytes(btcec.S256(),self.FromList[0].Pri)
	var fromAddressPubKey *btcutil.AddressPubKey
	fromAddressPubKey, err = btcutil.NewAddressPubKey(fromPublicKey.SerializeUncompressed(),self.Net)
	if err != nil {
		return "",err
	}

	var fromAddress btcutil.Address
	fromAddress, err = btcutil.DecodeAddress(fromAddressPubKey.EncodeAddress(), self.Net)
	if err != nil {
		return "",err
	}
	fromAddressP2PKHScript, err := txscript.PayToAddrScript(fromAddress)
	if err != nil {
		return "",err
	}

	redeemTxOut1 := wire.NewTxOut(balanceFrom-fee-balanceTo, fromAddressP2PKHScript)
	redeemTx.AddTxOut(redeemTxOut1)

	self.RedeemTx=redeemTx
	return "",nil
}