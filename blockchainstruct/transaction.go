package blockchainstruct

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob" //gob is the library used for encoding data (serialisation which can be done through protobufs as well for data streams in binary format
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"github.com/cyprus09/blockchain/wallets"
)

const subsidy = 10

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

// Serialize returns a serialized Transaction
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// HashValue returns the hash of the transaction
func (tx *Transaction) HashValue() []byte {
	var hashValue [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hashValue = sha256.Sum256(txCopy.Serialize())

	return hashValue[:]
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, VIn := range tx.VIn {
		if prevTXs[hex.EncodeToString(VIn.TxId)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, VIn := range txCopy.VIn {
		prevTX := prevTXs[hex.EncodeToString(VIn.TxId)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTX.VOut[VIn.VOut].PubKeyHash
		txCopy.ID = txCopy.HashValue()
		txCopy.VIn[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)
		tx.VIn[inID].Signature = signature
	}
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x", tx.ID))

	for i, input := range tx.VIn {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TxId))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.VOut))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.VOut {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")

}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, VIn := range tx.VIn {
		inputs = append(inputs, TxInput{VIn.TxId, VIn.VOut, nil, nil})
	}

	for _, VOut := range tx.VOut {
		outputs = append(outputs, TxOutput{VOut.Value, VOut.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, VIn := range tx.VIn {
		if prevTXs[hex.EncodeToString(VIn.TxId)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, VIn := range tx.VIn {
		prevTX := prevTXs[hex.EncodeToString(VIn.TxId)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTX.VOut[VIn.VOut].PubKeyHash
		txCopy.ID = txCopy.HashValue()
		txCopy.VIn[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(VIn.Signature)
		r.SetBytes(VIn.Signature[:(sigLen / 2)])
		s.SetBytes(VIn.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(VIn.PubKey)
		x.SetBytes(VIn.PubKey[:(keyLen / 2)])
		y.SetBytes(VIn.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}

// NewCoinbaseTx creates a new coinbase transaction
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTxOutput(subsidy, to)

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	tx.ID = tx.HashValue()

	return &tx
}

// NewUTXOTTransaction creates a new transaction
func NewUTXOTTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

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
			input := TxInput{txId, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTxOutput(amount, to))
	if acc > amount {
		// generate change
		outputs = append(outputs, *NewTxOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.HashValue()
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}