package blockchain

import (
	"time"

	"github.com/nlern/go-blockchain/block"
	"github.com/nlern/go-blockchain/proofofwork"
	"github.com/nlern/go-blockchain/transaction"
)

// NewBlock creates and returns a new block contaiing transactions
func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte, height int) *block.Block {
	block := &block.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Height:        height,
	}
	pow := proofofwork.NewProofOfWork(block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
