package blockchain

import ()

// UTXOPool represents a list of unspent transaction output
type UTXOPool struct {
  outputs map[UTXO]Output
}

// UTXO represents unspent transaction output
type UTXO struct {
  txHash  string
  index   int32
}

// GetUpdatedUTXOPool creates a new UTXOPool after processing transactions
func (up *UTXOPool) GetUpdatedUTXOPool(txs []*Transaction, coinbase *Transaction) *UTXOPool {
  newPool := *up
  for _, tx := range txs {
    newPool.processTx(tx);
  }
  return &newPool
}

func (up *UTXOPool) processTx(tx *Transaction) {
  for _, input := range tx.Inputs {
    utxo := UTXO{
      input.PrevTxHash,
      input.PrevTxOutputIdx,
    }
    delete(up.outputs, utxo)
  }

  for i, output := range tx.Outputs {
    utxo := UTXO{tx.Hash, int32(i)}
    up.outputs[utxo] = output
  }
}
