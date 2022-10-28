package testing

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

// TestListener tests that server honors the sleep and size params; runs
// multiple GETs concurrently for speedier testing.
func TestListener(t *testing.T) {
	type Test struct {
		sleep time.Duration
		sizeB int
	}

	testEnpoint := func(test Test, srv *http.Server, done chan<- bool) {
		// Use %s formatter for sleep (time.Duration)
		url := fmt.Sprintf("http://%s/test?sleep=%s&size=%d", srv.Addr, test.sleep, test.sizeB)

		beg := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("Request of \"%s\" error'd out: %v", url, err)
		}
		dur := time.Since(beg)

		if dur < test.sleep {
			t.Errorf("Expected test endpoint to sleep at least %d secs, only slept for %v\n", test.sleep, dur)
		}
		if dur > test.sleep+500*time.Millisecond {
			t.Errorf("Expected test endpoint to sleep %d secs plus another 0.5 secs, actually slept for %v\n", test.sleep, dur)
		}

		defer resp.Body.Close()
		nbytes, err := io.Copy(io.Discard, resp.Body)
		if err != nil {
			t.Errorf("Could not read body: %q\n", err)
		}
		if nbytes != int64(test.sizeB) {
			t.Errorf("Expected response-body size of %d bytes, instead got %d bytes\n", test.sizeB, nbytes)
		}

		done <- true
	}

	var tests = []Test{
		{100 * time.Millisecond, 0},
		{250 * time.Millisecond, 1024},
		{500 * time.Millisecond, 1024 * 1024},
	}

	done := make(chan bool, len(tests))

	srv := StartListener()

	for _, test := range tests {
		if testing.Verbose() {
			log.Printf("testing, sleep:%s, size:%d", test.sleep, test.sizeB)
		}
		go testEnpoint(test, srv, done)
	}
	for len(done) < len(tests) {
	}

	StopListener(srv)
}

func TestEmptyParams(t *testing.T) {

}
