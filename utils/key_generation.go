package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

// KeyPair generates a new (privae key, public key) pair
func KeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publicKey := append(privateKey.X.Bytes(),
		privateKey.Y.Bytes()...)
	return *privateKey, publicKey
}
