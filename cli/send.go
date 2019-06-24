package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/utxoset"

	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/transaction"
	txoperations "github.com/nlern/go-blockchain/transaction/operations"
	"github.com/nlern/go-blockchain/wallet"
)

func (cli *CLI) send(from, to string, amount int) {
	fmt.Print("Verifying sender and receiver addresses...")
	if wallet.ValidateAddress(from) == false {
		log.Panic("ERROR: Sender address is not valid")
	}
	if wallet.ValidateAddress(to) == false {
		log.Panic("ERROR: Recipient address is not valid")
	}
	fmt.Println(" Done")
	bc := blockchain.NewBlockchain()
	UTXOSet := utxoset.UTXOSet{Blockchain: bc}
	defer bc.CloseDB()

	fmt.Print("Creating a new transaction...")
	tx := txoperations.NewUTXOTransaction(from, to, amount, &UTXOSet)
	fmt.Println(" Done")

	fmt.Print("Creating a new coinbase transaction...")
	cbTx := transaction.NewCoinbaseTX(from, "")
	fmt.Println(" Done")

	txs := []*transaction.Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)

	fmt.Print("Updating UTXO sets...")
	UTXOSet.Update(newBlock)
	fmt.Println(" Done")

	fmt.Println("Transaction successful!")
}
