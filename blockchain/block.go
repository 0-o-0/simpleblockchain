package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

// Block represents a block in blockchain
type Block struct {
	Height        int64
	Hash          string
	PrevBlockHash string
	Txs           []*Transaction
	Coinbase      *Transaction
	Nonce         uint32
}

// NewGenesisBlock creates a new genesis block
func NewGenesisBlock(addr string) *Block {
	coinbase := NewCoinbaseTx("") // TODO: add addr
	genesis := Block{
		Height:        0,
		PrevBlockHash: "",
		Txs:           []*Transaction{},
		Coinbase:      coinbase,
	}
	genesis.Hash = genesis.HashStr()
	return &genesis
}

// HashBytes returns the block hash as an array of bytes
func (blk *Block) HashBytes() []byte {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	encoder.Encode(blk.Height)
	encoder.Encode(blk.PrevBlockHash)
	encoder.Encode(blk.Txs)
	encoder.Encode(blk.Coinbase)
	encoder.Encode(blk.Nonce)
	bytes := sha256.Sum256(b.Bytes())
	return bytes[:]
}

// HashStr returns readable hash string in hex
func (blk *Block) HashStr() string {
	return hex.EncodeToString(blk.HashBytes())
}

// Print prints the block for debugging
func (blk *Block) Print() {
	fmt.Println("==========Block==========")
	fmt.Println("Height: ", blk.Height)
	fmt.Println("Hash: ", blk.Hash)
	fmt.Println("Prev: ", blk.PrevBlockHash)
	fmt.Println("Coinbase: ")
	blk.Coinbase.Print()
	fmt.Println("Txs: ")
	for _, tx := range blk.Txs {
		tx.Print()
	}
}
