package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/nlern/go-blockchain/utxoset"

	blck "github.com/nlern/go-blockchain/block"

	"github.com/nlern/go-blockchain/transaction"

	"github.com/nlern/go-blockchain/utils"

	"github.com/nlern/go-blockchain/blockchain"
)

var (
	blocksInTransit = [][]byte{}
	memPool         = make(map[string]transaction.Transaction)
)

func handleBlock(request []byte, bc *blockchain.Blockchain) {
	var payload block
	err := utils.Deserialize(nil, request[commandLength:], &payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blck.DeserializeBlock(blockData)

	fmt.Println("Received a new block")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := utxoset.UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()
	}
}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "block":
		handleBlock(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "inv":
		handleInv(request, bc)
	case "version":
		handleVersion(request, bc)
	case "tx":
		handleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleGetBlocks(request []byte, bc *blockchain.Blockchain) {
	var payload version
	err := utils.Deserialize(nil, request[commandLength:], &payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *blockchain.Blockchain) {
	var payload getData
	err := utils.Deserialize(nil, request[commandLength:], &payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memPool[txID]

		sendTx(payload.AddrFrom, &tx)
	}
}

func handleInv(request []byte, bc *blockchain.Blockchain) {
	var payload inv
	err := utils.Deserialize(nil, request[commandLength:], &payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, block := range blocksInTransit {
			if bytes.Compare(block, blockHash) != 0 {
				newInTransit = append(newInTransit, block)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if _, present := memPool[hex.EncodeToString(txID)]; present == false {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleVersion(request []byte, bc *blockchain.Blockchain) {
	var payload version
	err := utils.Deserialize(nil, request[commandLength:], &payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if nodeIsKnown(payload.AddrFrom) == false {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func handleTx(request []byte, bc *blockchain.Blockchain) {
	var payload tx
	var tx transaction.Transaction
	err := utils.Deserialize(nil, request[commandLength:], &payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	err = utils.Deserialize(nil, txData, &tx)
	if err != nil {
		log.Panic(err)
	}
	memPool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(memPool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*transaction.Transaction

			for id := range memPool {
				tx := memPool[id]
				if bc.VerifyTransaction(&tx) == true {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid, waiting for new ones...")
				return
			}

			cbTx := transaction.NewCoinbaseTX(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := utxoset.UTXOSet{Blockchain: bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(memPool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(memPool) > 0 {
				goto MineTransactions
			}
		}
	}
}
