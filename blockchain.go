package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const ( // boltdb info
	BlocksBucket = "blocks"
	dbFile       = "chain.DB"
)

// normal blockchain

func CreateBlockchain(address string) *Blockchain {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var l []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BlocksBucket))
		if err != nil {
			log.Panic(err)
		}

		genesis := NewBlock([]*Transaction{NewCoinbaseTX("", address)}, []byte{})

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		l = genesis.Hash

		return nil
	})

	return &Blockchain{db, l}
}

func GetBlockchain() *Blockchain {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var l []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		l = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{db, l}
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	block := NewBlock(transactions, bc.L)

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
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Println("PoW: ", pow.Validate(block))

		fmt.Println()
	}
}
