package main

import (
	"bytes"
	"log"

	"github.com/boltdb/bolt"
)

func NewBlockchainIterator(bc *Blockchain) *BlockchainIterator {
	return &BlockchainIterator{bc.DB, bc.L}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		encodedBlock := b.Get(i.Hash)
		block = DeserializeBlock(encodedBlock)

		i.Hash = block.PrevBlockHash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return block
}

func (i *BlockchainIterator) HasNext() bool {
	return bytes.Compare(i.Hash, []byte{}) != 0
}