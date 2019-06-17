/*
Package proofofwork implements blockchain's Proof-of-Work(PoW) functionality
*/
package proofofwork

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/nlern/go-blockchain/blockchain/block"
	"github.com/nlern/go-blockchain/blockchain/proofofwork/utils"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 17

/*
ProofOfWork type is datastructure for PoW
*/
type ProofOfWork struct {
	block  *block.Block
	target *big.Int
}

/*
NewProofOfWork creates a new PoW for a block.  It takes
pointer to the block and returns pointer to new PoW
*/
func NewProofOfWork(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

/*
Run executes PoW algorithm
*/
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}

		nonce = nonce + 1
	}

	fmt.Print("\n\n")

	return nonce, hash[:]
}

/*
Validate validates PoW
*/
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
