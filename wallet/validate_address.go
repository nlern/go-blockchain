package wallet

import (
	"bytes"

	"github.com/nlern/go-blockchain/constants"
	"github.com/nlern/go-blockchain/utils"
)

// ValidateAddress checks if address is valid
func ValidateAddress(address string) bool {
	checksumLen := constants.AddressChecksumLen
	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChekSum := pubKeyHash[len(pubKeyHash)-checksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLen]
	targetCheckSum := utils.Checksum(append([]byte{version}, pubKeyHash...),checksumLen)

	return bytes.Compare(actualChekSum, targetCheckSum) == 0
}
