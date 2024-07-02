package transactions

import "bytes"

// TxInput represents a transaction input
type TxInput struct {
	TxId      []byte
	VOut      int
	Signature []byte
	PubKey		[]byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}