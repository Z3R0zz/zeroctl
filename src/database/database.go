package database

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

var dbPath = filepath.Join(os.Getenv("HOME"), ".config/zeroctl/zero.db")
var DB *bolt.DB

func InitBoltDB() error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}

	var err error
	DB, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("cache"))
		return err
	})
}

func CloseBoltDB() {
	if DB != nil {
		DB.Close()
	}
}

func StoreJsonData(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cache"))
		if b == nil {
			return errors.New("cache bucket does not exist")
		}
		return b.Put([]byte(key), data)
	})
}

func GetJsonData(key string, dest interface{}) error {
	return DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cache"))
		if b == nil {
			return errors.New("cache bucket does not exist")
		}

		data := b.Get([]byte(key))
		if data == nil {
			return errors.New("key not found")
		}

		return json.Unmarshal(data, dest)
	})
}
