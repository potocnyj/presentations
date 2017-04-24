// demo4 shows an example of using a custom profile
// to track in-progress requests in a simple http service.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
	"sync"
	"time"
)

const profileName = "activeRequests"

var myProfile *pprof.Profile

func init() {
	// Check that the profile doesn't exist;
	// if we try to create it twice we get a panic!
	if myProfile = pprof.Lookup(profileName); myProfile == nil {
		myProfile = pprof.NewProfile(profileName)
	}
}

func main() {
	http.HandleFunc("/w", handleIncCounter)
	http.HandleFunc("/r", handleReadCounter)
	http.ListenAndServe("localhost:8080", nil)
}

// trackRequest wraps an http handler with tracing information
// so we can profile in-flight requests.
func trackRequest() func() {
	// Allocate a byte and use its memory address as a key for the profile.
	key := new(byte)

	// Add the key to track that we have another request in progress.
	// Pass skip=1 to omit the myProfile.Add call in the stack recorded;
	// higher values will remove more calls and start higher up the stack.
	myProfile.Add(key, 1)

	// Introduce artificial delay in the request we're tracking,
	// since they don't do much work.
	time.Sleep(time.Millisecond * 500)

	return func() {
		// Remove the counter for the now completed request.
		myProfile.Remove(key)
	}
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

func handleIncCounter(w http.ResponseWriter, r *http.Request) {
	defer trackRequest()()
	incrementCounter()
}

func handleReadCounter(w http.ResponseWriter, r *http.Request) {
	defer trackRequest()()
	globalCounter.Lock()
	_, err := io.WriteString(w, fmt.Sprintln(globalCounter.count))
	if err != nil {
		log.Println("error writing response:", err)
	}
	globalCounter.Unlock()
}
