package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/0-o-0/simpleblockchain/blockchain"
	"github.com/0-o-0/simpleblockchain/miner"
	"github.com/0-o-0/simpleblockchain/peer"
)

const (
	// KeyBitSize defines the number of bits of RSA key
	KeyBitSize = 1024
	// NumKeys defines the number of different keys used on the server
	NumKeys = 10
	// MaxTxFee defines the maximum transaction fee
	MaxTxFee = 1
	// TxGenSecs defines the number of seconds the server waits to generate
	// a new transaction
	TxGenSecs = 2
)

type server struct {
	blockchain *blockchain.Blockchain
	miner      *miner.Miner
	keys       []*rsa.PrivateKey
	peers      *peer.Peer
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
	go s.miner.Start()
	s.genTxs()
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

func (s *server) genTxs() {
	for {
		utxoPool := s.blockchain.GetUTXOPoolCopy()
		utxo, spent := utxoPool.GetRandomUTXO()

		input := blockchain.Input{
			PrevTxHash:      utxo.TxHash,
			PrevTxOutputIdx: utxo.Index,
			Signature:       nil,
		}
		var sigKey *rsa.PrivateKey
		for _, key := range s.keys {
			if key.PublicKey == spent.Address {
				sigKey = key
				break
			}
		}

		addr := s.getRandomKey().PublicKey
		txFee := mrand.Float64() * MaxTxFee
		output := blockchain.Output{
			Value:   spent.Value - txFee,
			Address: addr,
		}

		tx := blockchain.NewTransaction(input, output, sigKey)

		if s.blockchain.ProcessTx(tx) {
			fmt.Println("A new tx has been generated")
			tx.Print()
		}

		time.Sleep(TxGenSecs * time.Second)
	}
}
