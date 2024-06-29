package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 24

// ProofOfWork represents proof-of-work for a blockchain
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns the proof of work for the block
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData is a private function that helps to join the header values of the block before hashing
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToBytes(pow.block.Timestamp),
			IntToBytes(int64(targetBits)),
			IntToBytes(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run performs the proof of work from the data generated from prepareData
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hashValue [32]byte
	nonce := 0

	fmt.Printf("Mining new block now...")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hashValue = sha256.Sum256(data)
		fmt.Printf("\r%x", hashValue)
		hashInt.SetBytes(hashValue[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")
	return nonce, hashValue[:]
}

// Validate validates the block's PoW which takes way lesser time than the actual process of generating the hash
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
