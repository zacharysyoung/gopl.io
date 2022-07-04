package main

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	gopltest "gopl.io/testing"
)

func buildFetchAll() {
	buildCmd := exec.Command("go", "build", "-o", "fetchall", "main.go")
	_, err := buildCmd.Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit rc =", e.ExitCode())
		default:
			panic(err)
		}
	}
}

func TestConcurrentFetchAll(t *testing.T) {
	log.Println("main_test.go - starting build")
	buildFetchAll()
	log.Println("main_test.go - finished build")

	// Elapsed should be around 2s, no more than the longest wait
	fetchAllCmd := exec.Command(
		"./fetchall",
		"http://localhost:8889/test?timeout=2s&size=50",
		"http://localhost:8889/test?timeout=1.5s&size=2550",
	)

	srv := gopltest.StartListener()
	defer gopltest.StopListener(srv)

	beg := time.Now()
	log.Println("main_test.go - starting exec")
	fetchAllOutput, err := fetchAllCmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	elapsed := time.Since(beg)

	if testing.Verbose() {
		fmt.Print(string(fetchAllOutput))
	}

	min := 2 * time.Second
	max := min + 150*time.Millisecond // allow for some latency in the processes foo
	if elapsed < min || elapsed > max {
		t.Errorf("expected all requests to take no less than %s and no more than %s, it actually took %s", min, max, elapsed)
	}
}
