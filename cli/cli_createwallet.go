package cli

import (
	"fmt"

	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := wallets.NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
}