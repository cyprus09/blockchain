package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) listAddresses() {
	wallets, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}

	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
