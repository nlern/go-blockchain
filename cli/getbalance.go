package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/wallet"

	"github.com/nlern/go-blockchain/utils"

	"github.com/nlern/go-blockchain/blockchain"
)

func (cli *CLI) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	balance := 0
	decodedAddress := utils.Base58Decode([]byte(address))
	pubKeyHash := utils.GetPublicKeyHash(decodedAddress)
	UTXOs := bc.FindUTXOs(pubKeyHash)

	for _, utxo := range UTXOs {
		balance = balance + utxo.Value
	}

	fmt.Printf("Balance of %q : %d\n", address, balance)
}
