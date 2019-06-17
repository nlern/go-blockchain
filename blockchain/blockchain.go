/*
Package blockchain contains datastructures, methods related to blockchain
*/
package blockchain

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nlern/go-blockchain/blockchain/block"
	"github.com/nlern/go-blockchain/blockchain/proofofwork"
)

const (
	dbfile       = "blockchain.db"
	blocksBucket = "blocks"
)

/*
Blockchain type represents a blokchain datastructure which
consists of an array of blockchain blocks
*/
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// Iterator is used to iterate over blockchain blocks
type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

/*
Next returns the next block from the blockchain
*/
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

/*
AddBlock method adds a new block to the blockchain.  It takes
as input the data for the new block
*/
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
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
Iterate instantiates a blockchain iterator and returns the iterator
*/
func (bc *Blockchain) Iterate() *Iterator {
	i := &Iterator{bc.tip, bc.db}

	return i
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
	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found, creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
