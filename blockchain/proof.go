package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 24

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	return &ProofOfWork{b, target}
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWork) InitNonce(nonce int) []byte {
	return bytes.Join([][]byte{
		pow.Block.PrevHash,
		pow.Block.HashTransactions(),
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	}, []byte{})
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	fmt.Println("Performing pow...")
	for nonce < math.MaxInt64 {
		data := pow.InitNonce(nonce)
		hash = sha256.Sum256(data)

		//fmt.Printf("\rPow: %x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	data := pow.InitNonce(pow.Block.Nonce)
	hash := sha256.Sum256(data)

	var intHash big.Int
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
