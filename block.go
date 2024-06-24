package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block struct helps define the structure of a block
type Block struct {
	Timestamp     int64
	BlockData     []byte
	PrevBlockHash []byte
	CurrHash      []byte
}

// SetHash calculates the hash of the current block
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.BlockData, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.CurrHash = hash[:]
}

// NewBlock creates and returns a new Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// NewGenesisBlock creates and returns the genesis block (first block of the blockchain)
func NewGenesisBlock() *Block {
	return NewBlock("New Genesis Block", []byte{})
}
