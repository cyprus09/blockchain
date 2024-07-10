package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) startNode(nodeID, minerAddess string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddess) > 0 {
		if wallets.ValidateAddress(minerAddess) {
			fmt.Println("Mining is on. Address to recieve rewards: ", minerAddess)
		}else {
			log.Panic("Wrong miner address!")
		}
	}
	StartServer(nodeID, minerAddess)
}