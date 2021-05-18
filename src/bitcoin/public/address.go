package public

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// 未压缩 p2pkh
func ImportPrikey(pri string) (string,string) {
	s,_:=hex.DecodeString(pri)
	priKey,pubKey:=btcec.PrivKeyFromBytes(btcec.S256(),s)

	privateKey := hex.EncodeToString(priKey.Serialize())
	publicKey := hex.EncodeToString(pubKey.SerializeUncompressed())
	return privateKey,publicKey
}

// 未压缩 p2pkh
func NewKey() (string,string) {
	priKey,err:=btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println(err)
	}
	privateKey := hex.EncodeToString(priKey.Serialize())
	publicKey := hex.EncodeToString(priKey.PubKey().SerializeUncompressed())
	return privateKey,publicKey
}

func PubkeyToAddress(pub string) string {
	b,_:=hex.DecodeString(pub)
	netParams := chaincfg.MainNetParams
	fromAddressPubKey, err := btcutil.NewAddressPubKey(b, &netParams)
	if err!=nil {
		fmt.Println(err)
	}
	return fromAddressPubKey.EncodeAddress()
}
