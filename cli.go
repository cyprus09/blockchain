package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI struct helps process command line arguments
type CLI struct {}

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Blockchain creation successful.")
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Println("Balance of '%s': %d\n", address, balance)
}

// printUsage prints all the available commands with their usage
func (cli *CLI) printUsage() {
	fmt.Println("Commands:")
	fmt.Println("  getbalance -address <address> : Get balance of address")
	fmt.Println("")
	fmt.Println("  createblockchain -address <address> - Create a blockchain and genesis block reward to address")
	fmt.Println("")
	fmt.Println("  printchain                   :  Displays all the blocks in the blockchain in order.")
	fmt.Println("")
	fmt.Println()
}

// validateArgs helps in validating the number of arguments within the cli
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// addBlock implements the cli command for adding block through the cli
func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Block added successfully!")
}

// printChain iterates through the entire blockchain starting from the tip all the way to the start and prints the values
func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("BlockData: %s\n", block.BlockData)
		fmt.Printf("Current Block Hash: %x\n", block.CurrHash)
		pow := NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}


// Run parses the cli arguments and processes these commands
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
