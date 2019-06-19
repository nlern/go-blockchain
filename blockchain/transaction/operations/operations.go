package operations

import (
	"encoding/hex"
	"log"

	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/blockchain/transaction"
)

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *blockchain.Blockchain) *transaction.Transaction {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	actualBalance, unspentOutputs := bc.FindSpendableOutputs(from, amount)

	if actualBalance < amount {
		log.Panic("ERROR: not enough balance in sender")
	}

	// Build a list of inputs
	for txid, outs := range unspentOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := transaction.TxInput{
				Txid:      txID,
				Vout:      out,
				ScriptSig: from,
			}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, transaction.TxOutput{
		Value:        amount,
		ScriptPubKey: to,
	})
	if actualBalance > amount {
		outputs = append(outputs, transaction.TxOutput{
			Value:        actualBalance - amount,
			ScriptPubKey: from,
		})
	}

	tx := transaction.Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetID()

	return &tx
}
