package kv

import (
	"log"

	"github.com/dgraph-io/badger/v2"
)

type Database struct {
	Path string
	DB   *badger.DB
}

// new db with path
func NewDB(path string) *Database {
	return &Database{
		Path: path,
		DB:   nil,
	}
}

// open db
func (d *Database) Open() error {
	opts := badger.DefaultOptions(d.Path)
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	d.DB = db

	return nil
}

// close db
func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) MultiPut(keys [][]byte, values [][]byte) error {
	return d.DB.Update(func(txn *badger.Txn) error {
		var err error
		for i := 0; i < len(keys); i++ {
			if err = txn.Set(keys[i], values[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Database) Put(key []byte, value []byte) error {
	return d.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}

func (d *Database) Delete(key []byte) error {
	return d.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
}

func (d *Database) Get(key []byte) ([]byte, error) {
	var result []byte
	err := d.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			result = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Database) Has(key []byte) (bool, error) {
	has := false
	err := d.DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
			return nil
		}
		has = true
		return nil
	})
	return has, err
}

func (d *Database) Update(key []byte, updateFunc func(value []byte) ([]byte, error)) error {
	return d.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		var oldValue []byte
		err = item.Value(func(val []byte) error {
			oldValue = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		newValue, err := updateFunc(oldValue)
		if err != nil {
			return err
		}
		return txn.Set(key, newValue)
	})
}

// get all value into a map
func (d *Database) GetAllValues() map[string]string {
	res := make(map[string]string)

	err := d.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				//fmt.Printf("key=%s, value=%s\n", k, v)
				res[string(k)] = string(v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Failed to iterator keys and values from the cache.", "error", err)
	}

	return res
}
