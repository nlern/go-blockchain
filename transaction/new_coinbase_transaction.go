package transaction

import (
	"crypto/rand"
	"fmt"
	"log"
)

// NewCoinbaseTX creates a new coinbase transaction for address `to`
// and with `data` as transaction data
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(subsidy, to)
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()

	return &tx
}
