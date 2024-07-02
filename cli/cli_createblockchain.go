package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/blockchainstruct"
	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) createBlockchain(address string) {
	if !wallets.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchainstruct.CreateBlockchain(address)
	bc.DB.Close()
	
	fmt.Println("Done!")
}