package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//Modify Block to store nonce so validation func can be implemented
/* Part 4: Edit Block Struct (block.go): Replace 'Data' with an array of Txns
Each block must have at least one tx. Can have many.
Edit Create Block and Genesis Funcs*/

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

/* Part 4: Our POW algo needs to consider the transactions in a block
so we need to create a new funct which allows us to use a hashing mechanism
 to provide a unique representation of all our txns combined */

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Add 0 for initial nonce (after prevHash)
// Modify CreateBlock so that it runs the pow algo on each block we create
// Execute the run func on that pow which will return the nonce and hash
// Put the nonce and hash into the block structure
// Return the block
// Now we can pass around data properly :)
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
