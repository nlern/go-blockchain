/*
Package block contains datastructures, methods related to block of a blockchain
*/
package block

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

/*
Block is basic datastructure of blockchain.  It represents the block of
a blockchain
*/
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

/*
SetHash method computes the sha256 hash of the block headers, consisting of
block's previousBlockHsh, Data and Timestamp. It stores the hash in
the block Hash field
*/
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

/*
NewBlock function creates a new block.  It takes the data for block and its
previous block hash and returns pointer to the new block
*/
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()

	return block
}
