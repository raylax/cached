package main

import (
	"github.com/dgraph-io/badger/v3"
	"time"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("data").WithLoggingLevel(badger.WARNING))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte("answer"), []byte("42")).
			WithTTL(10 * time.Second)
		return txn.SetEntry(entry)
	})

	if err != nil {
		panic(err)
	}

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			println(string(val))
			return nil
		})
	})

	if err != nil {
		panic(err)
	}

}
