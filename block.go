package main

import (
	"time"
)

// Block struct helps define the structure of a block
type Block struct {
	Timestamp     int64
	BlockData     []byte
	PrevBlockHash []byte
	CurrHash      []byte
	Nonce         int
}

// Deprecated: SetHash calculates the hash of the current block
// Not used anymore since we use proof of work concept to generate hash for each block
// func (b *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.BlockData, timestamp}, []byte{})
// 	hash := sha256.Sum256(headers)

// 	b.CurrHash = hash[:]
// }

// NewBlock creates and returns a new Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, currHash := pow.Run()

	block.CurrHash = currHash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns the genesis block (first block of the blockchain)
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
