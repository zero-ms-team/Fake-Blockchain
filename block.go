package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data), 0}
	pow := NewProofOfWork(block)
	block.Nonce, block.Hash = pow.Run()

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
