package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

type walletDTO struct {
	D   []byte
	Pub []byte
}

func (ws *Wallets) SaveFile() {
	// 1) Build a map[string]walletDTO
	dtos := make(map[string]walletDTO, len(ws.Wallets))
	for addr, w := range ws.Wallets {
		dtos[addr] = walletDTO{
			D:   w.PrivateKey.D.Bytes(),
			Pub: w.PublicKey,
		}
	}

	// 2) Gob‑encode that dto map
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(dtos); err != nil {
		log.Panic(err)
	}

	// 3) Write to disk (owner read/write; group/other read-only)
	if err := os.WriteFile(walletFile, buf.Bytes(), 0644); err != nil {
		log.Panic(err)
	}
}

func (ws *Wallets) LoadFile() error {
	// If the file doesn’t exist, bail out
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	// 1) Read the raw bytes
	data, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	// 2) Decode into our dto map
	var dtos map[string]walletDTO
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&dtos); err != nil {
		return err
	}

	// 3) Rebuild the real Wallets map
	ws.Wallets = make(map[string]*Wallet, len(dtos))
	for addr, dto := range dtos {
		// a) Reconstruct the private key
		priv := ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P256(),
			},
			D: new(big.Int).SetBytes(dto.D),
		}
		// b) Derive X,Y from D
		x, y := elliptic.P256().ScalarBaseMult(dto.D)
		priv.PublicKey.X = x
		priv.PublicKey.Y = y

		// c) Fill in Wallet.PublicKey (the byte form)
		w := &Wallet{
			PrivateKey: priv,
			PublicKey:  dto.Pub,
		}
		ws.Wallets[addr] = w
	}

	return nil
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAllAddress() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile()
	return &wallets, err
}
