package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI struct helps process command line arguments
type CLI struct{}

// printUsage prints all the available commands with their usage
func (cli *CLI) printUsage() {
	fmt.Println("Commands:")
	fmt.Println("")
	fmt.Println("  createblockchain -address <address>                                   : Create a blockchain and genesis block reward to address")
	fmt.Println("")
	fmt.Println("  printchain                                                            :  Displays all the blocks in the blockchain in order.")
	fmt.Println("")
	fmt.Println("  createwallet                                                          : Generates a new key-pair and saves it into the wallet file")
	fmt.Println("")
	fmt.Println("  listaddresses                                                         : Lists all addresses from the wallet file")
	fmt.Println("")
	fmt.Println("  getbalance -address <address>                                         : Get balance of address")
	fmt.Println("")
	fmt.Println("  reindexutxo                                                     : Rebuilds the UTXO set")
	fmt.Println("")
	fmt.Println("  sendcoin -from <from_address> -to <to_address> -amount <amount> -mine : Send amount of coins from from_address to to_address. Mine on the same node, when -mine is set.")
	fmt.Println("")
	fmt.Println(" startnode -miner <address> : Start a node with ID specified in nodeID env. var. -miner enables mining")
}

// validateArgs helps in validating the number of arguments within the cli
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run parses the cli arguments and processes these commands
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Panic("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	sendCoinCmd := flag.NewFlagSet("sendcoin", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCoinCmd.String("from", "", "Source wallet address")
	sendTo := sendCoinCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCoinCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCoinCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining node and send reward to <address>")

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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "sendcoin":
		err := sendCoinCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
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
		cli.getBalance(*getBalanceAddress, nodeID)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress, nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if sendCoinCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCoinCmd.Usage()
			os.Exit(1)
		}
		cli.sendCoin(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
	}

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}
}
