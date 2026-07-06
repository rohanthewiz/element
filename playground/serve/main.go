// Command serve hosts the playground directory locally for development.
// Run playground/build.sh first, then: go run ./playground/serve
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dir := flag.String("dir", "playground", "directory to serve")
	flag.Parse()
	log.Printf("playground on http://localhost%s (serving %s)", *addr, *dir)
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir(*dir))))
}
