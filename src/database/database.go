package database

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	dbPath = filepath.Join(os.Getenv("HOME"), ".config/zeroctl/zero.db")
	DB     *bolt.DB
	Bucket = "cache"

	ErrBucketNotFound = errors.New("bucket not found")
	ErrKeyNotFound    = errors.New("key not found")
)

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
		_, err := tx.CreateBucketIfNotExists([]byte(Bucket))
		return err
	})
}

func CloseBoltDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetValue(key string) (string, error) {
	var value string
	err := DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		if bucket == nil {
			return ErrBucketNotFound
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return ErrKeyNotFound
		}
		value = string(v)
		return nil
	})
	return value, err
}

func DeleteValue(key string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func StoreValue(key, value string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func StoreJsonData(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Bucket))
		if b == nil {
			return ErrBucketNotFound
		}
		return b.Put([]byte(key), data)
	})
}

func GetJsonData(key string, dest interface{}) error {
	return DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Bucket))
		if b == nil {
			return ErrBucketNotFound
		}

		data := b.Get([]byte(key))
		if data == nil {
			return ErrKeyNotFound
		}

		return json.Unmarshal(data, dest)
	})
}
