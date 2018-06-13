package miner

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/0-o-0/simpleblockchain/blockchain"
)

const (
	// MaxNonce defines the maximum possible value of nonce
	MaxNonce = ^uint32(0) // 2^32-1
	// BlockchainSyncSecs defines the number of seconds each miner waits to sync blockchain
	BlockchainSyncSecs = 10
)

var (
	defaultNumWorkers = uint32(runtime.NumCPU())
)

// Miner represents a node working on mining blocks
type Miner struct {
	sync.Mutex
	numWorkers       uint32
	bc               *blockchain.Blockchain
	quit             chan bool
	publishBlockLock sync.Mutex
}

// New creates a new miner with default number of workers and the given blockchain
func New(bc *blockchain.Blockchain) *Miner {
	return &Miner{
		numWorkers: uint32(defaultNumWorkers),
		bc:         bc,
	}
}

// Start kicks off workers to solve the hash puzzle in a concurrency-safe manner
func (m *Miner) Start() {
	var workers []chan bool
	for i := uint32(0); i < m.numWorkers; i++ {
		quit := make(chan bool)
		workers = append(workers, quit)
		go m.mine(quit)
	}

	for {
		select {
		case <-m.quit:
			for _, quit := range workers {
				close(quit)
			}
			return
		}
	}
}

func (m *Miner) mine(quit chan bool) {
	ticker := time.NewTicker(time.Second * BlockchainSyncSecs)
	defer ticker.Stop()
	for {
		block := m.bc.NewBlock()
		if m.findNonce(block, ticker, quit) {
			block.Hash = block.HashStr()
			m.publishBlock(block)
		}
	}
}

func (m *Miner) findNonce(block *blockchain.Block, ticker *time.Ticker,
	quit chan bool) bool {

	for i := uint32(0); i < MaxNonce; i++ {
		select {
		case <-quit:
			return false
		case <-ticker.C:
			if m.bc.Height()+1 > block.Height {
				return false
			}
		default:
			// keep tryinng
		}

		block.Nonce = i
		hash := block.HashBytes()
		if blockchain.ValidateBlockHash(hash) {
			return true
		}
	}

	return false
}

func (m *Miner) publishBlock(blk *blockchain.Block) {
	m.publishBlockLock.Lock()
	defer m.publishBlockLock.Unlock()

	if m.bc.ProcessBlock(blk) {
		fmt.Println("A new block has been mined!")
		m.bc.Print()
	}
}
