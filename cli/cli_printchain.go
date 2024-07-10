package cli

import (
	"fmt"
	"strconv"

	"github.com/cyprus09/blockchain/blockchainstruct"
)

// printChain iterates through the entire blockchain starting from the tip all the way to the start and prints the values
func (cli *CLI) printChain(nodeID string) {

	bc := blockchainstruct.NewBlockchain(nodeID)
	defer bc.DB.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============= Block %x =============\n", block.CurrHash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. Block: %x\n", block.PrevBlockHash)
		pow := blockchainstruct.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}

		fmt.Println()
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
