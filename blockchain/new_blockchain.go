package blockchain

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	dbutils "github.com/nlern/go-blockchain/utils/db"
)

// NewBlockchain returns a new pointer to existing blockchain
func NewBlockchain(nodeID string) *Blockchain {
	if dbutils.Exists(nodeID) == false {
		fmt.Println("No existing blockchain found, create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := dbutils.Open(nodeID)

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
