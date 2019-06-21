package wallets

import "os"

func checkWalletFileExists() (bool, error) {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
