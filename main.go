package main

import (
	"fmt"
	"github.com/ladderseeker/simple-blockchain/blockchain"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	validateChain(chain)
}

func validateChain(chain *blockchain.Blockchain) {
	for _, block := range chain.Blocks {
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow Validated %x: %s\n", block.Hash, strconv.FormatBool(pow.Validate()))
	}
}
