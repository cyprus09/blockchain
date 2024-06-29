package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob" //gob is the library used for encoding data (serialisation which can be done through protobufs as well for data streams in binary format
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

// TxInput represents a transaction input
type TxInput struct {
	TxId      []byte
	VOut      int
	ScriptSig string
}

// TxOutput represents a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// Transaction struct represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	VIn  []TxInput
	VOut []TxOutput
}

// IsCoinbase checks whether a transaction is a coinbase transaction or not
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxId) == 0 && tx.VIn[0].VOut == -1
}

// SetId sets the transaction ID
func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hashValue [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hashValue = sha256.Sum256(encoded.Bytes())
	tx.ID = hashValue[:]
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTx creates a new coinbase transaction
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{subsidy, to}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetId()

	return &tx
}

// NewUTXOTTransaction creates a new transaction
func NewUTXOTTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds in your account")
	}

	// Build a list of inputs
	for txId, outs := range validOutputs {
		txId, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txId, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		// generate change
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()

	return &tx
}
