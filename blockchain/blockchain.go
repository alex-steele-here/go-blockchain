package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte //Stores last hash of last block in chain
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts) //New badger db
	Handle(err)

	//Because we are initialising the blockchain, we need write capabilities: Hence, Update func
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(nil)
			return err //Return 1st err and handle 2nd err above
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db} //Create new blockchain in memory
	return &blockchain                     //This way we can use it further in our app
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	//Execute a read-only type of txn on our BadgerDB (call the database from our blockchain by calling on chain.Database)
	//Call View (read-only), which takes in a closure with a pointer to a Badger transaction
	//Returns an err

	err := chain.Database.View(func(txn *badger.Txn) error {
		//Get the current last hash out of the database (the hash of the last block in our database)
		//Call txn.Get to get the item, then unwrap the value from item and put it into our lastHash var
		//Return error if there is one
		//Handle the first err if there is one

		//New code from Badger DB below
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
		//New code ends here
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)
	//With our new block now created, we want to do a read/write type txn on our db
	//so we can put the new block into the database and assign the new block's hash to our last hash key
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
	//Succesfully created a layer of persistence for our blockchain
	//However, lost the ability to go through our blockchain and print it out like we have before
	//All the blocks are in the data layer. Can't just print them out
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	//Because we are starting with the last hash of our bc we will be iterating backwards through the blocks
	//starting with the newest and working our way back to the genesis block
	//Create next func for this

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
