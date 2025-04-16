# Simple Blockchain in Go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A demonstration of a Proof‑of‑Work blockchain with UTXO transaction model, wallet management (ECDSA & Base58), and a simple command‑line interface implemented in Go.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting‑started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
- [Usage](#usage)
    - [Building the Binary](#building‑the‑binary)
    - [CLI Commands](#cli‑commands)
- [Design & Architecture](#design‑--architecture)
    - [Block & Blockchain](#block‑--blockchain)
    - [Proof of Work](#proof‑of‑work)
    - [UTXO Transaction Model](#utxo‑transaction‑model)
    - [Wallet & Address Generation](#wallet‑--address‑generation)
    - [Data Persistence (BadgerDB)](#data‑persistence‑badgerdb)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Proof-of-Work** consensus with adjustable difficulty (24) 
- **UTXO‑based transactions**, including coinbase (mining reward) and standard transfers
- **ECDSA‑backed wallets** with SHA‑256 & RIPEMD‑160 hashing, Base58Check encoding
- **Block storage** using BadgerDB in `./tmp/blocks`
- **Command‑line interface** enabling blockchain creation, wallet management, and transactions

---

## Getting Started

### Prerequisites

- Go 1.16 or later
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/ladderseeker/simple-blockchain.git
cd simple-blockchain

# Fetch dependencies and build
go mod download
go build -o simple-blockchain
```

---

## Usage

### Building the Binary

```bash
go build -o simple-blockchain main.go
```

### CLI Commands

Once built, run:

```bash
./simple-blockchain [command] [flags]
```

| Command                               | Description                                                                      |
|---------------------------------------|----------------------------------------------------------------------------------|
| `createWallet`                        | Generate a new ECDSA wallet and print its Base58Check address                    |
| `listAddress`                         | List all addresses stored in the wallet file                                     |
| `createBlockchain -address ADDRESS`   | Initialize a new blockchain; mine the genesis block to the given address         |
| `getBalance -address ADDRESS`         | Display the balance (sum of UTXOs) for the specified address                     |
| `send -from FROM -to TO -amount AMOUNT` | Create a transaction transferring coins and mine a new block                    |
| `printChain`                          | Iterate and print all blocks, including transaction details and Proof‑of‑Work validation |

_To see usage details for a command run it without required flags._

---

## Design & Architecture

### Block & Blockchain

- Each **Block** contains:
    - `PrevHash` – hash of the previous block
    - `Transactions` – list of `*Transaction` (incl. coinbase)
    - `Nonce` & `Hash` computed via PoW
- **Blockchain** struct wraps the latest block hash and a BadgerDB instance for persistence.
- Blocks are serialized/deserialized with Go’s `encoding/gob` 

### Proof of Work

- **Difficulty**: leading zeros target set by `Difficulty = 24`.
- **PoW Data**: concatenates `PrevHash`, tx Merkle hash, `nonce`, and `difficulty`.
- Loop increments nonce until SHA‑256 hash < target.
- Validation re‑computes hash and compares to target. 

### UTXO Transaction Model

- **Transaction**:
    - `Inputs` reference previous outputs (`TxInput`) with signature match (`Sig == PubKey`).
    - `Outputs` (`TxOutput`) carry a `Value` and `PubKey` (recipient address).
- **CoinbaseTx** creates a single input with `Out = -1` and rewards `100` coins.
- **UTXO Set**: traversed via `FindUTXOs` to collect unspent outputs; spent outputs are tracked in a map.
- Change is returned when `totalUTXO > amount`.

### Wallet & Address Generation

- **Key Pair**: ECDSA on P‑256 curve.
- **Address**:
    1. SHA‑256 hash of public key
    2. RIPEMD‑160 of result
    3. Prepend version byte `0x00`
    4. Append 4‑byte checksum (double SHA‑256)
    5. Base58 encode the final payload

### Data Persistence (BadgerDB)

- Blocks stored under key = block hash; special key `"lh"` holds the latest block hash.
- DB directory: `./tmp/blocks`.
- Uses Dgraph’s BadgerDB for a lightweight, fast key‑value store.

---

## Examples

```bash
# Create a new wallet
./simple-blockchain createWallet
# ➜ New address is 1HAbC...

# List all addresses
./simple-blockchain listAddress
# ➜ 1HAbC...
# ➜ 1XYZ9...

# Initialize blockchain
./simple-blockchain createBlockchain -address 1HAbC...
# ➜ Genesis created
# ➜ Finished Creating chain

# Check balance
./simple-blockchain getBalance -address 1HAbC...
# ➜ Balance of 1HAbC...: 100

# Send transaction
./simple-blockchain send -from 1HAbC... -to 1XYZ9... -amount 10
# ➜ Transaction send 10 from 1HAbC... to 1XYZ9... success...

# Print the entire chain
./simple-blockchain printChain
# ➜ Block details and PoW validation...
```

---

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.

