package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/blockchainstruct"
	"github.com/cyprus09/blockchain/utils"
	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) getBalance(address string) {
	if !wallets.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := blockchainstruct.NewBlockchain(address)
	defer bc.DB.Close()

	balance := 0

	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", pubKeyHash, balance)
	fmt.Println()
}
