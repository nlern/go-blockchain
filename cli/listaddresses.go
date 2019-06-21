package cli

import (
	"fmt"
	"log"

	"github.com/nlern/go-blockchain/wallets"
)

func (cli *CLI) listAddresses() {
	ws, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}

	addresses := ws.GetAddresses()

	fmt.Println("Printing addresses...")
	fmt.Println()
	
	for _, address := range addresses {
		fmt.Println(address)
	}
}
