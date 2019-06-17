// Package transaction contains datastructures, methods related to blockchain
// transaction
package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction represents a blockchain transaction
type Transaction struct {
	// ID is transaction id
	ID []byte
	// Vin is array of transaction inputs
	Vin []TxInput
	// Vout is array of transaction outputs
	Vout []TxOutput
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// TxInput is a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TxOutput is a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// NewCoinbaseTX creates a new coinbase transaction for address `to`
// and with `data` as transaction data
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %q", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{subsidy, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}
