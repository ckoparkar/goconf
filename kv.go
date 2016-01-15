package main

import (
	"encoding/base64"
	"strings"
)

type KV struct {
	Key, Value string
}

func (k KV) StartsWith(prefix string) bool {
	if strings.HasPrefix(k.Key, prefix) {
		return true
	}
	return false
}

type KVMatcher interface {
	Match(kv KV) bool
}

type StartsWithMatcher struct {
	prefix string
}

func (m StartsWithMatcher) Match(kv KV) bool {
	return kv.StartsWith(m.prefix)
}

func filterKV(kvs <-chan KV, matcher KVMatcher) <-chan KV {
	out := make(chan KV)
	go func() {
		for kv := range kvs {
			if ok := matcher.Match(kv); ok {
				out <- kv
			}
		}
		close(out)
	}()
	return out
}

func base64ToStringKV(kv KV) KV {
	decKey, _ := base64.StdEncoding.DecodeString(kv.Key)
	decVal, _ := base64.StdEncoding.DecodeString(kv.Value)
	return KV{Key: string(decKey), Value: string(decVal)}
}

type mapperKVFunc func(kv KV) KV

func mapKV(kvs <-chan KV, fn mapperKVFunc) <-chan KV {
	out := make(chan KV)
	go func() {
		for kv := range kvs {
			out <- fn(kv)
		}
		close(out)
	}()
	return out
}
