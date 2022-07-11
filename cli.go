package main

import (
	"flag"
	"os"
)

func (c *CLI) CreateBlockchain() {
	bc := NewBlockchain()
	bc.DB.Close()
}

func (c *CLI) AddBlock(data string) {
	bc := NewBlockchain()
	defer bc.DB.Close()

	bc.AddBlock(data)
}

func (c *CLI) list() {
	bc := NewBlockchain()
	defer bc.DB.Close()

	bc.List()
}

func (c *CLI) Run() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	addBlockData := addCmd.String("data", "", "")

	switch os.Args[1] {
	case "new":
		newCmd.Parse(os.Args[2:])
	case "add":
		addCmd.Parse(os.Args[2:])
	case "list":
		listCmd.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}

	if newCmd.Parsed() {
		c.CreateBlockchain()
	}

	if addCmd.Parsed() {
		if *addBlockData == "" {
			addCmd.Usage()
			os.Exit(1)
		}
		c.AddBlock(*addBlockData)
	}

	if listCmd.Parsed() {
		c.list()
	}
}
