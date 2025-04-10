package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const reward = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

type TxOutput struct {
	Value  int
	PubKey string
}

func CoinbaseTx(toAddress, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}

	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{reward, toAddress}

	return &Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
}

func (in *TxInput) CanUnLock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnLocked(data string) bool {
	return out.PubKey == data
}

type UTXO struct {
	TxID   string
	OutIdx int
	OutPut TxOutput
}

func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	totalAmount, utxos := chain.FindSpendableOutputs(from, amount)
	if totalAmount < amount {
		log.Panic("Error, Not enough funds!")
	}

	var inputs []TxInput
	var outputs []TxOutput
	var transaction *Transaction

	for _, utxo := range utxos {
		txId, err := hex.DecodeString(utxo.TxID)
		Handle(err)
		inputs = append(inputs, TxInput{txId, utxo.OutIdx, from})
	}

	outputs = append(outputs, TxOutput{amount, to})

	if totalAmount-amount > 0 {
		outputs = append(outputs, TxOutput{totalAmount - amount, from})
	}

	transaction = &Transaction{nil, inputs, outputs}
	transaction.SetID()

	return transaction
}
