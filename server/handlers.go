package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"

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
	case "getblocks":
		handleGetBlocks(request, bc)
	case "inv":
		handleInv(request, bc)
	case "version":
		handleVersion(request, bc)
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
