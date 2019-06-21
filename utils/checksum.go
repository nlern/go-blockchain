package utils

import "crypto/sha256"

// Checksum returns checksum for given payload and address length
func Checksum(payload []byte, addressLen int) []byte {
	firstSHA := sha256.Sum224(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressLen]
}
