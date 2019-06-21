package utils

import (
	"github.com/nlern/go-blockchain/constants"
)

// GetPublicKeyHash returns public key hash without
// version and checksum data
func GetPublicKeyHash(pubKeyHashWithPayload []byte) []byte {
	minLen := constants.HashVersionLength
	maxLen := len(pubKeyHashWithPayload) - constants.AddressChecksumLen
	
	return pubKeyHashWithPayload[minLen:maxLen]
}
