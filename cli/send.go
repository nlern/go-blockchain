package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/transaction"
	txoperations "github.com/nlern/go-blockchain/transaction/operations"
	"github.com/nlern/go-blockchain/wallet"
)

func (cli *CLI) send(from, to string, amount int) {
	if wallet.ValidateAddress(from) == false {
		log.Panic("ERROR: Sender address is not valid")
	}
	if wallet.ValidateAddress(to) == false {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	fmt.Printf("Sending amount %d from %q to %q...\n\n", amount, from, to)

	tx := txoperations.NewUTXOTransaction(from, to, amount, bc)
	cbTx := transaction.NewCoinbaseTX(from, "")
	txs := []*transaction.Transaction{cbTx, tx}

	bc.MineBlock(txs)
	fmt.Printf("Successfully sent amount %d from %q to %q\n", amount, from, to)
}
