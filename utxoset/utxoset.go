package utxoset

import (
	"encoding/hex"
	"log"

	"github.com/nlern/go-blockchain/transaction"

	"github.com/boltdb/bolt"
	"github.com/nlern/go-blockchain/blockchain"
)

const utxoBucket = "chainstate"

// UTXOSet represents the UTXO set
type UTXOSet struct {
	blockchain *blockchain.Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accBalance := 0
	db := u.blockchain.GetDB()
	bucketName := []byte(utxoBucket)

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; cursor.Next() {
			txID := hex.EncodeToString(k)
			outs, err := transaction.DeserializeOutputs(v)
			if err != nil {
				return err
			}

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accBalance < amount {
					accBalance = accBalance + out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accBalance, unspentOutputs
}

// Reindex rebuilds the UTXOSet
func (u UTXOSet) Reindex() {
	db := u.blockchain.GetDB()
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucket(bucketName)
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
		bucket := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			err = bucket.Put(key, outs.Serialize())
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
