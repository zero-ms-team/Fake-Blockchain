package main

import (
	"crypto/ecdsa"

	"github.com/boltdb/bolt"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	Value      uint64
	PubKeyHash []byte
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type Blockchain struct {
	DB *bolt.DB
	L  []byte
}

type BlockchainIterator struct {
	DB   *bolt.DB
	Hash []byte
}

type Block struct {
	PrevBlockHash []byte
	Hash          []byte
	Timestamp     int64
	Transactions  []*Transaction
	Nonce         int64
}

type CLI struct{}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

type KeyStore struct {
	Wallets map[string]*Wallet
}
