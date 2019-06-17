package main

import (
	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/blockchain/cli"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	cli := cli.CLI{BC: bc}
	cli.Run()
}
