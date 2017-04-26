// demo1 shows an example of a simple http service
// with profiling information exposed.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	// Importing net/http/pprof initializes
	// handlers for pprof under /debug/pprof
	_ "net/http/pprof"

	"sync"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8081", nil)
}

var globalCounter struct {
	sync.Mutex
	count uint64
}

func incrementCounter() uint64 {
	globalCounter.Lock()
	defer globalCounter.Unlock()
	globalCounter.count++
	latest := globalCounter.count

	return latest
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, fmt.Sprint(incrementCounter()))
	if err != nil {
		log.Println("error writing response:", err)
	}
}
