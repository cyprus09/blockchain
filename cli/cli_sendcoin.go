package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/blockchainstruct"
	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) sendCoin(from, to string, amount int) {
	if !wallets.ValidateAddress(from) {
		log.Panic("ERROR: Sender Address is not valid")
	}

	if !wallets.ValidateAddress(to) {
		log.Panic("ERROR: Recipient Address is not valid")
	}

	bc := blockchainstruct.NewBlockchain()
	UTXOSet := blockchainstruct.UTXOSet{Blockchain: bc}
	defer bc.DB.Close()

	tx := blockchainstruct.NewUTXOTTransaction(from, to, amount, &UTXOSet)
	cbTx := blockchainstruct.NewCoinbaseTx(from, "")
	txs := []*blockchainstruct.Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	
	fmt.Printf("Success sent %d coins from %s to %s\n", amount, from, to)
}