package main

import (
	"fmt"
)

func main() {
	bc := NewBlockChain()

	bc.AddBlock("Send 1 BTC to Mayank")
	bc.AddBlock("Send 2 BTC to Mayank")

	for _, block := range bc.blocks {
		fmt.Printf("Prev Hash Value: %x\n", block.PrevBlockHash)
		fmt.Printf("Block Data: %x\n", block.BlockData)
		fmt.Printf("Current Hash Value: %x\n", block.CurrHash)
		fmt.Println()
	}
}
