package blockchain

// Blockchain represents the current blockchain including all the blocks,
// received transactions to be processed, and UTXO pool for each block
type Blockchain struct {
	lastBlock *Block
	blocks    map[string]*Block
	utxoPools map[string]*UTXOPool
	txPool    []*Transaction
}

// New creates a new blockchain with genesis block
func New() *Blockchain {
	genesis := NewGenesisBlock("") // TODO: add addr
	bc := Blockchain{
		lastBlock: genesis,
		blocks:    map[string]*Block{genesis.Hash: genesis},
		utxoPools: map[string]*UTXOPool{
			genesis.Hash: &UTXOPool{
				outputs: map[UTXO]Output{UTXO{genesis.Hash, 0}: genesis.Coinbase.Outputs[0]},
			},
		},
	}
	return &bc
}

func (bc *Blockchain) getLastBlock() *Block {
	return bc.lastBlock
}

// ProcessBlock validates the block and adds it to the blockchain if valid
func (bc *Blockchain) ProcessBlock(blk *Block) bool {
	//TODO: add validation
	bc.lastBlock = blk
	bc.blocks[blk.Hash] = blk
	bc.utxoPools[blk.Hash] = bc.utxoPools[blk.PrevBlockHash].GetUpdatedUTXOPool(blk.Txs, blk.Coinbase)
	return true
}

func (tx *Transaction) addTx() bool {
	return false // TODO
}

// NewBlock creates a new block on top of the tip block
func (bc *Blockchain) NewBlock() *Block {
	return &Block{
		Height:        bc.lastBlock.Height + 1,
		PrevBlockHash: bc.lastBlock.Hash,
		Txs:           bc.txPool,
		Coinbase:      NewCoinbaseTx(""), // TODO: add addr
		Nonce:         0,
	}
}

// Height returns the max height of the blocks in the blockchain
func (bc *Blockchain) Height() int64 {
	return bc.lastBlock.Height
}

// Print prints the blockchain for debugging
func (bc *Blockchain) Print() {
	for _, blk := range bc.blocks {
		blk.Print()
	}
}
