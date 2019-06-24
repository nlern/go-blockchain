package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/nlern/go-blockchain/utils"

	"github.com/nlern/go-blockchain/blockchain"
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
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleGetBlocks(request []byte, bc *blockchain.Blockchain) {
	var payload version
	err := utils.Deserialize(nil, request, &payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleVersion(request []byte, bc *blockchain.Blockchain) {
	var payload version
	err := utils.Deserialize(nil, request, &payload)
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
