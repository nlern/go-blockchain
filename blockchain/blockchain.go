/*
Package blockchain contains datastructures, methods related to blockchain
*/
package blockchain

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nlern/blockchain_go/blockchain/transaction"

	"github.com/boltdb/bolt"
	"github.com/nlern/blockchain_go/blockchain/block"
	"github.com/nlern/blockchain_go/blockchain/proofofwork"
)

const (
	dbfile              = "blockchain.db"
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain implements interaction with a DB
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

func dbExists() bool {
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		return false
	}
	return true
}

// MineBlock mines a new block with provided transactions
func (bc *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

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

// NewBlock creates and returns transaction
func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *block.Block {
	block := &block.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
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

// CloseDB closes the blockchain db
func (bc *Blockchain) CloseDB() {
	bc.db.Close()
}

// NewGenesisBlock creates and returns genesis block
func NewGenesisBlock(coinbase *transaction.Transaction) *block.Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain() *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found, create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		genesisTx := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(genesisTx)

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

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
