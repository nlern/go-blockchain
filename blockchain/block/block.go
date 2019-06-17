/*
Package block contains datastructures, methods related to block of a blockchain
*/
package block

import (
	"bytes"
	"encoding/gob"
	"log"
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
	Nonce         int
}

/*
Serialize serializes block structure into byte array and returns slice
of the array
*/
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	if err != nil {
		log.Fatal(err)
	}

	return result.Bytes()
}

/*
DeserializeBlock deserializes an encoded block and returns decoded block
*/
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err != nil {
		log.Fatal(err)
	}

	return &block
}
