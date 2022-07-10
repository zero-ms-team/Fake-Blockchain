package main

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewBlock("Genesis Block", []byte{})}}
}

func (bc *Blockchain) AddBlock(data string) {
	b := NewBlock(data, bc.blocks[len(bc.blocks)-1].Hash)
	bc.blocks = append(bc.blocks, b)
}
