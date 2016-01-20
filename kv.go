package main

import (
	"encoding/base64"
	"strings"
)

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (k KV) StartsWith(prefix string) bool {
	if strings.HasPrefix(k.Key, prefix) {
		return true
	}
	return false
}

// filter functions

type KVMatcher interface {
	Match(kv KV) bool
}

type StartsWithMatcher struct {
	prefix string
}

func (m StartsWithMatcher) Match(kv KV) bool {
	return kv.StartsWith(m.prefix)
}

type ExactMatcher struct {
	prefix string
}

func (m ExactMatcher) Match(kv KV) bool {
	return kv.Key == m.prefix
}

type DoesNotStartWithMatcher struct {
	prefixes []string
}

func (m DoesNotStartWithMatcher) Match(kv KV) bool {
	for _, p := range m.prefixes {
		if p != "" && kv.StartsWith(p) {
			return false
		}
	}
	return true
}

func filterKV(kvs []KV, matcher KVMatcher) []KV {
	filteredKVs := make([]KV, 0)
	for _, kv := range kvs {
		if ok := matcher.Match(kv); ok {
			filteredKVs = append(filteredKVs, kv)
		}
	}
	return filteredKVs
}

// mapper functions

func base64ToStringKV(kv KV) KV {
	decKey, _ := base64.StdEncoding.DecodeString(kv.Key)
	decVal, _ := base64.StdEncoding.DecodeString(kv.Value)
	return KV{Key: string(decKey), Value: string(decVal)}
}

func stringKVToBase64(kv KV) KV {
	encKey := base64.StdEncoding.EncodeToString([]byte(kv.Key))
	encVal := base64.StdEncoding.EncodeToString([]byte(kv.Value))
	return KV{Key: string(encKey), Value: string(encVal)}
}

type mapperKVFunc func(kv KV) KV

func mapKV(kvs []KV, fn mapperKVFunc) []KV {
	mappedKVs := make([]KV, 0)
	for _, kv := range kvs {
		mappedKVs = append(mappedKVs, kv)
	}
	return mappedKVs
}
