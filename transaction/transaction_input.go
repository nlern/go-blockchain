package transaction

import (
	"bytes"

	"github.com/nlern/go-blockchain/utils"
)

// TxInput represents a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks wheteher the address initiated the transaction
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := utils.HashPublicKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
