package blockchain

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/nlern/go-blockchain/transaction"
	dbutils "github.com/nlern/go-blockchain/utils/db"
)

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbutils.Exists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := dbutils.Open()

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
