package block

import (
	"log"

	"github.com/nlern/go-blockchain/utils"
)

// DeserializeBlock deserializes an encoded block and returns decoded block
func DeserializeBlock(d []byte) *Block {
	var block Block

	err := utils.Deserialize(nil, d, &block)

	if err != nil {
		log.Fatal(err)
	}

	return &block
}
