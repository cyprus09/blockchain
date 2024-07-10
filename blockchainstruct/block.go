package blockchainstruct

import (
	"bytes"
	"encoding/gob" //gob is the library used for encoding data (serialisation which can be done through protobufs as well for data streams in binary format
	"log"
	"time"

	"github.com/cyprus09/blockchain/merkletree"
)

// Block struct helps define the structure of a block
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	CurrHash      []byte
	Nonce         int
	Height        int
}

// Deprecated: SetHash calculates the hash of the current block
// Not used anymore since we use proof of work concept to generate hash for each block
// func (b *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.BlockData, timestamp}, []byte{})
// 	hash := sha256.Sum256(headers)

// 	b.CurrHash = hash[:]
// }

// HashTransaction return a hash of the transaction within a block
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.SerializeTransaction())
	}
	mTree := merkletree.NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// NewBlock creates and returns a new Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, currHash := pow.Run()

	block.CurrHash = currHash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns the genesis block (first block of the blockchain)
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// SerializeBlock serializes the value of the block to be able to store in BoltDb
func (b *Block) SerializeBlock() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes the block value got from the db
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
