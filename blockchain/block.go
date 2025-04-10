package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
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

	// Print detailed transaction information.
	fmt.Println("Transactions:")
	for _, tx := range b.Transactions {
		fmt.Printf("  Transaction ID: %x\n", tx.ID)

		// Print inputs for the transaction.
		fmt.Println("  Inputs:")
		for _, in := range tx.Inputs {
			fmt.Printf("    - TxInput: ID: %x, Out: %d, Sig: %s\n", in.ID, in.Out, in.Sig)
		}

		// Print outputs for the transaction.
		fmt.Println("  Outputs:")
		for _, out := range tx.Outputs {
			fmt.Printf("    - TxOutput: Value: %d, PubKey: %s\n", out.Value, out.PubKey)
		}
		fmt.Println()
	}

	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("Nonce: %d\n", b.Nonce)
	fmt.Println("-----------")
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	block.PrintBlock()
	return block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
