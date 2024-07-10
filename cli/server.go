package cli

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/cyprus09/blockchain/blockchainstruct"
)

const (
	protocol    = "tcp"
	nodeVersion = 1
	commandLen  = 12
)

var (
	nodeAddress     string
	miningAddress   string
	knownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	memPool         = make(map[string]blockchainstruct.Transaction)
)

type address struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type tx struct {
	AddrFrom    string
	Transaction []byte
}

type version struct {
	AddrFrom   string
	Version    int
	BestHeight int
}

func commandToBytes(command string) []byte {
	var bytes [commandLen]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLen]
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendAddress(addr string) {
	nodes := address{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("address"), payload...)

	sendData(addr, request)
}

func sendBlock(address string, b *blockchainstruct.Block) {
	data := block{nodeAddress, b.SerializeBlock()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(address, request)
}

func sendData(address string, data []byte) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		fmt.Printf("%s is not available\n", address)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != address {
				updatedNodes = append(updatedNodes, node)
			}
		}
		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	payload := gobEncode(inv{nodeAddress, kind, items})
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(address string, txn *blockchainstruct.Transaction) {
	payload := gobEncode(tx{nodeAddress, txn.SerializeTransaction()})
	request := append(commandToBytes("tx"), payload...)

	sendData(address, request)
}

func sendVersion(address string, bc *blockchainstruct.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeAddress, nodeVersion, bestHeight})
	request := append(commandToBytes("version"), payload...)

	sendData(address, request)
}

func handleAddress(request []byte) {
	var buff bytes.Buffer
	var payload address

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now.\n", len(knownNodes))
	requestBlocks()
}

func handleBlock(request []byte, bc *blockchainstruct.Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block :=blockchainstruct.DeserializeBlock(blockData)

	fmt.Println("Received a new block.")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.CurrHash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchainstruct.UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

func handleInv(request []byte, bc *blockchainstruct.Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if memPool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *blockchainstruct.Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *blockchainstruct.Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			log.Panic(err)
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memPool[txID]

		sendTx(payload.AddrFrom, &tx)
	}
}

func handleTx(request []byte, bc *blockchainstruct.Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := blockchainstruct.DeserializeTransaction(txData)
	memPool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(memPool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*blockchainstruct.Transaction

			for id := range memPool {
				tx := memPool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTX := blockchainstruct.NewCoinbaseTx(miningAddress, "")
			txs = append(txs, cbTX)

			newBlock := bc.MineBlock(txs)
			UTXOSet := blockchainstruct.UTXOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(memPool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.CurrHash})
				}
			}

			if len(memPool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *blockchainstruct.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLen:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignBestHeight := payload.BestHeight

	if myBestHeight < foreignBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func handleConnenction(conn net.Conn, bc *blockchainstruct.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLen])
	fmt.Printf("Receieved %s command\n", command)

	switch command {
	case "address":
		handleAddress(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

// StartServer starts a node
func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := blockchainstruct.NewBlockchain(nodeID)

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnenction(conn, bc)
	}
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(address string) bool {
	for _, node := range knownNodes {
		if node == address {
			return true
		}
	}
	return false
}
