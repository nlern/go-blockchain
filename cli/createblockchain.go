package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/utxoset"
	"github.com/nlern/go-blockchain/wallet"
)

func (cli *CLI) createBlockchain(address string) {
	if wallet.ValidateAddress(address) == false {
		log.Panic("ERROR: Address is not valid")
	}

	fmt.Println("Creating a new blockchain...")
	fmt.Println()

	bc := blockchain.CreateBlockchain(address)
	UTXOSet := utxoset.UTXOSet{Blockchain: bc}
	defer bc.CloseDB()

	UTXOSet.Reindex()

	fmt.Println("Successfully created blockchain!")
}
