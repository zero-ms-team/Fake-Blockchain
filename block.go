package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data)}
	block.SetHash()

	type Block struct {
		PrevBlockHash []byte
		Hash          []byte
		Timestamp     int64
		Data          []byte
	}
	return block
}

func (b *Block) SetHash() {
	header := bytes.Join([][]byte{
		b.PrevBlockHash,
		b.Data,
		IntToHex(b.Timestamp),
	}, []byte{})

	hash := sha256.Sum256(header)
	b.Hash = hash[:]
}
