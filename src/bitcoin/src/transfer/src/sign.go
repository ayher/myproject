package src

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"myproject/public/fmt"
)

func (self *BtcTransfer)findUtxo(hash string,index int64) string {
	for _,v:=range self.FromList{
		for _,vv:=range v.Utxos{
			if vv.TxId==hash && vv.Index==index{
				return vv.PkScript
			}
		}
	}
	return ""
}

func (self *BtcTransfer)getPri(id int) []byte {
	for _,v:=range self.FromList{
		for _,vv:=range v.TxIndex {
			if vv==id{
				return v.Pri
			}
		}
	}
	return nil
}

func (self *BtcTransfer)Sign() error{
	// 循环inputs签名
	compress := false
	signHashType := txscript.SigHashAll

	for idx, v := range self.RedeemTx.TxIn {

		sourcePKScript, err := hex.DecodeString(self.findUtxo(v.PreviousOutPoint.Hash.String(),int64(v.PreviousOutPoint.Index)))
		if err != nil {
			return err
		}

		hash, err := txscript.CalcSignatureHash(sourcePKScript, signHashType, self.RedeemTx, idx)
		if err != nil {
			return err
		}

		pri,fromPublicKey:=btcec.PrivKeyFromBytes(btcec.S256(),self.getPri(idx))
		fmt.Println(hex.EncodeToString(pri.Serialize()),v.PreviousOutPoint.Hash.String())
		signature,e:=pri.Sign(hash)
		if e!=nil{
			return e
		}

		var pkData []byte
		if compress {
			pkData = fromPublicKey.SerializeCompressed()
		} else {
			pkData = fromPublicKey.SerializeUncompressed()
		}

		rawTxInSignature := append(signature.Serialize(), byte(signHashType))

		signatureScript, err := txscript.NewScriptBuilder().AddData(rawTxInSignature).AddData(pkData).Script()
		if err != nil {
			return err
		}

		// 添加签名
		self.RedeemTx.TxIn[idx].SignatureScript = signatureScript
	}

	// 总raw输出
	var signedTx bytes.Buffer
	err := self.RedeemTx.Serialize(&signedTx)
	if err != nil {
		return err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	self.HexSignedTx=hexSignedTx

	return nil
}