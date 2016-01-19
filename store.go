package main

type Store interface {
	GetAllKV() <-chan KV
	GetACL(token string) []string

	SetKV(kv KV) error
}
