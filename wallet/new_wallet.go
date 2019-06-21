package wallet

import "github.com/nlern/go-blockchain/utils"

// NewWallet creates and returns a new wallet
func NewWallet() *Wallet {
	privateKey, publicKey := utils.KeyPair()
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}
