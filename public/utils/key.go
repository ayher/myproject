package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// 由 x,y 点转成pubKey
func PubKeyFromXY(x, y *big.Int) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: GetCurve(),
		X:     x,
		Y:     y,
	}
}

// 由 x,y string 转成pubKey
func PubKeyFromXYString(x string, y string) (*ecdsa.PublicKey, error) {
	X, ok := new(big.Int).SetString(x, 10)
	if !ok {
		return nil, fmt.Errorf("cannot convert x string to big int, string: %s", x)
	}
	Y, ok := new(big.Int).SetString(y, 10)
	if !ok {
		return nil, fmt.Errorf("cannot convert y string to big int, string: %s", y)
	}
	return PubKeyFromXY(X, Y), nil
}

// 由hexString转成pubKey，压缩未压缩均可
func PubKeyFromHexString(hexString string) (*ecdsa.PublicKey, error) {
	if len(hexString) == 128 {
		xByte, _ := hex.DecodeString(hexString[0:64])
		yByte, _ := hex.DecodeString(hexString[64:128])
		x := new(big.Int).SetBytes(xByte)
		y := new(big.Int).SetBytes(yByte)
		recoveredPubKey, _ := PubKeyFromXYString(x.String(), y.String())
		return recoveredPubKey, nil
	} else if len(hexString) == 66 {
		recoveredPKByte, err := hex.DecodeString(hexString)
		if err != nil {
			return nil, err
		}
		recoveredPubKey, err := DecompressPubKey(recoveredPKByte)
		if err != nil {
			return nil, err
		}
		return recoveredPubKey, nil
	}
	return nil, fmt.Errorf("hex string length error")
}


// PubKeyToHexString encodes a public key to the 128-length string.
func PubKeyToHexString(pubKey *ecdsa.PublicKey) string {
	pubKeyBytes := append([]byte{}, common.LeftPadBytes(pubKey.X.Bytes(), 32)...)
	pubKeyBytes = append(pubKeyBytes, common.LeftPadBytes(pubKey.Y.Bytes(), 32)...)
	pubKeyHexString := hex.EncodeToString(pubKeyBytes)
	return pubKeyHexString
}

// PubKeyToCompressedHexString encodes a public key to the 66-length compressed string.
func PubKeyToCompressedHexString(pubKey *ecdsa.PublicKey) string {
	compressedPKByte := CompressPubKey(pubKey)
	compressedPKHexString := hex.EncodeToString(compressedPKByte)
	return compressedPKHexString
}

// CompressPubKey encodes a public key to the 33-byte compressed format.
func CompressPubKey(pubKey *ecdsa.PublicKey) []byte {
	return secp256k1.CompressPubkey(pubKey.X, pubKey.Y)
}

// DecompressPubKey parses a public key in the 33-byte compressed format.
func DecompressPubKey(pubKey []byte) (*ecdsa.PublicKey, error) {
	x, y := secp256k1.DecompressPubkey(pubKey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: x, Y: y, Curve: GetCurve()}, nil
}

// GetCurve returns the curve we used in this example.
func GetCurve() elliptic.Curve {
	// For simplicity, we use S256 curve.
	return btcec.S256()
}
