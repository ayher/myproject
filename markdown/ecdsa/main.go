package main

import "C"
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"math/big"
	"myproject/public/fmt"
)


func main()  {
	c()
}

func c() error {
	//s1,_:=new(big.Int).SetString("31905475329698698180477949829636524701656852887399070638189431125636837201920",10)
	//s2,_:=new(big.Int).SetString("83886613907617497243093035179051383151180711391675833744415732015881324292417",10)
	//
	//fmt.Debug(new(big.Int).Add(s2,s1).String(),btcec.S256().N.String())
	ss,_:=new(big.Int).SetString("23543543",10)
	x1,y1:=btcec.S256().ScalarBaseMult(ss.Bytes())
	x2,y2:=btcec.S256().ScalarBaseMult(new(big.Int).Sub(btcec.S256().N,ss).Bytes())
	fmt.Error(x1,y1)
	fmt.Error(x2,y2)
	return nil



	/*
	[main.c:21] D 63481556249910665218183086135171669472351116862418108094404752369633882326524
	[main.c:22] X 107179828267535759119928587570055036604023070975664677263063377314101154829485
	[main.c:23] Y 66752936639974785450363218555146671568979457989628465565115352319396063423168
	*/

	p,_:=new(big.Int).SetString("63481556249910665218183086135171669472351116862418108094404752369633882326524",10)
	x,_:=new(big.Int).SetString("107179828267535759119928587570055036604023070975664677263063377314101154829485",10)
	y,_:=new(big.Int).SetString("66752936639974785450363218555146671568979457989628465565115352319396063423168",10)

	pr:=&ecdsa.PrivateKey{D:p,PublicKey:ecdsa.PublicKey{
		Curve:btcec.S256(),
		X:x,
		Y:y,
	}}
	fmt.Println(pr.D,pr.X,pr.Y)

	hashStr :="123"
	hash,_:=hex.DecodeString(hashStr)

	//r,s,err:=ecdsa.Sign(rand.Reader,pr,hash)
	//if err!=nil{
	//	return err
	//}

	// 37212900268878498885863964294118546528925284290187789966423253513035753701278 20100562003947792483683139206700719820816481223767626211732084303475948057717
	r,_:=new(big.Int).SetString("37212900268878498885863964294118546528925284290187789966423253513035753701278",10)
	s,_:=new(big.Int).SetString("20100562003947792483683139206700719820816481223767626211732084303475948057717",10)

	//s=new(big.Int).Sub(btcec.S256().N,s) // 半减
	//s=new(big.Int).Sub(btcec.S256().N,s)

	fmt.Println(r.String(),s.String())

	fmt.Println(Verify(&pr.PublicKey,hash,r,s))

	return nil
}


type combinedMult interface {
	CombinedMult(bigX, bigY *big.Int, baseScalar, scalar []byte) (x, y *big.Int)
}

type invertible interface {
	// Inverse returns the inverse of k in GF(P)
	Inverse(k *big.Int) *big.Int
}

func Verify(pub *ecdsa.PublicKey, hash []byte, r, s *big.Int) bool {
	// See [NSA] 3.4.2
	c := pub.Curve
	N := c.Params().N

	if r.Sign() <= 0 || s.Sign() <= 0 {
		return false
	}
	if r.Cmp(N) >= 0 || s.Cmp(N) >= 0 {
		return false
	}
	e := hashToInt(hash, c)

	var w *big.Int
	if in, ok := c.(invertible); ok {
		w = in.Inverse(s)
	} else {
		w = new(big.Int).ModInverse(s, N)
	}


	u1 := e.Mul(e, w)
	u1.Mod(u1, N)
	u2 := w.Mul(r, w)
	u2.Mod(u2, N)

	// Check if implements S1*g + S2*p
	var x, y *big.Int
	if opt, ok := c.(combinedMult); ok {
		x, y = opt.CombinedMult(pub.X, pub.Y, u1.Bytes(), u2.Bytes())
	} else {
		fmt.Debug(u1.String())
		x1, y1 := c.ScalarBaseMult(u1.Bytes())
		fmt.Debug(x1.String())
		x2, y2 := c.ScalarMult(pub.X, pub.Y, u2.Bytes())
		x, y = c.Add(x1, y1, x2, y2)
	}

	if x.Sign() == 0 && y.Sign() == 0 {
		return false
	}


	x.Mod(x, N)
	return x.Cmp(r) == 0
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}