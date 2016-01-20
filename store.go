package main

import "io"

type Store interface {
	GetAllKV() []KV
	GetAllACL() []KV

	GetACL(token string) []string

	SetKVs(kvs []KV) error
	SetACLs(kvs []KV) error

	DeleteKVs(kvs []KV) error
	DeleteACLs(kvs []KV) error

	Backup(w io.Writer) (int, error)
}
