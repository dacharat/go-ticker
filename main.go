package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Create a context that will be used for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a ticker that ticks every 1 seconds
	ticker := time.NewTicker(time.Second)

	// Create waitgroup to wait goroutine
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case t := <-ticker.C:
				time.Sleep(300 * time.Millisecond)
				// Perform your action here, e.g., printing the current time
				fmt.Println("Current time:", t.Format(time.RFC3339Nano))
			case <-ctx.Done():
				// The context was canceled, indicating a graceful shutdown request
				fmt.Println("Ticker stopped gracefully")
				return
			}

		}
	}()

	// Set up a signal handler to catch termination signals
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	// Wait for a termination signal
	<-gracefulStop

	// Signal the Goroutine to stop by canceling the context
	cancel()

	// Stop the ticker and wait for it to finish
	ticker.Stop()

	// Wait for waitgroup
	wg.Wait()

	// You can perform additional cleanup or shutdown tasks here
	fmt.Println("Ticker stopped successfully.")
}
