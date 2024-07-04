package blockchainstruct

import (
	"bytes"
	"github.com/cyprus09/blockchain/wallets"
)

// TxInput represents a transaction input
type TxInput struct {
	TxId      []byte
	VOut      int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallets.NewWallet().HashPubKey(in.PubKey)

	return bytes.Equal(lockingHash, pubKeyHash)
}
