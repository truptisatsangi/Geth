package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// validators keeps track of open validators and balances
var validators = make(map[string]uint)
var chain [1000]Block
var chainLength int

type Block struct {
	Number    uint64
	Hash      string
	PrevHash  string
	Timestamp uint64
	Nonce     uint64
	Data      string
	Validator string
}

// calculateBlockHash calculates the hash of a block.
func calculateBlockHash(block Block) string {
	record := strconv.FormatUint(block.Number, 10) +
		strconv.FormatUint(block.Timestamp, 10) +
		block.PrevHash + block.Data +
		strconv.FormatUint(block.Nonce, 10)

	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, address string) (Block, error) {
	newBlock := Block{
		Number:    oldBlock.Number + 1,
		Timestamp: uint64(time.Now().Unix()),
		PrevHash:  oldBlock.Hash,
		Data:      "0x",
		Nonce:     0,
		Validator: address,
	}
	newBlock.Hash = calculateBlockHash(newBlock)

	return newBlock, nil
}

// stakeEth increases the validator's stake
func stakeEth(amount uint, address string) {
	validators[address] += amount
}

// verifyBlock checks if a block is valid by comparing it with the last block in the chain
func verifyBlock(block Block) bool {
	if chainLength == 0 {
		return false
	}
	lastBlock := chain[chainLength-1]
	if block.Number != lastBlock.Number+1 || block.PrevHash != lastBlock.Hash || block.Hash != calculateBlockHash(block) {
		fmt.Println("Invalid block")
		return false
	}
	return true
}

func main() {
	var stakingPeriod = make(map[string]uint)
	var chosenValidator string

	// Simulate validator deposits
	validators["0"] = 32
	validators["1"] = 39
	validators["2"] = 96

	// Simulate staking durations (e.g., time difference from deposit to now)
	stakingPeriod["0"] = 110
	stakingPeriod["1"] = 80
	stakingPeriod["2"] = 10

	// Create and add genesis block
	genesisBlock := Block{
		Number:    0,
		Timestamp: uint64(time.Now().Unix()),
		PrevHash:  "0x",
		Data:      "Genesis Block",
		Nonce:     0,
		Validator: "GENESIS",
	}
	genesisBlock.Hash = calculateBlockHash(genesisBlock)
	chain[0] = genesisBlock
	chainLength++

	// Choose validator with highest score
	maxScore := 0
	for address, deposit := range validators {
		score := int((deposit*50 + stakingPeriod[address]*50) / 100)
		if score > maxScore {
			maxScore = score
			chosenValidator = address
		}
	}
	fmt.Println("Chosen Validator:", chosenValidator)

	// Generate and verify a new block
	newBlock, err := generateBlock(chain[chainLength-1], chosenValidator)
	if err != nil {
		fmt.Println("Error generating block:", err)
		return
	}

	if verifyBlock(newBlock) {
		chain[chainLength] = newBlock
		chainLength++
		fmt.Println("Block added to the chain by", newBlock.Validator)
	} else {
		fmt.Println("Failed to add block: Verification failed")
	}
}
