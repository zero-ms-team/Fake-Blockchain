package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
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
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)

	return NewTransaction([]TXInput{txin}, []TXOutput{*txout})
}

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []*Transaction {
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
				if bytes.Equal(out.PubKeyHash, pubKeyHash) {
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
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
	return bytes.Equal(tx.Vin[0].Txid, []byte{}) && tx.Vin[0].Vout == -1 && len(tx.Vin) == 1
}

func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTXs {
		for _, out := range tx.Vout {
			if bytes.Equal(out.PubKeyHash, pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) GetBalance(address string) uint64 {
	var balance uint64

	pubKeyHash, _, err := base58.CheckDecode(address)
	if err != nil {
		log.Panic(err)
	}
	for _, out := range bc.FindUTXO(pubKeyHash) {
		balance += out.Value
	}

	return balance
}

func (bc *Blockchain) Send(value uint64, from, to string) *Transaction {
	var txin []TXInput
	var txout []TXOutput
	keyStore := NewKeyStore()

	wallet := keyStore.Wallets[from]
	UTXs := bc.FindUnspentTransactions(HashPubKey(wallet.PublicKey))

	var acc uint64

Work:
	for _, tx := range UTXs {
		for outIdx, out := range tx.Vout {
			if bytes.Equal(out.PubKeyHash, HashPubKey(wallet.PublicKey)) && acc < value {
				acc += out.Value
				txin = append(txin, TXInput{tx.ID, outIdx, nil, wallet.PublicKey})
			}
			if acc >= value {
				break Work
			}
		}
	}

	if value > acc {
		log.Panic("ERROR: NOT ENOUGH FUNDS")
	}

	txout = append(txout, *NewTXOutput(value, to))
	if acc > value {
		txout = append(txout, *NewTXOutput(acc-value, from))
	}

	tx := NewTransaction(txin, txout)
	bc.SignTransaction(wallet.PrivateKey, tx)

	return tx
}

func NewTXOutput(value uint64, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock(address)

	return txo
}

func (out *TXOutput) Lock(address string) {
	pubKeyHash, _, err := base58.CheckDecode(address)
	if err != nil {
		log.Panic(err)
	}
	out.PubKeyHash = pubKeyHash
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockinghash := HashPubKey(in.PubKey)
	return bytes.Equal(pubKeyHash, lockinghash)
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, in := range txCopy.Vin {
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTXs[hex.EncodeToString(in.Txid)].Vout[in.Vout].PubKeyHash
		txCopy.SetID()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}

		tx.Vin[inID].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Vin {
		inputs = append(inputs, TXInput{in.Txid, in.Vout, nil, nil})
	}
	for _, out := range tx.Vout {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}

	return &Transaction{nil, inputs, outputs}
}

func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, in := range tx.Vin {
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTXs[hex.EncodeToString(in.Txid)].Vout[in.Vout].PubKeyHash
		txCopy.SetID()
		txCopy.Vin[inID].PubKey = nil

		var r, s big.Int

		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:sigLen/2])
		s.SetBytes(in.Signature[sigLen/2:])

		var x, y big.Int

		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:keyLen/2])
		y.SetBytes(in.PubKey[keyLen/2:])

		pubKey := ecdsa.PublicKey{curve, &x, &y}

		if isVerified := ecdsa.Verify(&pubKey, txCopy.ID, &r, &s); !isVerified {
			return false
		}
	}

	return true
}
