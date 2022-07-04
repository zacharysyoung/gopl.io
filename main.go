package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	gopltest "gopl.io/testing"
)

func main() {
	srv := gopltest.StartListener()

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// if err := srv.ListenAndServe(); err != http.ErrServerClosed {
	// 	// Error starting or closing listener:
	// 	log.Fatalf("HTTP server ListenAndServe: %v", err)
	// }

	<-idleConnsClosed
}
