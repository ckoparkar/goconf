package main

import (
	"io"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	dbKV  = []byte("kv")
	dbACL = []byte("acl")
)

type BoltStore struct {
	*bolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	b, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	store := &BoltStore{b}
	if err := store.initialize(); err != nil {
		return nil, err
	}
	return store, nil
}

func (b *BoltStore) initialize() error {
	tx, err := b.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create all the buckets
	if _, err := tx.CreateBucketIfNotExists(dbKV); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists(dbACL); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *BoltStore) GetAllKV() []KV {
	return b.getAllFromBucket(dbKV)
}

func (b *BoltStore) GetAllACL() []KV {
	return b.getAllFromBucket(dbACL)
}

func (b *BoltStore) getAllFromBucket(bucketName []byte) []KV {
	tx, _ := b.Begin(true)
	defer tx.Rollback()
	bucket := tx.Bucket(bucketName)
	cursor := bucket.Cursor()

	kvs := make([]KV, 0)
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		kvs = append(kvs, KV{Key: string(k), Value: string(v)})
	}
	return kvs
}

func (b *BoltStore) GetACL(token string) []string {
	tx, _ := b.Begin(true)
	defer tx.Rollback()

	bucket := tx.Bucket(dbACL)
	acl := bucket.Get([]byte(token))
	if acl == nil {
		acl = bucket.Get([]byte("anonymous"))
	}
	return strings.Split(string(acl), ",")
}

func (b *BoltStore) SetKVs(kvs []KV) error {
	return b.setInBucket(kvs, dbKV)
}

func (b *BoltStore) SetACLs(kvs []KV) error {
	return b.setInBucket(kvs, dbACL)
}

func (b *BoltStore) setInBucket(kvs []KV, bucketName []byte) error {
	tx, err := b.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketName)
	for _, kv := range kvs {
		if err := bucket.Put([]byte(kv.Key), []byte(kv.Value)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (b *BoltStore) DeleteKVs(kvs []KV) error {
	return b.deleteInBucket(kvs, dbKV)
}

func (b *BoltStore) DeleteACLs(kvs []KV) error {
	return b.deleteInBucket(kvs, dbACL)
}

func (b *BoltStore) deleteInBucket(kvs []KV, bucketName []byte) error {
	tx, err := b.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketName)
	for _, kv := range kvs {
		bucket.Delete([]byte(kv.Key))
	}
	return tx.Commit()
}

func (b *BoltStore) Backup(w io.Writer) (int, error) {
	tx, err := b.Begin(true)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	n, err := tx.WriteTo(w)
	if err != nil {
		return -1, err
	}
	return int(n), nil
}
