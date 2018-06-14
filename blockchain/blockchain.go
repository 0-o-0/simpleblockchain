package blockchain

import (
	"crypto/rsa"
	"fmt"
	"math/big"
)

const (
	// CutOffAge defines the acceptable maximum height gap to the current tip block
	CutOffAge = 2
	// DefaultDifficulty defines the number of leading zeros of the targeted difficulty
	DefaultDifficulty = 19
)

// Blockchain represents the current blockchain including all the blocks,
// received transactions to be processed, and UTXO pool for each block
type Blockchain struct {
	lastBlock *Block
	blocks    map[string]*Block
	utxoPools map[string]*UTXOPool
	txPool    []*Transaction
	utxoPool  *UTXOPool
	txFee     float64
}

// New creates a new blockchain with genesis block
func New(pk rsa.PublicKey) *Blockchain {
	genesis := NewGenesisBlock(pk)
	utxoPool := &UTXOPool{
		outputs: map[UTXO]Output{
			UTXO{genesis.Coinbase.Hash, 0}: genesis.Coinbase.Outputs[0],
		},
	}
	return &Blockchain{
		lastBlock: genesis,
		blocks:    map[string]*Block{genesis.Hash: genesis},
		utxoPools: map[string]*UTXOPool{
			genesis.Hash: utxoPool,
		},
		txPool:   []*Transaction{},
		utxoPool: utxoPool.copy(),
	}
}

// GetUTXOPoolCopy returns a copy of latest UTXO pool of the blockchain
func (bc *Blockchain) GetUTXOPoolCopy() *UTXOPool {
	return bc.utxoPool.copy()
}

// ProcessBlock validates the block and adds it to the blockchain if valid
func (bc *Blockchain) ProcessBlock(blk *Block) bool {
	if !bc.validateBlock(blk) {
		return false
	}

	newUTXOPool := bc.utxoPools[blk.PrevBlockHash].copy()
	if !newUTXOPool.ProcessBlockTxs(blk) {
		return false
	}

	bc.blocks[blk.Hash] = blk
	bc.utxoPools[blk.Hash] = newUTXOPool
	if blk.Height > bc.Height() {
		bc.lastBlock = blk
		bc.txPool = []*Transaction{}
		bc.utxoPool = newUTXOPool.copy()
		bc.txFee = 0.0
	}
	return true
}

func (bc *Blockchain) validateBlock(blk *Block) bool {
	// if the block exist in the blockchain already, then reject
	if _, exist := bc.blocks[blk.Hash]; exist {
		return false
	}

	prevBlk, ok := bc.blocks[blk.PrevBlockHash]
	// if the previous block doesn't exist in the current blockchain, then reject
	if !ok {
		return false
	}

	// the block's height has to be the previous one's + 1
	if blk.Height != prevBlk.Height+1 {
		return false
	}

	// if the height is below the cut off age, then reject
	if bc.Height()-blk.Height > CutOffAge {
		return false
	}

	// recalculate the hash of the block to make sure it matches the block hash
	if blk.HashStr() != blk.Hash {
		return false
	}

	// check the block hash meets the target difficulty
	hashBytes := blk.HashBytes()
	if !ValidateBlockHash(hashBytes) {
		return false
	}

	return true
}

// ProcessTx handles the transaction on top of the tip block of the blockchain
func (bc *Blockchain) ProcessTx(tx *Transaction) bool {
	pool := bc.utxoPool.copy()
	valid, txFee := pool.ProcessTx(tx)
	if !valid {
		return false
	}
	bc.utxoPool = pool
	bc.txPool = append(bc.txPool, tx)
	bc.txFee += txFee
	return true
}

// NewBlock creates a new block on top of the tip block
func (bc *Blockchain) NewBlock(pk rsa.PublicKey) *Block {
	return &Block{
		Height:        bc.lastBlock.Height + 1,
		PrevBlockHash: bc.lastBlock.Hash,
		Txs:           bc.getTxs(),
		Coinbase:      NewCoinbaseTx(pk, bc.txFee),
		Nonce:         0,
	}
}

// Height returns the max height of the blocks in the blockchain
func (bc *Blockchain) Height() int64 {
	return bc.lastBlock.Height
}

// ValidateBlockHash validates the hash meets the target difficulty
func ValidateBlockHash(bytes []byte) bool {
	target := new(big.Int)
	target.SetBit(target, 256-DefaultDifficulty, 1)
	hash := new(big.Int).SetBytes(bytes)
	return hash.Cmp(target) <= 0
}

func (bc *Blockchain) getTxs() []*Transaction {
	txs := make([]*Transaction, len(bc.txPool))
	copy(txs, bc.txPool)
	return txs
}

// GetTxCount returns the number of transactions in the pool of the blockchain
func (bc *Blockchain) GetTxCount() int {
	return len(bc.txPool)
}

// Print prints the blockchain for debugging
func (bc *Blockchain) Print() {
	fmt.Println("==========Chain==========")
	fmt.Println("Total blocks: ", len(bc.blocks))
	bc.lastBlock.Print()
	bc.utxoPools[bc.lastBlock.Hash].Print()
}
