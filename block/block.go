/*
Package block contains datastructures, methods related to block of a blockchain
*/
package block

import (
	"bytes"
	"crypto/sha256"
	"log"

	"github.com/nlern/go-blockchain/utils"

	"github.com/nlern/go-blockchain/transaction"
)

/*
Block is basic datastructure of blockchain.  It represents the block of
a blockchain
*/
type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// Serialize serializes block structure into byte array and returns slice
// of the array
func (b *Block) Serialize() []byte {
	serialized, err := utils.Serialize(nil, b)

	if err != nil {
		log.Fatal(err)
	}

	return serialized
}
