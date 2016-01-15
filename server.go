package main

import (
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
		} else {
			http.Error(w, "No route found.", http.StatusNotFound)
		}
	}
}

func (s *Server) serveGetKV(w http.ResponseWriter, r *http.Request) {
	//	recurse := r.URL.Query().Get("recurse")
	//	token := r.URL.Query().Get("token")
	prefix := strings.Replace(r.URL.Path, "/v1/kv/", "", -1)
	kvs1 := make([]KV, 0)
	matcher := StartsWithMatcher{prefix: prefix}
	for kv := range filterKV(mapKV(s.rconn.GetAllKV(), base64ToStringKV), matcher) {
		kvs1 = append(kvs1, kv)
	}
	fmt.Fprintln(w, kvs1)
}
