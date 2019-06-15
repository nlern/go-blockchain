/*
Package blockchain contains datastructures, methods related to blockchain
*/
package blockchain

import "./block"

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
	newBlock := block.NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

/*
NewGenesisBlock method generates the first block or Genesis block
of a blockchain and returns a pointer to the block
*/
func NewGenesisBlock() *block.Block {
	return block.NewBlock("Genesis block", []byte{})
}

/*
NewBlockchain method creates a new blockchain and returns a pointer
to the blockchain
*/
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{NewGenesisBlock()}}
}
