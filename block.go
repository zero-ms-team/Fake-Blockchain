package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
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

func (b *Block) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var buf bytes.Buffer
	var block Block

	buf.Write(d)
	dec := gob.NewDecoder(&buf)

	err := dec.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
