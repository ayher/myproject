package account

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pubUtils "myproject/public/utils"
)

func CreateAddress() (*ecdsa.PrivateKey,*ecdsa.PublicKey,error){
	pr,err:=ecdsa.GenerateKey(btcec.S256(),rand.Reader)
	if err!=nil{
		return nil,nil,err
	}
	return pr,&pr.PublicKey,nil
}

func GetAddressFromPublic(pu string) (*common.Address,error){
	b,err:=hex.DecodeString(pu)
	if err!=nil{
		return nil,err
	}
	publicKey,err:=btcec.ParsePubKey(b,btcec.S256())
	if err!=nil{
		return nil,err
	}
	pubKey:=publicKey.ToECDSA()
	compressedPKHexString := pubUtils.PubKeyToCompressedHexString(pubKey)

	recoveredPK, _ := pubUtils.PubKeyFromHexString(compressedPKHexString)
	recoveredPK.Curve = pubUtils.GetCurve()

	pubKey = recoveredPK

	fromAddress := crypto.PubkeyToAddress(*pubKey)
	return &fromAddress,nil
}
