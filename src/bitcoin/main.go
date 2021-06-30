package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
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

//8353c1e67ba0dc3061dd1220a5ebbc5ad8cbc5fdb2c7f4d428696ab0a9ea2176
//0469a1c7cd7bc162b98c98533b024d5a0c6fdb55ffa508e5a9f22f9a659409c850fccc2ec15549cfd2583c61247dd71123ec9becc180deecd6d5bd97e7f2f7dfda
//14EHU1tLH2jFJCTq47qWSkFLD6172k2FdY
//mikEm4yK64AW5JwSmgotGfTf55bovEuusq


//[main.main:21] cc9ecfd789b203e441fbfba8bd875461bd5a2eca91fc89b8ee806df4ec9ff649
//[main.main:22] 0415d39902ea6421d9c8dfba1d584215cabaebd8c08fc3307793b8aa3c02a775889aa526332d3fe8cf1ab288ba1c331855817e59da02d31db9c8f63d3a531ed4a7
//[main.main:23] 1AqwTk1hVbhNJxMkD3kAYKfrdvF7XUxeTR
//mqMtko6gJd8d64qMvciYNEtBVuqpNgt8EQ

//[main.main:28] c0432ce2f4a189ae2b50c13cd6a2f033d83361b3747f70ac2d0907edbf46aa51
//[main.main:29] 041358615d7182b5c1be5a00395bbf591f7060dccc6c215ae93758f42137a1dfa636b13e65a44f26641c87c9d14acb25d0a01dc6549a9f51bfa374a6dc4a578720
//[main.main:30] 1GbQZyDmBSbYCfeTqkuYHzRodp1gpeBb6U
//mw7Ms2JjzU2nyn85ZKsv7ue8VocPmE5aQ8


type UTXO struct {
	TxId     string `json:"tx_hash"`
	Index    int64  `json:"tx_output_n"`
	Value    int64  `json:"value"`
	PkScript string `json:"script"`
	Spend    bool   `json:"spent"`
}

