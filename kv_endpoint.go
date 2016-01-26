package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) serveGetKV(w http.ResponseWriter, r *http.Request) {
	// create a filter to remove kv's protected by ACL
	token := r.URL.Query().Get("token")
	acls := s.store.GetACL(token)
	aclFilter := DoesNotStartWithMatcher{prefixes: acls}

	// create a matcher depending on recurse option
	recurse := r.URL.Query().Get("recurse")
	prefix := kvPattern.ReplaceAllString(r.URL.Path, "$2")
	var matcher KVMatcher
	if recurse != "" {
		matcher = StartsWithMatcher{prefix: prefix}
	} else {
		matcher = ExactMatcher{prefix: prefix}
	}

	kvs := make([]KV, 0)
	for kv := range filterKV(filterKV(s.store.GetAllKV(), matcher), aclFilter) {
		kvs = append(kvs, kv)
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	json.NewEncoder(gz).Encode(kvs)
	gz.Close()
}

func (s *Server) servePostKV(w http.ResponseWriter, r *http.Request) {
	// create a filter to remove kv's protected by ACL
	token := r.URL.Query().Get("token")
	acls := s.store.GetACL(token)
	aclFilter := DoesNotStartWithMatcher{prefixes: acls}

	var kvs []KV
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &kvs)

	kvChan := make(chan KV, 10)
	go func() {
		for _, kv := range kvs {
			kvChan <- kv
		}
		close(kvChan)
	}()

	filteredKVs := make([]KV, 0)
	for kv := range filterKV(kvChan, aclFilter) {
		filteredKVs = append(filteredKVs, kv)
	}
	if err := s.store.SetKVs(filteredKVs); err != nil {
		log.Println("[ERR] " + err.Error())
	}
}

func (s *Server) serveDeleteKV(w http.ResponseWriter, r *http.Request) {
	// create a filter to remove kv's protected by ACL
	token := r.URL.Query().Get("token")
	acls := s.store.GetACL(token)
	aclFilter := DoesNotStartWithMatcher{prefixes: acls}

	var keys []string
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &keys)

	for _, k := range keys {
		kvs := make([]KV, 0)
		matcher := StartsWithMatcher{prefix: k}
		for kv := range filterKV(filterKV(s.store.GetAllKV(), matcher), aclFilter) {
			kvs = append(kvs, kv)
		}
		if err := s.store.DeleteKVs(kvs); err != nil {
			log.Println("[ERR] " + err.Error())
		}

	}
}
