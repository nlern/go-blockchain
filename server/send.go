package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/nlern/go-blockchain/transaction"

	blck "github.com/nlern/go-blockchain/block"
	"github.com/nlern/go-blockchain/blockchain"
	"github.com/nlern/go-blockchain/utils"
)

func sendBlock(address string, b *blck.Block) {
	data := block{address, b.Serialize()}
	payload, err := utils.Serialize(nil, data)
	if err != nil {
		log.Panic(err)
	}

	request := append(commandToBytes("block"), payload...)
	sendData(address, request)
}

func sendData(address string, data []byte) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		fmt.Printf("%s is not available\n", address)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != address {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendGetBlocks(address string) {
	payload, err := utils.Serialize(nil, getBlocks{address})
	if err != nil {
		log.Panic(err)
	}
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload, err := utils.Serialize(nil, &getData{address, kind, id})
	if err != nil {
		log.Panic(err)
	}

	request := append(commandToBytes("getdata"), payload...)
	sendData(address, request)
}

func sendInv(address string, kind string, items [][]byte) {
	inventory := inv{address, kind, items}

	payload, err := utils.Serialize(nil, inventory)
	if err != nil {
		log.Panic(err)
	}

	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendTx(address string, tnx *transaction.Transaction) {
	data := tx{address, tnx.Serialize()}

	payload, err := utils.Serialize(nil, data)
	if err != nil {
		log.Panic(err)
	}

	request := append(commandToBytes("tx"), payload...)

	sendData(address, request)
}

func sendVersion(address string, bc *blockchain.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload, err := utils.Serialize(nil, version{nodeVersion, bestHeight, nodeAddress})
	if err != nil {
		log.Panic(err)
	}

	request := append(commandToBytes("version"), payload...)

	sendData(address, request)
}
