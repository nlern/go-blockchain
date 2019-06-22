package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/utxoset"
	"github.com/nlern/go-blockchain/wallet"

	"github.com/nlern/go-blockchain/utils"

	"github.com/nlern/go-blockchain/blockchain"
)

func (cli *CLI) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchain.NewBlockchain()
	UTXOSet := utxoset.UTXOSet{Blockchain: bc}
	defer bc.CloseDB()

	balance := 0

	fmt.Print("Decoding address...")
	decodedAddress := utils.Base58Decode([]byte(address))
	pubKeyHash := utils.GetPublicKeyHash(decodedAddress)
	fmt.Println(" Done")

	fmt.Print("Fetching balances...")
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	fmt.Println(" Fetched")

	for _, utxo := range UTXOs {
		balance = balance + utxo.Value
	}
	
	fmt.Printf("Balance of %q : %d\n", address, balance)
}
