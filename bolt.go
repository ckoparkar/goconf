package main

import (
	"io"
	"strings"

	"github.com/boltdb/bolt"
)

type BoltStore struct {
	*bolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	b, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	return &BoltStore{b}, nil
}

func (b *BoltStore) GetAllKV() <-chan KV {
	return b.getAllFromBucket("kv")
}

func (b *BoltStore) GetAllACL() <-chan KV {
	return b.getAllFromBucket("acl")
}

func (b *BoltStore) getAllFromBucket(bucket string) <-chan KV {
	out := make(chan KV, 10)
	go func() {
		b.View(func(tx *bolt.Tx) error {
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

func (b *BoltStore) GetACL(token string) []string {
	var aclByte []byte
	b.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("acl"))
		aclByte = b.Get([]byte(token))
		return nil
	})
	if aclByte == nil {
		b.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("acl"))
			aclByte = b.Get([]byte("anonymous"))
			return nil
		})
	}
	return strings.Split(string(aclByte), ",")
}

func (b *BoltStore) SetKV(kv KV) error {
	return b.setInBucket(kv, "kv")
}

func (b *BoltStore) SetACL(kv KV) error {
	return b.setInBucket(kv, "acl")
}

func (b *BoltStore) setInBucket(kv KV, bucket string) error {
	err := b.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
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

func (b *BoltStore) DeleteKV(kv KV) error {
	return b.deleteInBucket(kv, "kv")
}

func (b *BoltStore) deleteInBucket(kv KV, bucket string) error {
	err := b.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		if err := b.Delete([]byte(kv.Key)); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (b *BoltStore) Backup(w io.Writer) (int, error) {
	var n int64
	err := b.View(func(tx *bolt.Tx) error {
		var err error
		n, err = tx.WriteTo(w)
		return err
	})
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
