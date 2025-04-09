package blockchain

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
)

const dbPath = "./tmp/blocks"

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (iterator *BlockchainIterator) HasNext() bool {
	return !(len(iterator.CurrentHash) == 0)
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)

		err = item.Value(func(val []byte) error {
			block = block.Deserialize(val)
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

func InitBlockChain() *Blockchain {
	var lastHash []byte

	// Init DB
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)

	// Retrieve lastHash or create Genesis if no existing chain
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); errors.Is(err, badger.ErrKeyNotFound) {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")

			err := txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)
			Handle(err)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)

			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			Handle(err)
			return err
		}
	})

	Handle(err)

	blockchain := &Blockchain{lastHash, db}
	return blockchain
}
