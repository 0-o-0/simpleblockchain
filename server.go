package main

import (
	"github.com/0-o-0/simpleblockchain/blockchain"
	"github.com/0-o-0/simpleblockchain/miner"
)

type server struct {
	blockchain *blockchain.Blockchain
	miner      *miner.Miner
}

func newServer() *server {
	s := server{}
	s.blockchain = blockchain.New()
	s.blockchain.Print()
	s.miner = miner.New(s.blockchain)
	return &s
}

// Run starts the server to mine blocks and handle network requests
func (s *server) Run() {
	// TODO: handle peers
	s.miner.Start()
}
