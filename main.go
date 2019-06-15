package main

import (
	"fmt"
	"strconv"

	"./blockchain"
	"./blockchain/proofofwork"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock("Send 1 BTC to Neo")
	bc.AddBlock("Send 2 more BTC to Neo")

	for _, block := range bc.Blocks {
		pow := proofofwork.NewProofOfWork(block)
		
		fmt.Printf("Prev. Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
