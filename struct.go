package main

import "github.com/boltdb/bolt"

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
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
	Data          []byte
	Nonce         int64
}

type CLI struct{}
