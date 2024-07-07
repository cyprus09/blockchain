package blockchainstruct

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	dbFile              = "blockchain.db"
	blocksBucket        = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain keeps a sequence of Blocks in the blockchain
type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// MineBlock saves the provided data as a block in the blockchain
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.CurrHash, newBlock.SerializeBlock())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.CurrHash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.CurrHash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return newBlock
}

// FindUTXO finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.VOut {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.VIn {
					inTxID := hex.EncodeToString(in.TxId)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VOut)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.DB}

	return bci
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("Transaction not found.")
}

// NewBlockChain creates a new Blockchain with the genesis block
func NewBlockchain() *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	cbtx := NewCoinbaseTx(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.CurrHash, genesis.SerializeBlock())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.CurrHash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.CurrHash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, VIn := range tx.VIn {
		prevTX, err := bc.FindTransaction(VIn.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, VIn := range tx.VIn {
		prevTX, err := bc.FindTransaction(VIn.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

// Deprecated: defined utxo_set separately which takes care of FindUnspentTransactions returns a list of transactions containing unspent outputs
// func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
// 	var unspentTXs []Transaction
// 	spentTXOs := make(map[string][]int)
// 	bci := bc.Iterator()

// 	for {
// 		block := bci.Next()

// 		for _, tx := range block.Transactions {
// 			txID := hex.EncodeToString(tx.ID)

// 		Outputs:
// 			for outIdx, out := range tx.VOut {
// 				// Check if the output was spent or not
// 				if spentTXOs[txID] != nil {
// 					for _, spentOut := range spentTXOs[txID] {
// 						if spentOut == outIdx {
// 							continue Outputs
// 						}
// 					}
// 				}

// 				if out.IsLockedWithKey(pubKeyHash) {
// 					unspentTXs = append(unspentTXs, *tx)
// 				}
// 			}

// 			if !tx.IsCoinbase() {
// 				for _, in := range tx.VIn {
// 					if in.UsesKey(pubKeyHash) {
// 						inTxID := hex.EncodeToString(in.TxId)
// 						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VOut)
// 					}
// 				}
// 			}
// 		}

// 		if len(block.PrevBlockHash) == 0 {
// 			break
// 		}
// 	}

// 	return unspentTXs
// }
