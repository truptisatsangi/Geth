# Go Proof-of-Stake (PoS) Consensus Simulator

This project is a simple Go implementation of a Proof-of-Stake (PoS) consensus mechanism. It demonstrates core blockchain principles such as validator selection, block generation, and block verification in a minimal, readable format.

## 🚀 Features

- Genesis block initialization  
- Validator node registration with stake and staking period  
- Weighted validator selection algorithm  
- Block generation and hashing  
- Block verification by peer nodes  

## 🧠 How It Works

### 1. Genesis Block Creation
A genesis block is initialized and added to the blockchain.

### 2. Validator Initialization
Three validators are registered with:
- Unique address
- Stake amount
- Staking duration

### 3. Validator Selection
The validator with the highest weighted score is selected:

```score = (stake * 50 + stakePeriod * 50) / 100```
### 4. Block Generation
The selected validator generates a new block based on the previous block’s hash and metadata.

### 5. Block Verification
The remaining validators verify the new block’s integrity by checking hashes and block continuity.

## 📄 Code Structure

- `Node`: Represents a validator (address, stake, stakePeriod).
- `Block`: Contains block data, hash, validator info, etc.
- `chooseValidator`: Selects the best validator by score.
- `generateBlock`: Produces a block with a calculated hash.
- `verifyBlock`: Ensures the new block is valid and properly linked.
- `main`: Executes a single consensus round.

## ✅ Sample Output
```
Validator 0x0000000000000000000000000000000000000003 accepted block 1
Node 0x0000000000000000000000000000000000000001 has verified block number 1 successfully
Node 0x0000000000000000000000000000000000000002 has verified block number 1 successfully
A block is successfully added
```

## 📦 Dependencies

- Go (no third-party libraries required)

## ▶️ Running the Code

Make sure you have Go installed, then run:

```bash
go run main.go