func main()  {
	var err error
	b,_:=hex.DecodeString("41932a6271f15818af221ab9f6ad79949f3e742e070eb8fa278e83858bd1bde7")
	pri,fromPublicKey:=btcec.PrivKeyFromBytes(btcec.S256(),b)



	// 网络
	//netParams := chaincfg.MainNetParams
	netParams:=chaincfg.Params{}

	netParams.PrivateKeyID=33
	netParams.PubKeyHashAddrID=55
	netParams.ScriptHashAddrID=80

	//fmt.Println("fromPublicKey", "X", fromPublicKey.X)
	//fmt.Println("fromPublicKey", "Y", fromPublicKey.Y)
	//fmt.Println("fromPublicKey", "Hex", fmt.Sprintf("%x", fromPublicKey.SerializeCompressed()))


	// to地址
	destination := "PNDyeKXLSsfgD182soXoHas4rBJC5TZuAs"
	var amount int64
	amount=123120


	// 由公钥xy生成AddressPubKey对象
	var fromAddressPubKey *btcutil.AddressPubKey
	fromAddressPubKey, err = btcutil.NewAddressPubKey(fromPublicKey.SerializeCompressed(), &netParams)
	if err != nil {
		fmt.Println("btcutil.NewAddressPubKey", "err", err)
	}
	//fmt.Println()("NewAddressPubKey", "String", fromAddressPubKey.AddressPubKeyHash().String())
	fmt.Println("From", "EncodeAddress", fromAddressPubKey.EncodeAddress()) // 转账地址

	// 由地址生成address对象
	var toAddress btcutil.Address
	toAddress, err = btcutil.DecodeAddress(destination, &netParams)
	if err != nil {
		fmt.Println("btcutil.DecodeAddress", "err", err)
	}
	//fmt.Println()("DecodeAddress", "String", toAddress.String())
	fmt.Println("To", "EncodeAddress", toAddress.EncodeAddress()) // 转账地址

	// 获取所有input的UTXO
	var UTXOs []*UTXO
	UTXOs, err = GetUTXOs(fromAddressPubKey.EncodeAddress())

	// 构建转账数据
	redeemTx, err := NewTx()
	if err != nil {
		fmt.Println("NewTx", "err", err)
	}

	fmt.Println("UTXOs:",UTXOs)

	// 由UTXO添加inputs
	balance := int64(0)
	//scriptHashFill := make([]byte, 35)
	for i, utxo := range UTXOs {

		// 累加余额
		balance = balance + utxo.Value

		utxoHash, err := chainhash.NewHashFromStr(utxo.TxId)
		if err != nil {
			fmt.Println("chainhash.NewHashFromStr", "err", err)
		}
		//fmt.Println("input", i, utxoHash)
		pkScript, err := hex.DecodeString(utxo.PkScript)
		if err != nil {
			fmt.Println("Pkscript", "hex.DecodeString err", err)
		}
		if disasm, err := txscript.DisasmString(pkScript); err != nil {
			fmt.Println("Pkscript", "txscript.DisasmString err", err)
		} else {
			fmt.Println("Pkscript", `i`, i, `disasm`, disasm)
		}

		prevOutPoint := wire.NewOutPoint(utxoHash, uint32(utxo.Index))
		txIn := wire.NewTxIn(prevOutPoint, nil, nil)
		redeemTx.AddTxIn(txIn)

		//// 预填空签名，用于后面计算转账信息长度，用来算手续费
		//redeemTx.TxIn[i].SignatureScript = scriptHashFill
	}

	// 估算手续费
	// 获取手续费率，每kb多少聪
	feeRate, err := EstimateFee()
	fmt.Println("EstimateFee", "feeRate", feeRate)
	if err != nil {
		fmt.Println("EstimateFee", "err", err)
	}

	// 这里用的第3种方法
	size := len(redeemTx.TxIn)*180 + 2*34 + 10 + len(redeemTx.TxIn)
	fmt.Println("EstimateFee", "bytes", size)
	fee := int64(math.Ceil(float64(int64(size) * feeRate)/1000))
	fmt.Println("EstimateFee", "fee", fee)

	fmt.Println("balance < amount+fee",balance , amount,fee)
	// 检查余额
	if balance < amount+fee {
		fmt.Println("balance check", "err", "the balance of the account is not sufficient")
	}
	fmt.Println("tx", "amount", amount)

	// 添加outputs

	// 添加to地址到output0
	toAddressP2PKHScript, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		fmt.Println("destination txscript.PayToAddrScript", "err", err)
	}
	fmt.Println("%x", toAddressP2PKHScript)

	redeemTxOut0 := wire.NewTxOut(amount, toAddressP2PKHScript)
	if disasm, err := txscript.DisasmString(toAddressP2PKHScript); err != nil {
		fmt.Println("Pkscript", "DisasmString err", err)
	} else {
		fmt.Println("out0", "Value", redeemTxOut0.Value, "PkScript", disasm)
	}
	redeemTx.AddTxOut(redeemTxOut0)

	// 把from添加到output1, 找零
	var fromAddress btcutil.Address
	fromAddress, err = btcutil.DecodeAddress(fromAddressPubKey.EncodeAddress(), &netParams)
	if err != nil {
		fmt.Println("btcutil.DecodeAddress", "err", err)
	}
	fromAddressP2PKHScript, err := txscript.PayToAddrScript(fromAddress)
	if err != nil {
		fmt.Println("from txscript.PayToAddrScript", "err", err)
	}
	fmt.Println("%x", fromAddressP2PKHScript)

	redeemTxOut1 := wire.NewTxOut(balance-fee-amount, fromAddressP2PKHScript)
	if disasm, err := txscript.DisasmString(fromAddressP2PKHScript); err != nil {
		fmt.Println("Pkscript", "DisasmString err", err)
	} else {
		fmt.Println("out1", "Value", redeemTxOut1.Value, "PkScript", disasm)
	}
	redeemTx.AddTxOut(redeemTxOut1)





	// 循环inputs签名
	compress := false
	signHashType := txscript.SigHashAll
	for idx, _ := range redeemTx.TxIn {
		// 计算摘要
		sourcePKScript, err := hex.DecodeString(UTXOs[idx].PkScript)
		if err != nil {
			fmt.Println("Pkscript", "DisasmString err", err)
		}

		hash, err := txscript.CalcSignatureHash(sourcePKScript, signHashType, redeemTx, idx)
		if err != nil {
			fmt.Println("sign", "txscript.CalcSignatureHash err", err)
		}
		digestString := hex.EncodeToString(hash)
		fmt.Println("sign", "digest", digestString)

		// 这里用门限签把digest签名得到r、s
		//var i string
		//fmt.Println("r_s_v:")
		//_, _ = fmt.Scanln(&i)
		//tmp := strings.Split(i, "_")
		//if len(tmp) != 3 {
		//	fmt.Println("Sign", "err", "invalid signer result")
		//}
		//R, _ := new(big.Int).SetString(tmp[0], 10)
		//S, _ := new(big.Int).SetString(tmp[1], 10)
		signature,e:=pri.Sign(hash)
		if e!=nil{
			fmt.Println(e)
		}
		//signature := &btcec.Signature{
		//	R: R,
		//	S: S,
		//}

		var pkData []byte
		if compress {
			pkData = fromPublicKey.SerializeCompressed()
		} else {
			pkData = fromPublicKey.SerializeCompressed()
		}

		rawTxInSignature := append(signature.Serialize(), byte(signHashType))

		signatureScript, err := txscript.NewScriptBuilder().AddData(rawTxInSignature).AddData(pkData).Script()
		if err != nil {
			fmt.Println("sign", "txscript.signatureScript err", err)
		}

		// 添加签名
		redeemTx.TxIn[idx].SignatureScript = signatureScript

	}
	fmt.Println(redeemTx)
	fmt.Println("hash",redeemTx.TxIn[0].PreviousOutPoint.Hash.String())
	bb,_:=json.Marshal(redeemTx)
	fmt.Println("ppppppppp",string(bb))
	// 总raw输出
	var signedTx bytes.Buffer
	err = redeemTx.Serialize(&signedTx)
	if err != nil {
		fmt.Println("sign", "redeemTx.Serialize", err)
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	fmt.Println("Final", "raw signed transaction", hexSignedTx)
	fmt.Println("Final", "Decode", "https://live.blockcypher.com/btc/decodetx/")
	fmt.Println("Final", "Broadcast", "https://live.blockcypher.com/btc/pushtx")
}

