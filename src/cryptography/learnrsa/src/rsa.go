package src

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"github.com/farmerx/gorsa"
	"myproject/src/cryptography/keySaveInFile/src"
	"myproject/public/fmt"
)

// 获取公私钥，十六进制
func GetRsa() (pri string,pub string,err error) {
	Pri,err:=rsa.GenerateKey(rand.Reader,32)
	if err!=nil{
		return
	}

	priByte := x509.MarshalPKCS1PrivateKey(Pri)
	pri=hex.EncodeToString(priByte)
	err=src.Save("RSA PRIVATE KEY",pri,"./key/pricate.pem")
	if err!=nil{
		return
	}

	Pub := &Pri.PublicKey
	mpub, err := x509.MarshalPKIXPublicKey(Pub)
	if err!=nil{
		return
	}
	pub=hex.EncodeToString(mpub)
	err=src.Save("RSA PUBLIC KEY",pub,"./key/public.pem")
	if err!=nil{
		return
	}

	return
}

func Vitify(pri,pub string) (err error) {
	str:="ppp"

	err=gorsa.RSA.SetPrivateKey(pri)
	if err!=nil{
		fmt.Debug("12312312")
		return
	}
	err=gorsa.RSA.SetPublicKey(pub)
	if err!=nil{
		return
	}

	// vitify
	rb,err:=gorsa.RSA.PriKeyENCTYPT([]byte(str))
	if err!=nil{
		return
	}

	ub,err:=gorsa.RSA.PubKeyDECRYPT(rb)
	fmt.Println(string(ub))
	return
}