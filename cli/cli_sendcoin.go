package cli

import (
	"fmt"
	"log"

	"github.com/cyprus09/blockchain/blockchainstruct"
	"github.com/cyprus09/blockchain/wallets"
)

func (cli *CLI) sendCoin(from, to string, amount int, nodeID string, mineNow bool) {
	if !wallets.ValidateAddress(from) {
		log.Panic("ERROR: Sender Address is not valid")
	}

	if !wallets.ValidateAddress(to) {
		log.Panic("ERROR: Recipient Address is not valid")
	}

	if to == from {
		log.Panic("ERROR: You cannot send coins to yourself")
	} else if from == to {
		log.Panic("ERROR: You cannot send coins to yourself")
	}

	bc := blockchainstruct.NewBlockchain(nodeID)
	UTXOSet := blockchainstruct.UTXOSet{Blockchain: bc}
	defer bc.DB.Close()

	wallets, err := wallets.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := blockchainstruct.NewUTXOTTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := blockchainstruct.NewCoinbaseTx(from, "")
		txs := []*blockchainstruct.Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Printf("Success sent %d coins from %s to %s\n", amount, from, to)
}
