package transaction

import (
	"log"

	"github.com/nlern/go-blockchain/utils"
)

// TxOutputs represents a collection of TxOutput
type TxOutputs struct {
	Outputs []TxOutput
}

// Serialize serializes TxOutputs
func (out *TxOutputs) Serialize() []byte {
	serialized, err := utils.Serialize(nil, out)
	if err != nil {
		log.Panic(err)
	}

	return serialized
}
