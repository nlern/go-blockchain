package transaction

// NewTxOutput creates a new TxOutput
func NewTxOutput(value int, address string) *TxOutput{
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}