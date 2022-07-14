package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

func NewTransaction(vin []TXInput, vout []TXOutput) *Transaction {
	tx := &Transaction{nil, vin, vout}
	tx.SetID()

	return tx
}

func (tx *Transaction) SetID() {
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(buf.Bytes())
	tx.ID = hash[:]
}

func NewCoinbaseTX(data, to string) *Transaction {
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}

	return NewTransaction([]TXInput{txin}, []TXOutput{txout})
}

func (bc *Blockchain) FindUnspentTransactions(address string) []*Transaction {
	bci := NewBlockchainIterator(bc)

	spentTXOs := make(map[string][]int)
	var unspentTXs []*Transaction

	for bci.HasNext() {
		for _, tx := range bci.Next().Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.ScriptPubKey == address {
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.ScriptSig == address {
						hash := hex.EncodeToString(in.Txid)
						spentTXOs[hash] = append(spentTXOs[hash], in.Vout)
					}
				}
			}
		}
	}

	return unspentTXs
}

func (tx *Transaction) IsCoinbase() bool {
	return bytes.Compare(tx.Vin[0].Txid, []byte{}) == 0 && tx.Vin[0].Vout == -1 && len(tx.Vin) == 1
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTXs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTXs {
		for _, out := range tx.Vout {
			if out.ScriptPubKey == address {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) GetBalance(address string) uint64 {
	var balance uint64

	for _, out := range bc.FindUTXO(address) {
		balance += out.Value
	}

	return balance
}

func (bc *Blockchain) Send(value uint64, from, to string) *Transaction {
	var txin []TXInput
	var txout []TXOutput

	UTXs := bc.FindUnspentTransactions(from)
	var acc uint64

Work:
	for _, tx := range UTXs {
		for outIdx, out := range tx.Vout {
			if out.ScriptPubKey == from && acc < value {
				acc += out.Value
				txin = append(txin, TXInput{tx.ID, outIdx, from})
			}
			if acc >= value {
				break Work
			}
		}
	}

	if value > acc {
		log.Panic("ERROR: NOT ENOUGH FUNDS")
	}

	txout = append(txout, TXOutput{value, to})
	if acc > value {
		txout = append(txout, TXOutput{acc - value, from})
	}

	return NewTransaction(txin, txout)
}
