package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const ( // boltdb info
	BlocksBucket = "blocks"
	dbFile       = "chain.DB"
)

// normal blockchain

func NewBlockchain() *Blockchain { // Make new blockchain, if exist on db, load it
	DB, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var l []byte

	err = DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		if b == nil {
			b, err = tx.CreateBucket([]byte(BlocksBucket))
			if err != nil {
				log.Panic(err)
			}

			genesis := NewBlock("Genesis Block", []byte{})

			err = b.Put(genesis.Hash, genesis.Serialize())

			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}

			l = genesis.Hash
		} else {
			l = b.Get([]byte("l"))
		}
		if err != nil {
			log.Panic(err)
		}

		return nil
	})

	return &Blockchain{DB, l}
}

func (bc *Blockchain) AddBlock(data string) {
	block := NewBlock(data, bc.L)

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.L = block.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) List() {
	bci := NewBlockchainIterator(bc)

	for bci.HasNext() {
		block := bci.Next()

		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Println("PoW: ", pow.Validate(block))

		fmt.Println()
	}
}

// Blockchain block iterator

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
