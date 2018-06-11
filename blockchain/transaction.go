package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

const (
	// COINBASE defines the reward for creating a block
	COINBASE = 50
)

// Input represents the input part of the transaction
type Input struct {
	PrevTxHash      string
	PrevTxOutputIdx int32
	Signature       string
}

// Output represents the output part of the transaction
type Output struct {
	Value   float64
	Address string // PK
}

// Transaction represents the transactions recorded in blocks
type Transaction struct {
	Inputs     []Input
	Outputs    []Output
	Hash       string
	IsCoinbase bool
}

// NewCoinbaseTx creates a new coinbase transaction which rewards the block creator
func NewCoinbaseTx(addr string) *Transaction {
	coinbase := Transaction{
		Inputs: []Input{},
		Outputs: []Output{
			Output{COINBASE, addr},
		},
		IsCoinbase: true,
	}
	coinbase.Hash = coinbase.hashStr()
	return &coinbase
}

func (tx *Transaction) hashStr() string {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	encoder.Encode(tx.Inputs)
	encoder.Encode(tx.Outputs)
	encoder.Encode(tx.IsCoinbase)
	hash := sha256.Sum256(b.Bytes())
	return hex.EncodeToString(hash[:])
}

// Print prints the transaction for debugging
func (tx *Transaction) Print() {
	if tx.IsCoinbase {
		fmt.Printf("  =Transaction= coinbase to %s for %f\n", tx.Outputs[0].Address, tx.Outputs[0].Value)
	} else {
		fmt.Printf("  =Transaction= from %s to %s for %f\n", tx.Inputs[0].PrevTxHash, tx.Outputs[0].Address, tx.Outputs[0].Value)
	}
}
