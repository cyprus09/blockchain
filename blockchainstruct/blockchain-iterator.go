package blockchainstruct

import (
	"github.com/boltdb/bolt"
	"log"
)

// BlockchainIterator is used to iterate over the blockchain blocks
type BlockchainIterator struct {
	currHash []byte
	db       *bolt.DB
}

// NextBlock returns the next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currHash = block.PrevBlockHash

	return block
}
