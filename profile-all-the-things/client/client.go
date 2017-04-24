package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	uri := flag.String("uri", "http://127.0.0.1:8080", "The URL to request.")
	flag.Parse()
	runClient(*uri)
}

func runClient(uri string) {
	for {
		resp, err := http.Get(uri)
		if err != nil {
			log.Println(err)
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
}
