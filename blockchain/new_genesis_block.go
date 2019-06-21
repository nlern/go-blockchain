package blockchain

import (
	"github.com/nlern/go-blockchain/block"
	"github.com/nlern/go-blockchain/transaction"
)

// NewGenesisBlock creates and returns genesis block
func NewGenesisBlock(coinbase *transaction.Transaction) *block.Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}
