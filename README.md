Blockchain
====

This is a Go implementation of a simple blockchain.

## Features
* Blockchain
  * Data structures of blockchain, block, transaction, UTXO, and UTXO pool
  * SHA256 for all the hashing
  * Transactions are signed with RSA digital signature
* Consensus
  * PoW with the same puzzle as Bitcoin
  * Validations for blocks, transactions, double-spendings, digital signatures, etc
* Mining
  * Multiple goroutines working on the puzzle in a concurrency-safe manner
* Network
  * The client generates transactions automatically to simulate the network
  * Single machine for now, will implement a gossip protocol
* Logs
  * Event logging for block mining/processing, transaction generation/processing

## Run
``` bash
go build
./simpleblockchain
```
