package main

import (
	"encoding/json"
	"fmt"
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

	kvs := filterKV(filterKV(s.store.GetAllKV(), matcher), aclFilter)
	j, _ := json.Marshal(kvs)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(j))
}

func (s *Server) servePostKV(w http.ResponseWriter, r *http.Request) {
	// create a filter to remove kv's protected by ACL
	token := r.URL.Query().Get("token")
	acls := s.store.GetACL(token)
	aclFilter := DoesNotStartWithMatcher{prefixes: acls}

	var kvs []KV
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &kvs)

	kvs = filterKV(kvs, aclFilter)
	if err := s.store.SetKVs(kvs); err != nil {
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
		matcher := StartsWithMatcher{prefix: k}
		kvs := filterKV(filterKV(s.store.GetAllKV(), matcher), aclFilter)
		if err := s.store.DeleteKVs(kvs); err != nil {
			log.Println("[ERR] " + err.Error())
		}
	}
}
