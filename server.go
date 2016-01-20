package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	kvPattern  = regexp.MustCompile("(/v1/kv/?)(.*)")
	aclPattern = regexp.MustCompile("/v1/acl/?")
)

type Server struct {
	store Store
}

func NewServer() (*Server, error) {
	db, err := NewBoltStore("conf.db")
	if err != nil {
		return nil, err
	}
	return &Server{store: db}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/ui"):
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	case kvPattern.MatchString(r.URL.Path):
		switch r.Method {
		case "GET":
			s.serveGetKV(w, r)
			return
		case "POST":
			s.servePostKV(w, r)
			return
		case "DELETE":
			s.serveDeleteKV(w, r)
			return
		default:
			http.Error(w, "No route found.", http.StatusNotFound)
			return
		}
	case aclPattern.MatchString(r.URL.Path):
		switch r.Method {
		case "GET":
			s.authorizeACL(s.serveGetACL)(w, r)
			return
		case "POST":
			s.authorizeACL(s.servePostACL)(w, r)
			return
		case "DELETE":
			s.authorizeACL(s.serveDeleteACL)(w, r)
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

func (s *Server) serveBackup(w http.ResponseWriter, r *http.Request) {
	n, err := s.store.Backup(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="conf.db"`)
	w.Header().Set("Content-Length", strconv.Itoa(n))
}
