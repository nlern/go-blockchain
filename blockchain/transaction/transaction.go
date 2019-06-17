// Package transaction contains datastructures, methods related to blockchain
// transaction
package transaction

// Transaction represents a blockchain transaction
type Transaction struct {
	// ID is transaction id
	ID []byte
	// Vin is array of transaction inputs
	Vin []TxInput
	// Vout is array of transaction outputs
	Vout []TxOutput
}

// TxInput is a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TxOutput is a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}
