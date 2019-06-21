package utils

import (
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

// HashPublicKey double hashes public key (using sha256
// algorithm) and returns the hash
func HashPublicKey(publicKey []byte) []byte {
	pubSHA256 := sha256.Sum256(publicKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}
