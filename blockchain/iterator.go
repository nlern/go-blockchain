package blockchain

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/nlern/go-blockchain/block"
)

// Iterator is used to iterate over blockchain blocks
type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns the next block from the blockchain
func (i *Iterator) Next() *block.Block {
	var currBlock *block.Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		encodeBlock := b.Get(i.currentHash)
		currBlock = block.DeserializeBlock(encodeBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = currBlock.PrevBlockHash

	return currBlock
}
