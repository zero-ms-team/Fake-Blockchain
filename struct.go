package main

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
	blocks []*Block
}

type Block struct {
	PrevBlockHash []byte
	Hash          []byte
	Timestamp     int64
	Data          []byte
	Nonce         int64
}
