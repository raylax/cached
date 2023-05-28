package database

import "github.com/dgraph-io/badger/v3"

type engineImpl struct {
	badger *badger.DB
}

func (e *engineImpl) Close() error {
	return e.badger.Close()
}
