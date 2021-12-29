package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	data map[string]interface{}
	hash string
	previousHash string
	timestamp time.Time
	pow int
}

func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHashBytes := sha256.Sum256([]byte(blockData))
	blockHash := fmt.Sprintf("%x", blockHashBytes)
	log.Tracef("Calculated hash: %s", blockHash)
	return blockHash
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

type BlockChain struct {
	genesisBlock Block
	chain []Block
	difficulty int
}

func (bc *BlockChain) AddBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from": from,
		"to": to,
		"amount": amount,
	}
	lastBlock := bc.chain[len(bc.chain) - 1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(bc.difficulty)
	bc.chain = append(bc.chain, newBlock)
}

func (bc BlockChain) isValid() bool {
	for i := range bc.chain[1:] {
		previousBlock := bc.chain[i]
		currentBlock := bc.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func NewBlockChain(difficulty int) *BlockChain {
	genesisBlock := Block{
		hash: "0",
		timestamp: time.Now(),
	}

	return &BlockChain{
		genesisBlock: genesisBlock,
		chain:        []Block{genesisBlock},
		difficulty:   difficulty,
	}
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := NewBlockChain(1)

	log.Debug("Block chain created")
	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.AddBlock("Alice", "Bob", 5)
	log.Debug("Added block transaction from Alice to Bob")

	blockchain.AddBlock("John", "Bob", 2)
	log.Debug("Added block transaction from John to Bob")

	// check if the blockchain is valid; expecting true
	log.Debug("Validating blockchain")
	log.Info(blockchain.isValid())
}
