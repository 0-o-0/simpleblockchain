package blockchain

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"
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
	Address rsa.PublicKey
}

// Transaction represents the transactions recorded in blocks
type Transaction struct {
	Inputs     []Input
	Outputs    []Output
	Hash       string
	IsCoinbase bool
	Ts         int64
}

// NewCoinbaseTx creates a new coinbase transaction which rewards the block creator
func NewCoinbaseTx(pk rsa.PublicKey) *Transaction {
	coinbase := Transaction{
		Inputs: []Input{},
		Outputs: []Output{
			Output{COINBASE, pk},
		},
		IsCoinbase: true,
		Ts:         time.Now().UnixNano(),
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
	encoder.Encode(tx.Ts)
	hash := sha256.Sum256(b.Bytes())
	return hex.EncodeToString(hash[:])
}

func (o Output) getAddressHashStr() string {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	encoder.Encode(o.Address)
	hash := sha256.Sum256(b.Bytes())
	return hex.EncodeToString(hash[:])
}

// Print prints the transaction for debugging
func (tx *Transaction) Print() {
	if tx.IsCoinbase {
		fmt.Printf("  =Tx %s coinbase to %s for %.1f\n", tx.Hash, tx.Outputs[0].getAddressHashStr(), tx.Outputs[0].Value)
	} else {
		fmt.Printf("  =Tx %s from %s to %s for %.1f\n", tx.Hash, tx.Inputs[0].PrevTxHash, tx.Outputs[0].getAddressHashStr(), tx.Outputs[0].Value)
	}
}
