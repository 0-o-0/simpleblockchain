Blockchain
====

This is a Go implementation of a simple blockchain.

## Features
* Blockchain
  * Data structures of blockchain, block, transaction, UTXO, and UTXO pool
  * Hashing in SHA256
* Consensus
  * PoW with the same puzzle as Bitcoin
* Mining
  * Multiple goroutines working on the puzzle in a concurrency-safe manner
* Network
  * Single machine for now, will implement a gossip protocol

## Run
``` bash
go build
./simpleblockchain
```
