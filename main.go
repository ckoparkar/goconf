package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listenAddr     = flag.String("listen", ":8080", "HTTP listen address.")
	aclMasterToken = flag.String("master-token", "master", "ACL master token")
)

func main() {
	flag.Parse()
	s, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %s\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
