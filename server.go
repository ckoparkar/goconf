package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	// redis "gopkg.in/redis.v3"
)

var (
	kvPattern = regexp.MustCompile("/v1/kv/?.*")
)

type Server struct {
	store Store
}

func NewServer() (*Server, error) {
	// db, err := NewRedisClient(&redis.Options{Addr: "localhost:6379"})
	db, err := NewBoltDB("conf.db")
	if err != nil {
		return nil, err
	}
	return &Server{store: db}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case kvPattern.MatchString(r.URL.Path):
		if r.Method == "GET" {
			s.serveGetKV(w, r)
			return
		} else {
			http.Error(w, "No route found.", http.StatusNotFound)
			return
		}
	}
}

func (s *Server) serveGetKV(w http.ResponseWriter, r *http.Request) {
	// create a matcher to filter out kv's protected by ACL
	token := r.URL.Query().Get("token")
	acls := s.store.GetACL(token)
	aclMatcher := DoesNotStartWithMatcher{prefixes: acls}

	// create a matcher depending on recurse option
	recurse := r.URL.Query().Get("recurse")
	prefix := strings.Replace(r.URL.Path, "/v1/kv/", "", -1)
	var matcher KVMatcher
	if recurse != "" {
		matcher = StartsWithMatcher{prefix: prefix}
	} else {
		matcher = ExactMatcher{prefix: prefix}
	}

	kvs := make([]KV, 0)
	// for kv := range filterKV(filterKV(mapKV(s.store.GetAllKV(), base64ToStringKV), matcher), aclMatcher) {
	//	kvs = append(kvs, kv)
	// }
	for kv := range filterKV(filterKV(s.store.GetAllKV(), matcher), aclMatcher) {
		kvs = append(kvs, kv)
	}
	fmt.Println(len(kvs))
	j, _ := json.Marshal(kvs)
	fmt.Fprintln(w, string(j))
}
