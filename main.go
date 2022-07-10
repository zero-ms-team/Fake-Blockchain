package main

import (
	"fmt"
)

func main() {
	fmt.Println("Blockchain Start!")
	fmt.Println()

	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, b := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)

		pow := NewProofOfWork(b)
		fmt.Println("pow: ", pow.Validate(b))

		fmt.Println()
	}
}
