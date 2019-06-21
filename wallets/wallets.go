package wallets

import (
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nlern/go-blockchain/utils"

	"github.com/nlern/go-blockchain/wallet"
)

const walletFile = "wallet.dat"

// Wallets store a collection of wallets
type Wallets struct {
	Wallets map[string]*wallet.Wallet
}

// CreateWallet adds a wallet to wallets
func (ws *Wallets) CreateWallet() string {
	newWallet := wallet.NewWallet()
	address := fmt.Sprintf("%s", newWallet.GetAddress())

	ws.Wallets[address] = newWallet

	return address
}

// GetAddresses returns an array of addresses stored in wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet returns a wallet by its address
func (ws *Wallets) GetWallet(address string) wallet.Wallet {
	return *ws.Wallets[address]
}

// LoadFromFile loads wallets from file
func (ws *Wallets) LoadFromFile() error {
	if exists, err := checkWalletFileExists(); exists == false {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets

	err = utils.Deserialize(elliptic.P256(), fileContent, &wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// SaveToFile save wallets to file and returns success or failure
// with error
func (ws *Wallets) SaveToFile() (bool, error) {
	serialized, err := utils.Serialize(elliptic.P256(), ws)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(walletFile, serialized, 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}
