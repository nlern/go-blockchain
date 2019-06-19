// Package blockchain contains datastructures, methods related to blockchain
package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nlern/go-blockchain/blockchain/transaction"

	"github.com/boltdb/bolt"
	"github.com/nlern/go-blockchain/blockchain/block"
	"github.com/nlern/go-blockchain/blockchain/proofofwork"
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

// FindUnspentTransactions returns a list of transactions containing
// unspent outputs for a given address
func (bc *Blockchain) FindUnspentTransactions(address string) []transaction.Transaction {
	var unspentTXs []transaction.Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterate()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if _, pres := spentTXOs[txID]; pres == true {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindUTXOs finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXOs(address string) []transaction.TxOutput {
	var UTXOs []transaction.TxOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs find and returns unspent outputs to reference in inputs
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	availableBalance := 0

FindUnspentOutputs:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if availableBalance >= amount {
				break FindUnspentOutputs
			}
			if out.CanBeUnlockedWith(address) {
				availableBalance = availableBalance + out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}

	return availableBalance, unspentOutputs
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

// NewBlockchain returns a new pointer to existing blockchain
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
