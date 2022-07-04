package main

import (
	"fmt"
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
	buildFetchAll()

	// Elapsed should be around 2s, no more than the longest wait
	fetchAllCmd := exec.Command(
		"./fetchall",
		"http://localhost:8889/test?timeout=2s&size=50",
		"http://localhost:8889/test?timeout=1.5s&size=2550",
	)

	srv := gopltest.StartListener()
	defer gopltest.StopListener(srv)

	beg := time.Now()
	_, err := fetchAllCmd.Output()
	if err != nil {
		t.Error(err)
	}
	elapsed := time.Since(beg)

	min := 2 * time.Second
	max := min + 150*time.Millisecond
	if elapsed < min || elapsed > max {
		t.Errorf("expected all requests to take no less than %s and no more than %s, it actually took %s", min, max, elapsed)
	}
}
