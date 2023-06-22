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

// Take the data from the block
// create a counter (nonce) which starts at 0
// create a hash of the data plus the counter
// check the hash to see if it meets a set of requirements
// Requirements:
// The First few bytes must contain 0s

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof basically creates a target and pairs it with a given block
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// InitData (along with ToHex) will replace the old DeriveHash func
// In DeriveHash: Hash derived from Data + PrevHash
// In InitData + ToHex: Hash derived from Data, PrevHash, Nonce, Diff.
// Nonce and Diff need to be cast in an int64 when calling ToHex on them
// Chapter 4: Replace Data with HashTransactions
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// Run is the main computational method
// It runs on the pow and outputs an int (nonce) and a slice of bytes (hash)
// Basically create an infinite loop
// Prep data using InitData
// Hash it into a sha256 format - Printed out to see the process
// Convert hash into big int
// Compare pow target with new big int
// If == -1, break out of loop (hash is less than target means block is signed)
// Else, increment nonce, create new hash
// Println for some space (outside loop)
// Return the nonce and hash
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}

	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex is a utility func that turns an int64 into a slice of bytes
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}
