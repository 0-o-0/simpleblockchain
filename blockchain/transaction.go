package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
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
	Signature       []byte
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
func NewCoinbaseTx(pk rsa.PublicKey, txFee float64) *Transaction {
	coinbase := Transaction{
		Inputs: []Input{},
		Outputs: []Output{
			Output{COINBASE + txFee, pk},
		},
		IsCoinbase: true,
		Ts:         time.Now().UnixNano(),
	}
	coinbase.Hash = coinbase.hashStr()
	return &coinbase
}

// NewTransaction creates a new transaction with the input, output and signature
// key provided
func NewTransaction(input Input, output Output, key *rsa.PrivateKey) *Transaction {
	hashed := GetHashToSign(input, []Output{output})
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return nil
	}
	input.Signature = sig
	tx := Transaction{
		Inputs:     []Input{input},
		Outputs:    []Output{output},
		IsCoinbase: false,
		Ts:         time.Now().UnixNano(),
	}
	tx.Hash = tx.hashStr()
	return &tx
}

// GetHashToSign returns the hash of the input and all the output
// to be signed by the private key
func GetHashToSign(input Input, outputs []Output) []byte {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	encoder.Encode(input.PrevTxHash)
	encoder.Encode(input.PrevTxOutputIdx)
	for _, output := range outputs {
		encoder.Encode(output.Address)
		encoder.Encode(output.Value)
	}
	hash := sha256.Sum256(b.Bytes())
	return hash[:]
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

// ValidateInput validates the input signature with the public key
func (tx *Transaction) ValidateInput(input Input, pk rsa.PublicKey) bool {
	hashed := GetHashToSign(input, tx.Outputs)
	err := rsa.VerifyPKCS1v15(&pk, crypto.SHA256, hashed, input.Signature)
	return err == nil
}

// Print prints the transaction for debugging
func (tx *Transaction) Print() {
	if tx.IsCoinbase {
		fmt.Printf("  =Tx %s coinbase to %s for %.1f\n", tx.Hash, tx.Outputs[0].getAddressHashStr(), tx.Outputs[0].Value)
	} else {
		fmt.Printf("  =Tx %s from %s to %s for %.1f\n", tx.Hash, tx.Inputs[0].PrevTxHash, tx.Outputs[0].getAddressHashStr(), tx.Outputs[0].Value)
	}
}
