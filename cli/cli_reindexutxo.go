package cli

import (
	"fmt"

	"github.com/cyprus09/blockchain/blockchainstruct"
)

func (cli *CLI) reindexUTXO() {
	bc := blockchainstruct.NewBlockchain()
	UTXOSet := blockchainstruct.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}