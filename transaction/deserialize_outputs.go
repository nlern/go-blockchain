package transaction

import (
	"github.com/nlern/go-blockchain/utils"
)

// DeserializeOutputs deserializes a serialized TxOutputs and returns
// deserialized TxOutputs
func DeserializeOutputs(serializedOuts []byte) (*TxOutputs, error) {
	txOuts := &TxOutputs{}
	err := utils.Deserialize(nil, serializedOuts, txOuts)

	return txOuts, err
}
