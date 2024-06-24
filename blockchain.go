package main

// Blockchain keeps a sequence of Blocks in the blockchain
type Blockchain struct {
	blocks []*Block
}

// AddBlock saves the provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.CurrHash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockChain creates a new Blockchain with the genesis block
func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
