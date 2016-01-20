package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	kvPattern = regexp.MustCompile("/v1/kv/?.*")
)

type Server struct {
	store Store
}

func NewServer() (*Server, error) {
	// db, err := NewRedisClient("localhost:6379")
	db, err := NewBoltDB("conf.db")
	if err != nil {
		return nil, err
	}
	return &Server{store: db}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case kvPattern.MatchString(r.URL.Path):
		switch r.Method {
		case "GET":
			s.serveGetKV(w, r)
			return
		case "POST":
			s.servePostKV(w, r)
			return
		default:
			http.Error(w, "No route found.", http.StatusNotFound)
			return
		}
	case r.URL.Path == "/v1/acl":
		switch r.Method {
		case "GET":
			s.serveGetACL(w, r)
			return
		case "POST":
			s.servePostACL(w, r)
			return
		default:
			http.Error(w, "No route found.", http.StatusNotFound)
			return
		}
	case r.URL.Path == "/backup":
		s.serveBackup(w, r)
		return
	default:
		http.Error(w, "No route found.", http.StatusNotFound)
		return
	}
}

func (s *Server) serveGetKV(w http.ResponseWriter, r *http.Request) {
	// create a filter to remove kv's protected by ACL
	token := r.URL.Query().Get("token")
	acls := s.store.GetACL(token)
	aclFilter := DoesNotStartWithMatcher{prefixes: acls}

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
	for kv := range filterKV(filterKV(s.store.GetAllKV(), matcher), aclFilter) {
		kvs = append(kvs, kv)
	}
	j, _ := json.Marshal(kvs)
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

	kvChan := make(chan KV, 10)
	go func() {
		for _, kv := range kvs {
			kvChan <- kv
		}
		close(kvChan)
	}()

	for kv := range filterKV(kvChan, aclFilter) {
		if err := s.store.SetKV(kv); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) serveGetACL(w http.ResponseWriter, r *http.Request) {
	var acls []KV
	for acl := range s.store.GetAllACL() {
		acls = append(acls, acl)
	}
	j, _ := json.Marshal(acls)
	fmt.Fprintln(w, string(j))
}

func (s *Server) servePostACL(w http.ResponseWriter, r *http.Request) {
	var acls []KV
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &acls)

	for _, acl := range acls {
		if err := s.store.SetACL(acl); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) serveBackup(w http.ResponseWriter, r *http.Request) {
	n, err := s.store.Backup(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
	w.Header().Set("Content-Length", strconv.Itoa(n))
}
