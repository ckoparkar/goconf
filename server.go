package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	redis "gopkg.in/redis.v3"
)

var (
	kvPattern = regexp.MustCompile("/v1/kv/?.*")
)

type Server struct {
	rconn *Redis
}

func NewServer() (*Server, error) {
	rconn, err := NewRedisClient(&redis.Options{Addr: "localhost:6379"})
	if err != nil {
		return nil, err
	}
	return &Server{rconn: rconn}, nil
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
	token := r.URL.Query().Get("token")
	tokenExists, _ := s.rconn.HExists("acl", token).Result()
	if token == "" || !tokenExists {
		token = "anonymous"
	}
	recurse := r.URL.Query().Get("recurse")
	prefix := strings.Replace(r.URL.Path, "/v1/kv/", "", -1)
	var matcher KVMatcher
	if recurse != "" {
		matcher = StartsWithMatcher{prefix: prefix}
	} else {
		matcher = ExactMatcher{prefix: prefix}
	}

	// filter out keys starting with acl
	aclString, _ := s.rconn.HGet("acl", token).Result()
	acls := strings.Split(aclString, ",")
	aclMatcher := DoesNotStartWithMatcher{prefixes: acls}

	kvs := make([]KV, 0)
	for kv := range filterKV(filterKV(mapKV(s.rconn.GetAllKV(), base64ToStringKV), matcher), aclMatcher) {
		kvs = append(kvs, kv)
	}
	j, _ := json.Marshal(kvs)
	fmt.Fprintln(w, string(j))
}
