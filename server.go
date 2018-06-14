package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	mrand "math/rand"

	"github.com/0-o-0/simpleblockchain/blockchain"
	"github.com/0-o-0/simpleblockchain/miner"
)

const (
	// KeyBitSize defines the number of bits of RSA key
	KeyBitSize = 1024
	// NumKeys defines the number of different keys used on the server
	NumKeys = 10
)

type server struct {
	blockchain *blockchain.Blockchain
	miner      *miner.Miner
	keys       []*rsa.PrivateKey
}

func newServer() *server {
	s := server{}
	s.keys = generateKeys()
	s.blockchain = blockchain.New(s.getRandomKey().PublicKey)
	s.blockchain.Print()
	s.miner = miner.New(s.blockchain, s.keys)
	return &s
}

// Run starts the server to mine blocks and handle network requests
func (s *server) Run() {
	// TODO: handle peers
	s.miner.Start()
}

func generateKeys() []*rsa.PrivateKey {
	var keys []*rsa.PrivateKey
	for i := 0; i < NumKeys; i++ {
		key, err := rsa.GenerateKey(rand.Reader, KeyBitSize)
		if err != nil {
			fmt.Println("RSA key generation failed")
		}
		keys = append(keys, key)
	}
	return keys
}

func (s *server) getRandomKey() *rsa.PrivateKey {
	return s.keys[mrand.Intn(NumKeys)]
}
