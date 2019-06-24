package server

import (
	"fmt"
	"log"
	"net"

	"github.com/nlern/go-blockchain/blockchain"
)

const (
	protocol    = "tcp"
	nodeVersion = 1
)

var (
	nodeAddress,
	miningAddress string
	knownNodes = []string{"localhost:3000"}
)

// StartServer starts a node
func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := blockchain.NewBlockchain(nodeID)

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}
