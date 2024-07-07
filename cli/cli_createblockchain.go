package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/blockchainstruct"
	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) createBlockchain(address string, nodeID) {
	if !wallets.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchainstruct.CreateBlockchain(address)
	defer bc.DB.Close()

	UTXOSet := blockchainstruct.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()
	
	fmt.Println("Done! Blockchain created successfully.")
}