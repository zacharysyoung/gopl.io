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

func test(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()

	timeoutStr := q.Get("timeout")
	if timeoutStr == "" {
		timeoutStr = "0s"
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		panic(err)
	}

	sizeStr := q.Get("size")
	if sizeStr == "" {
		sizeStr = "1024"
	}
	size, _ := strconv.Atoi(sizeStr)

	time.Sleep(timeout)

	fmt.Fprint(w, strings.Repeat("x", size))
}

// StartListener starts and returns an *http.Server, listening for calls to
// http://localhost:8889/test?timeout=XX&size=YY
//
// -timeout: a time.Duration formatted string, e.g.: 1.5s, 750ms
// -size: an int byte count
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
