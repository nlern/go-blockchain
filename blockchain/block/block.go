/*
Package block contains datastructures, methods related to block of a blockchain
*/
package block

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
