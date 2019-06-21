package wallet

import (
	"crypto/ecdsa"

	"github.com/nlern/go-blockchain/constants"
	"github.com/nlern/go-blockchain/utils"
)

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// GetAddress returns wallet address
func (w *Wallet) GetAddress() []byte {
	pubKeyHash := utils.HashPublicKey(w.PublicKey)

	versionedPayload := append([]byte{constants.HashVersion}, pubKeyHash...)
	checksum := utils.Checksum(versionedPayload, constants.AddressChecksumLen)

	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	return []byte(address)
}
