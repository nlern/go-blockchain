package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/wallets"
)

func (cli *CLI) createWallet() {
	ws, _ := wallets.NewWallets()
	address := ws.CreateWallet()
	if success, err := ws.SaveToFile(); success == false {
		log.Panic(err)
	}

	fmt.Printf("Your new address %q\n", address)
}
