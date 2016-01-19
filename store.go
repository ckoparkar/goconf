package main

type Store interface {
	GetAllKV() <-chan KV
	GetAllACL() <-chan KV

	GetACL(token string) []string

	SetKV(kv KV) error
	SetACL(kv KV) error
}
