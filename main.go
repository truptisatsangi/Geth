package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Node represents a validator node
type Node struct {
	Address     string
	Stake       uint
	StakePeriod uint
}

var chain []Block
var Nodes = map[string]*Node{}
var mu sync.Mutex

// Block represents a block in the chain
type Block struct {
	Number    uint64
	Hash      string
	PrevHash  string
	Timestamp uint64
	Nonce     uint64
	Data      string
	Validator string
}

// Update stake of nodes
func updateStake(id string, delta int) {
	mu.Lock()
	defer mu.Unlock()

	node := Nodes[id]
	node.Stake += uint(delta)
}

// calculateBlockHash calculates the hash of a block.
func calculateBlockHash(block Block) string {
	record := strconv.FormatUint(block.Number, 10) +
		strconv.FormatUint(block.Timestamp, 10) +
		block.PrevHash +
		block.Data +
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
		Validator: address,
		Data:      "0x",
		Nonce:     0,
	}
	newBlock.Hash = calculateBlockHash(newBlock)
	return newBlock, nil
}

func initiateNode(id string, address string, stake uint, stakePeriod uint) {
	node := &Node{Address: address, Stake: stake, StakePeriod: stakePeriod}
	Nodes[id] = node
}

// chooseValidator picks a validator based on weighted score
func chooseValidator(validators map[string]*Node) string {
	type scoreEntry struct {
		Address string
		Score   int
	}
	var scores []scoreEntry

	for _, node := range validators {
		score := int((node.Stake*50 + node.StakePeriod*50) / 100)
		scores = append(scores, scoreEntry{Address: node.Address, Score: score})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if len(scores) > 0 {
		return scores[0].Address
	}
	return ""
}

// generateAndValidateBlock creates a block and adds to chain if valid
func generateAndValidateBlock(validator string) {
	mu.Lock()
	defer mu.Unlock()
	lastBlock := chain[len(chain)-1]
	newBlock, _ := generateBlock(lastBlock, validator)
	chain = append(chain, newBlock)
	fmt.Printf("Validator %s accepted block %d\n", validator, newBlock.Number)

}

// verifyBlock checks block validity
func verifyBlock(block Block, node string) bool {
	lastBlock := chain[len(chain)-2]
	if block.Number != lastBlock.Number+1 || block.Hash != calculateBlockHash(block) {
		fmt.Printf("Invalid block detected by %s\n", node)
		return false
	}
	fmt.Printf("Node %s has verified block number %d sucessfully\n", node, block.Number)

	return true
}

func main() {
	genesisBlock := Block{
		Number:    0,
		Timestamp: uint64(time.Now().Unix()),
		PrevHash:  "0x",
		Data:      "Genesis Block",
		Nonce:     0,
	}
	genesisBlock.Hash = calculateBlockHash(genesisBlock)
	chain = append(chain, genesisBlock)

	// Initiate Validators
	initiateNode("1", "0x0000000000000000000000000000000000000001", 32, 110)
	initiateNode("2", "0x0000000000000000000000000000000000000002", 39, 80)
	initiateNode("3", "0x0000000000000000000000000000000000000003", 96, 10)

	// Choose Validator
	validator := chooseValidator(Nodes)
	generateAndValidateBlock(validator)

	// Verify block by other nodes
	for _, node := range Nodes {
		if node.Address != validator {
			isValid := verifyBlock(chain[len(chain)-1], node.Address)
			if !isValid {
				panic(fmt.Sprintf("Node %s detected an invalid block. Halting!", node.Address))

			}
		}
	}
	fmt.Print("A block is successfully added")
}
