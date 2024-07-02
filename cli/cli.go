package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/cyprus09/blockchain/blockchainstruct"
)

// CLI struct helps process command line arguments
type CLI struct{}

func (cli *CLI) createBlockchain(address string) {
	bc := blockchainstruct.CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Blockchain creation successful.")
}

func (cli *CLI) getBalance(address string) {
	bc := blockchainstruct.NewBlockChain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", pubKeyHash, balance)
	fmt.Println()
}

// printUsage prints all the available commands with their usage
func (cli *CLI) printUsage() {
	fmt.Println("Commands:")
	fmt.Println("  getbalance -address <address>                                   : Get balance of address")
	fmt.Println("")
	fmt.Println("  createblockchain -address <address>                             : Create a blockchain and genesis block reward to address")
	fmt.Println("")
	fmt.Println("  printchain                                                      :  Displays all the blocks in the blockchain in order.")
	fmt.Println("")
	fmt.Println("  sendcoin -from <from_address> -to <to_address> -amount <amount> : Send amount of coins from from_address to to_address")
}

// validateArgs helps in validating the number of arguments within the cli
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// printChain iterates through the entire blockchain starting from the tip all the way to the start and prints the values
func (cli *CLI) printChain() {

	bc := NewBlockChain("")
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Current Block Hash: %x\n", block.CurrHash)
		pow := NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) sendCoin(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.db.Close()

	tx := NewUTXOTTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Printf("Success! Coins sent successfuly to '%s'", to)
	fmt.Println()
}

// Run parses the cli arguments and processes these commands
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCoinCmd := flag.NewFlagSet("sendcoin", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCoinCmd.String("from", "", "Source wallet address")
	sendTo := sendCoinCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCoinCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "sendcoin":
		err := sendCoinCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCoinCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCoinCmd.Usage()
			os.Exit(1)
		}
		cli.sendCoin(*sendFrom, *sendTo, *sendAmount)
	}
}
