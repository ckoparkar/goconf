package main

type Store interface {
	GetAllKV() <-chan KV
	GetACL(token string) []string
}
