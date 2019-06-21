package transaction

import (
	"fmt"
)

// NewCoinbaseTX creates a new coinbase transaction for address `to`
// and with `data` as transaction data
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %q", to)
	}

	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(subsidy, to)
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()

	return &tx
}
