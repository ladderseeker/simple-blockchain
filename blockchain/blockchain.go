package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First transaction from Genesis"
)

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

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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

func (chain *Blockchain) FindUTXOs(address string) []*UTXO {
	var UTXOs []*UTXO

	spentTXNs := make(map[string][]int)

	iterator := chain.Iterator()
	for iterator.HasNext() {
		block := iterator.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Check outputs and find unspent outputs
			for outIdx, output := range tx.Outputs {
				if spentTXNs[txID] != nil {
					// Check if this output index is in spentIdxs
					if isOutIdxSpent(outIdx, spentTXNs[txID]) == true {
						continue
					}
				}

				if output.CanBeUnLocked(address) {
					UTXOs = append(UTXOs, &UTXO{txID, outIdx, output})
				}
			}

			// Loop through input, mark linked output spent
			if !tx.IsCoinBase() {
				for _, input := range tx.Inputs {
					if input.CanUnLock(address) {
						inputTxID := hex.EncodeToString(input.ID)
						spentTXNs[inputTxID] = append(spentTXNs[inputTxID], input.Out)
					}
				}
			}
		}
	}
	return UTXOs
}

func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, []*UTXO) {
	accumulated := 0
	totalUTXOs := chain.FindUTXOs(address)
	var toSpendUTXOs []*UTXO

	for _, utxo := range totalUTXOs {
		toSpendUTXOs = append(toSpendUTXOs, utxo)
		accumulated += utxo.OutPut.Value
		if accumulated > amount {
			break
		}
	}
	return accumulated, toSpendUTXOs
}

func isOutIdxSpent(outId int, spentIdxs []int) bool {
	for _, spentId := range spentIdxs {
		if spentId == outId {
			return true
		}
	}
	return false
}

func InitBlockChain(address string) *Blockchain {
	if DBExists(dbFile) {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	Handle(err)
	return &Blockchain{lastHash, db}
}

func ContinueBlockchain(address string) *Blockchain {
	if DBExists(dbFile) == false {
		fmt.Println("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
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

	return &Blockchain{lastHash, db}
}

func DBExists(db string) bool {
	if _, err := os.Stat(db); os.IsNotExist(err) {
		return false
	}
	return true
}
