package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

func (b *Block) PrintBlock() {
	fmt.Println("-----------")
	fmt.Printf("Previous hash: %x\n", b.PrevHash)
	fmt.Printf("Data: %s\n", b.Data)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("Nonce: %d\n", b.Nonce)
	fmt.Println("-----------")
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	block.PrintBlock()
	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis Block", []byte{})
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
