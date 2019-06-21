package cli

import (
	"fmt"
	"strconv"

	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/proofofwork"
)

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	bci := bc.Iterator()

	fmt.Printf("Printing chain...\n\n")

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := proofofwork.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
