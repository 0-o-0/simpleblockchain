package blockchain

import (
	"fmt"
	"math/rand"
)

// UTXOPool represents a list of unspent transaction output
type UTXOPool struct {
	outputs map[UTXO]Output
}

// UTXO represents unspent transaction output
type UTXO struct {
	TxHash string
	Index  int32
}

// ProcessBlockTxs updates the UTXOPool by processing transactions in the block
// and returns false if any of the transactio is invalid
func (up *UTXOPool) ProcessBlockTxs(blk *Block) bool {
	for _, tx := range blk.Txs {
		if valid, _ := up.ProcessTx(tx); !valid {
			return false
		}
	}

	if valid, _ := up.ProcessTx(blk.Coinbase); !valid {
		return false
	}

	return true
}

// ProcessTx updates the UTXOPool by processing the transaction
// It returns false if the transaction is invalid
func (up *UTXOPool) ProcessTx(tx *Transaction) (bool, float64) {
	txFee := 0.0
	for i, output := range tx.Outputs {
		utxo := UTXO{tx.Hash, int32(i)}
		_, exist := up.outputs[utxo]
		if exist {
			return false, txFee
		}
		up.outputs[utxo] = output
		txFee -= output.Value
	}

	for _, input := range tx.Inputs {
		utxo := UTXO{
			input.PrevTxHash,
			input.PrevTxOutputIdx,
		}
		output, exist := up.outputs[utxo]
		if !exist {
			return false, txFee
		}
		if !tx.ValidateInput(input, output.Address) {
			return false, txFee
		}
		delete(up.outputs, utxo)
		txFee += output.Value
	}

	return true, txFee
}

// GetRandomUTXO returns a random UTXO in the UTXO pool
func (up *UTXOPool) GetRandomUTXO() (*UTXO, *Output) {
	idx := rand.Intn(len(up.outputs))
	for utxo, output := range up.outputs {
		if idx == 0 {
			return &utxo, &output
		}
		idx--
	}
	return nil, nil
}

func (up *UTXOPool) copy() *UTXOPool {
	newPool := UTXOPool{outputs: map[UTXO]Output{}}
	for utxo, output := range up.outputs {
		newPool.outputs[utxo] = output
	}
	return &newPool
}

// Print prints the blockchain for debugging
func (up *UTXOPool) Print() {
	fmt.Println("==========UTXOP==========")
	for utxo, output := range up.outputs {
		fmt.Printf("%s idx: %d val: %1f\n", utxo.TxHash, utxo.Index, output.Value)
	}
}
