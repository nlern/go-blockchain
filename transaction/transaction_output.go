package transaction

import (
	"bytes"

	"github.com/nlern/go-blockchain/utils"
)

// TxOutput represents a transaction output
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock signs the output
func (out *TxOutput) Lock(address []byte) {
	decodedAddress := utils.Base58Decode(address)
	pubKeyHash := utils.GetPublicKeyHash(decodedAddress)
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner
// of the pubKey
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
