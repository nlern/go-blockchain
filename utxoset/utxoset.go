package utxoset

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
	"github.com/nlern/go-blockchain/blockchain"
)

const utxoBucket = "chainstate"

// UTXOSet represents the UTXO set
type UTXOSet struct {
	blockchain *blockchain.Blockchain
}

// Reindex rebuilds the UTXOSet
func (u UTXOSet) Reindex() {
	db := u.blockchain.GetDB()
	bucket := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucket(bucket)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			err = b.Put(key, outs.Serialize())
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
