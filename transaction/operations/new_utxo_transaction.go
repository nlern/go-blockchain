package txoperations

import (
	"encoding/hex"
	"log"

	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/transaction"
	"github.com/nlern/go-blockchain/utils"
	"github.com/nlern/go-blockchain/wallets"
)

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *blockchain.Blockchain) *transaction.Transaction {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	ws, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}

	senderWallet := ws.GetWallet(from)
	senderPubKeyHash := utils.HashPublicKey(senderWallet.PublicKey)
	actualBalance, unspentOutputs := bc.FindSpendableOutputs(senderPubKeyHash, amount)

	if actualBalance < amount {
		log.Panic("ERROR: Sender has insufficient balance")
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
				Signature: nil,
				PubKey:    senderWallet.PublicKey,
			}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *transaction.NewTxOutput(amount, to))
	if actualBalance > amount {
		// sender's change
		outputs = append(outputs, *transaction.NewTxOutput(actualBalance-amount, from))
	}

	tx := transaction.Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetID()
	bc.SignTransaction(&tx, senderWallet.PrivateKey)

	return &tx
}
