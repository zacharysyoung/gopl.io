package testing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// test responds with a specific byte-sized payload, after a specific amount of time.
// It accepts the following query params:
// - sleep: a time.Duration formatted string (e.g., 1.5s, 750ms); defaults to "0s"
// - size: an integer; defaults to "1024" (bytes)
func test(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()

	sleepStr := q.Get("sleep")
	if sleepStr == "" {
		sleepStr = "0s"
	}
	sleep, err := time.ParseDuration(sleepStr)
	if err != nil {
		panic(err)
	}

	sizeStr := q.Get("size")
	if sizeStr == "" {
		sizeStr = "1024"
	}
	size, _ := strconv.Atoi(sizeStr)

	time.Sleep(sleep)

	fmt.Fprint(w, strings.Repeat("x", size))

	log.Printf("server.go - slept for %s and wrote %d bytes\n", sleep, size)
}

// StartListener starts and returns an *http.Server listening
// for calls to http://localhost:8889
func StartListener() *http.Server {
	srv := &http.Server{
		Addr: "localhost:8889",
	}
	http.HandleFunc("/test", test)
	go srv.ListenAndServe()
	log.Println("server.go - listening and serving up", srv.Addr)
	return srv
}

// StopListener calls Shutdown(context.Background()) on srv
func StopListener(srv *http.Server) {
	srv.Shutdown(context.Background())
	log.Println("server.go - shut down", srv.Addr)
}