type ChainInfo struct {
	Name           string    `json:"name"`
	Height         uint32    `json:"height"`
	Time           time.Time `json:"time"`
	HighFeePerKb   int64     `json:"high_fee_per_kb"`
	MediumFeePerKb int64     `json:"medium_fee_per_kb"`
	LowFeePerKb    int64     `json:"low_fee_per_kb"`
}

func getBlockChainInfo() (*ChainInfo, error) {
	// https://api.blockcypher.com/v1/btc/test3

	newURL := `http://api.blockcypher.com/v1/btc/test3`

	response, err := http.Get(newURL)
	if err != nil {
		fmt.Println("error in EstimateFee, http.Get")
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var resp *ChainInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("error in EstimateFee, json.Unmarshal")
		return nil, err
	}

	fmt.Println("getBlockChainInfo", "info", resp)

	return resp, nil
}

func EstimateFee() (int64, error) {
	info, err := getBlockChainInfo()
	if err != nil {
		return 0, err
	}
	return info.HighFeePerKb, nil
}

// 根据api返回地址多个UTXO
func GetUTXOs(address string) ([]*UTXO, error) {
	return []*UTXO{
		{
			TxId:"12639765559e7e5154e4f7eaaab14a823d59d1d4f5e14c8c392e5378530a3bc7",
			Index:1,
			Value:545999000,
			PkScript:"76a9146eef1c01044cf0aef8557438f5391d9f7b4e928388ac",
			Spend:false,
		},
	},nil
	// https://api.blockcypher.com/v1/btc/test3/addrs/n4TBTn6xNDjNvqoWoaaJWi33dY4sMQ2DqV??unspentOnly=true&includeScript=true&limit=50

	newURL := fmt.Sprintf("https://api.blockcypher.com/v1/btc/test3/addrs/%s??unspentOnly=true&includeScript=true&limit=50", address)

	response, err := http.Get(newURL)
	if err != nil {
		fmt.Println("error in GetUTXOs, http.Get")
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var resp *UnSpent
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("error in GetUTXOs, json.Unmarshal")
		return nil, err
	}

	// 过滤有效txrefs
	var utxos []*UTXO
	balance := int64(0)
	for i, u := range resp.UTXOs {
		if u.Index >= 0 && u.Spend == false {
			utxos = append(utxos, u)
			balance += u.Value
			fmt.Println("GetUTXOs", "i", i, "UTXO", u)
		}
	}

	return utxos, nil
}

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

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}
