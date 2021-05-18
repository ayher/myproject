package src

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type Tf struct {
	Pri     []byte
	pub     []byte
	Address string
	Value   float64
	Utxos   []*UTXO
	TxIndex []int
}

type UTXO struct {
	TxId     string `json:"tx_hash"`
	Index    int64  `json:"tx_output_n"`
	Value    int64  `json:"value"`
	PkScript string `json:"script"`
	Spend    bool   `json:"spent"`
}

type BtcTransfer struct {
	FromList []Tf
	ToList   []Tf
	Net      *chaincfg.Params
	RedeemTx *wire.MsgTx
	HexSignedTx string
}

var precision=100000000

func bctPrecisionConvert(f float64) int64 {
	return int64(f*float64(precision))
}