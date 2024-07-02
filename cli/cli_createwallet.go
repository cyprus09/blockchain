package cli

import (
	"fmt"

	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) createWallet() {
	wallets, _ := wallets.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}