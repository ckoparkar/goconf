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
	//	token := r.URL.Query().Get("token")
	recurse := r.URL.Query().Get("recurse")
	prefix := strings.Replace(r.URL.Path, "/v1/kv/", "", -1)
	var matcher KVMatcher
	if recurse != "" {
		matcher = StartsWithMatcher{prefix: prefix}
	} else {
		matcher = ExactMatcher{prefix: prefix}
	}

	kvs := make([]KV, 0)
	for kv := range filterKV(mapKV(s.rconn.GetAllKV(), base64ToStringKV), matcher) {
		kvs = append(kvs, kv)
	}
	j, _ := json.Marshal(kvs)
	fmt.Fprintln(w, string(j))
}
