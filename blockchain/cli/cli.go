package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nlern/go-blockchain/blockchain"

	"github.com/nlern/go-blockchain/blockchain/proofofwork"
	"github.com/nlern/go-blockchain/blockchain/transaction"
	txoperations "github.com/nlern/go-blockchain/blockchain/transaction/operations"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createBlockchain(address string) {
	fmt.Printf("Creating a new blockchain...\n\n")

	bc := blockchain.CreateBlockchain(address)
	bc.CloseDB()

	fmt.Println("Successfully created blockchain!")
}

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	balance := 0
	UTXOs := bc.FindUTXOs(address)

	for _, utxo := range UTXOs {
		balance = balance + utxo.Value
	}

	fmt.Printf("Balance of %q : %d\n", address, balance)
}

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	iterator := bc.Iterate()

	fmt.Printf("Printing chain...\n\n")

	for {
		block := iterator.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := proofofwork.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if block.PrevBlockHash == nil {
			break
		}
	}

	fmt.Println("Chain printed")
}

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	fmt.Printf("Sending amount %d from %q to %q...\n\n", amount, from, to)

	tx := txoperations.NewUTXOTransaction(from, to, amount, bc)

	bc.MineBlock([]*transaction.Transaction{tx})
	fmt.Printf("Successfully sent amount %d from %q to %q\n", amount, from, to)
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Sender wallet address")
	sendTo := sendCmd.String("to", "", "Receiver wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount < 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
