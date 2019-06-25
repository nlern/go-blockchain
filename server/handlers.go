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
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

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
