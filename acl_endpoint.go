package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) authorizeACL(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token != *aclMasterToken {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func (s *Server) serveGetACL(w http.ResponseWriter, r *http.Request) {
	acls := s.store.GetAllACL()
	j, _ := json.Marshal(acls)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(j))
}

func (s *Server) servePostACL(w http.ResponseWriter, r *http.Request) {
	var acls []KV
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &acls)

	if err := s.store.SetACLs(acls); err != nil {
		log.Println("[ERR] " + err.Error())
	}
}

func (s *Server) serveDeleteACL(w http.ResponseWriter, r *http.Request) {
	var keys []string
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &keys)

	acls := make([]KV, 0)
	for _, k := range keys {
		acls = append(acls, KV{Key: k, Value: ""})
	}
	if err := s.store.DeleteACLs(acls); err != nil {
		log.Println("[ERR] ", err)
	}
}
