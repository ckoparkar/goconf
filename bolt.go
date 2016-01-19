package main

import (
	"strings"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	*bolt.DB
}

func NewBoltDB(path string) (*BoltDB, error) {
	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	return &BoltDB{db}, nil
}

func (db *BoltDB) GetAllKV() <-chan KV {
	return db.getAllFromBucket("kv")
}

func (db *BoltDB) GetAllACL() <-chan KV {
	return db.getAllFromBucket("acl")
}

func (db *BoltDB) getAllFromBucket(bucket string) <-chan KV {
	out := make(chan KV, 10)
	go func() {
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(bucket))
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				out <- KV{Key: string(k), Value: string(v)}
			}

			close(out)
			return nil
		})
	}()
	return out
}

func (db *BoltDB) GetACL(token string) []string {
	var aclByte []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("acl"))
		aclByte = b.Get([]byte(token))
		return nil
	})
	if aclByte == nil {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("acl"))
			aclByte = b.Get([]byte("anonymous"))
			return nil
		})
	}
	return strings.Split(string(aclByte), ",")
}

func (db *BoltDB) SetKV(kv KV) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("kv"))
		if err != nil {
			return err
		}
		if err := b.Put([]byte(kv.Key), []byte(kv.Value)); err != nil {
			return err
		}
		return nil
	})
	return err
}
