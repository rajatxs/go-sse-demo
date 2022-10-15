package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	hostname = "localhost"
	port     = 4040
)

func main() {
	var addr string = fmt.Sprintf("%s:%d", hostname, port)

	log.Printf("server starting at %s\n", addr)
	log.Fatal(
		"HTTP server error: ",
		http.ListenAndServe(addr, getServerMux()))
}
