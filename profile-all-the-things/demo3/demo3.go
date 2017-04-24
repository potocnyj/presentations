// demo3 shows an example of enabling mutex profiling
// to identify long critical-sections in a simple http service.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

func main() {
	// Sample 1 out of every 100 mutex contention events.
	runtime.SetMutexProfileFraction(100)

	http.HandleFunc("/w", handleIncCounter)
	http.HandleFunc("/r", handleReadCounter)
	http.ListenAndServe("localhost:8080", nil)
}

var globalCounter struct {
	sync.Mutex
	count uint64
}

func incrementCounter() uint64 {
	globalCounter.Lock()
	defer globalCounter.Unlock()
	// This counter increment is in a critical section,
	// but it's a fast operation so it does not generate much contention.
	globalCounter.count++
	latest := globalCounter.count

	return latest
}

func handleIncCounter(w http.ResponseWriter, r *http.Request) {
	incrementCounter()
}

func handleReadCounter(w http.ResponseWriter, r *http.Request) {
	globalCounter.Lock()
	// Writing the response in our critical section generates lock contention.
	_, err := io.WriteString(w, fmt.Sprintln(globalCounter.count))
	if err != nil {
		log.Println("error writing response:", err)
	}
	globalCounter.Unlock()
}
