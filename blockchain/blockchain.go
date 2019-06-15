/*
Package blockchain contains datastructures, methods related to blockchain
*/
package blockchain

import (
	"time"

	"./block"
	"./proofofwork"
)

/*
Blockchain type represents a blokchain datastructure which
consists of an array of blockchain blocks
*/
type Blockchain struct {
	Blocks []*block.Block
}

/*
AddBlock method adds a new block to the blockchain.  It takes
as input the data for the new block
*/
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

/*
NewBlock function creates a new block.  It takes the data for block and its
previous block hash and returns pointer to the new block
*/
func NewBlock(data string, prevBlockHash []byte) *block.Block {
	block := &block.Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}
	pow := proofofwork.NewProofOfWork(block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

/*
NewGenesisBlock method generates the first block or Genesis block
of a blockchain and returns a pointer to the block
*/
func NewGenesisBlock() *block.Block {
	return NewBlock("Genesis block", []byte{})
}

/*
NewBlockchain method creates a new blockchain and returns a pointer
to the blockchain
*/
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{NewGenesisBlock()}}
}
