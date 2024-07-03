package blockchainstruct

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/cyprus09/blockchain/utils"
)

// TxOutput represents a transaction output
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock signs the output
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTxOutput creates a new TXOutput
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// TxOutputs collects TxOutput
type TxOutputs struct {
	Outputs []TxOutput
}

// SerializeOutputs serializes TxOutputs
func (outs *TxOutputs) SerializeOutputs() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TxOutputs
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
