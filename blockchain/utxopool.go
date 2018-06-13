package blockchain

import (
	"fmt"
)

// UTXOPool represents a list of unspent transaction output
type UTXOPool struct {
  outputs map[UTXO]Output
}

// UTXO represents unspent transaction output
type UTXO struct {
  txHash  string
  index   int32
}

// ProcessBlockTxs updates the UTXOPool by processing transactions in the block
// and returns false if any of the transactio is invalid
func (up *UTXOPool) ProcessBlockTxs(blk *Block) bool {
  for _, tx := range blk.Txs {
    if !up.processTx(tx) {
      return false
    }
  }

  if !up.processTx(blk.Coinbase) {
    return false
  }

  return true
}

func (up *UTXOPool) processTx(tx *Transaction) bool {
  for i, output := range tx.Outputs {
    utxo := UTXO{tx.Hash, int32(i)}
    _, exist := up.outputs[utxo]
    if exist {
      return false
    }
    up.outputs[utxo] = output
  }

  for _, input := range tx.Inputs {
    utxo := UTXO{
      input.PrevTxHash,
      input.PrevTxOutputIdx,
    }
    _, exist := up.outputs[utxo]
    if !exist {
      return false
    }
    delete(up.outputs, utxo)
  }

  return true
}

// Print prints the blockchain for debugging
func (up *UTXOPool) Print() {
	fmt.Println("==========UTXOP==========")
  for utxo, output := range up.outputs {
    fmt.Println(utxo.txHash, utxo.index, output.Value)
  }
}
