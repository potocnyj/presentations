// demo2 shows an example of enabling block profiling
// to identify contention in a simple http service.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

func main() {
	// Sample blocking events once every 100ns.
	runtime.SetBlockProfileRate(100)

	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
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
	// Introduce a 1us block.
	<-time.After(time.Microsecond)
	_, err := io.WriteString(w, fmt.Sprint(incrementCounter()))
	if err != nil {
		log.Println("error writing response:", err)
	}
}
