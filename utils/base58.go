package utils

import (
	"log"

	"github.com/akamensky/base58"
)

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	encodedString := base58.Encode(input)

	return []byte(encodedString)
}

// Base58Decode decodes Base58 data
func Base58Decode(encoded []byte) []byte {
	encodedString := string(encoded)

	decodedString, err := base58.Decode(encodedString)	
	if err != nil {
		log.Panic(err)
	}

	return []byte(decodedString)
}
