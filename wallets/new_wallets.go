package wallets

import "github.com/nlern/go-blockchain/wallet"

// NewWallets creates Wallets and fills it from a file if it exists
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*wallet.Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}
