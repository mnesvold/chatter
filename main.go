package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	staticRoot = "./src/github.com/mnesvold/chatter/www"
)

var (
	port = flag.Int("port", 8000, "port to serve site over")
)

func main() {
	flag.Parse()

	server := NewServer()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(staticRoot))))
	http.Handle("/chat", server)

	log.Printf("Listening on port %d\n", *port)
	bind := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(bind, nil))
}
