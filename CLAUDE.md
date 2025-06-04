# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of a blockchain with Proof-of-Work consensus, transaction management, and peer-to-peer networking capabilities. The project implements a complete blockchain system including:

- **Blockchain Core**: Block creation, mining, transaction verification using elliptic curve cryptography
- **UTXO Management**: Unspent Transaction Output tracking for efficient balance calculations
- **Wallet System**: Key generation, address creation, and transaction signing
- **Network Layer**: P2P communication with transaction mempool and block synchronization
- **CLI Interface**: Command-line tools for blockchain operations

## Build and Development Commands

```bash
# Build the blockchain executable
go build -o blockchainexe .

# Run tests (uses Ginkgo/Gomega testing framework)
go test ./merkletree/...

# Run with specific node ID (required environment variable)
NODE_ID=3000 ./blockchainexe <command>
```

## Architecture

### Core Components

- **blockchainstruct/**: Core blockchain data structures and operations
  - `blockchain.go`: Main blockchain logic, mining, transaction verification
  - `block.go`: Block structure with Merkle tree root and proof-of-work
  - `transaction.go`: Transaction inputs/outputs with ECDSA signatures
  - `utxo_set.go`: UTXO set management for efficient balance queries
  - `proof-of-work.go`: SHA-256 based proof-of-work implementation

- **cli/**: Command-line interface and network server
  - `cli.go`: Main CLI dispatcher with all available commands
  - `server.go`: P2P networking layer with mempool and block propagation
  - Individual command handlers for blockchain operations

- **wallets/**: Cryptographic wallet management
  - ECDSA key pair generation and address derivation
  - Base58 encoding for Bitcoin-compatible addresses

- **merkletree/**: Merkle tree implementation for block transaction verification

### Network Architecture

The system supports multi-node operation with:
- **Node Communication**: TCP-based P2P protocol on configurable ports
- **Transaction Mempool**: Pending transactions before mining
- **Block Synchronization**: Automatic blockchain sync between nodes
- **Mining Coordination**: Distributed mining with configurable miner addresses

### Data Storage

- **BoltDB**: Persistent blockchain storage with node-specific databases (`blockchain_<nodeID>.db`)
- **Wallet Files**: Encrypted wallet storage per node
- **UTXO Set**: Cached unspent outputs for performance

## Node Operation

Each node requires a unique `NODE_ID` environment variable. The system supports:
- **Genesis Node**: `localhost:3000` (default known node)
- **Mining Nodes**: Nodes that can mine blocks when mempool reaches threshold
- **Regular Nodes**: Transaction relay and blockchain synchronization

## CLI Commands

All commands require `NODE_ID` environment variable to be set:

- `createblockchain -address <address>`: Initialize blockchain with genesis block
- `createwallet`: Generate new wallet and address
- `getbalance -address <address>`: Query address balance
- `sendcoin -from <addr> -to <addr> -amount <int> -mine`: Transfer coins
- `startnode -miner <address>`: Start P2P node with optional mining
- `printchain`: Display entire blockchain
- `listaddresses`: Show all wallet addresses
- `reindexutxo`: Rebuild UTXO set

## Testing

Uses Ginkgo BDD testing framework. Tests are located in `merkletree/merkle_tree_test.go` and follow the `*_test.go` naming convention.