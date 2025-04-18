package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"sync/atomic"
	"time"

	
	"github.com/prometheus/client_golang/prometheus/promhttp"
)



func main() {
	// Create trace file
	traceFile, err := os.Create("goroutine_trace.out")
	if err != nil {
		log.Fatalf("failed to create trace file: %v", err)
	}
	defer traceFile.Close()

	// Start tracing
	if err := trace.Start(traceFile); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	// Start pprof server
	go func() {
		log.Println("Starting pprof server on :6063")
		log.Println(http.ListenAndServe("localhost:6063", nil))
	}()

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting metrics server on :2144")
		log.Println(http.ListenAndServe(":2144", nil))
	}()

	var goroutineCount int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Total operation duration
	const operationDuration = 1 * time.Minute
	endTime := time.Now().Add(operationDuration)

	// Goroutine manager
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		phase := "increasing" // "increasing" or "decreasing"
		target := 0

		for range ticker.C {
			if time.Now().After(endTime) {
				cancel()
				return
			}

			elapsed := time.Since(endTime.Add(-operationDuration))
			progress := float64(elapsed) / float64(operationDuration)

			switch phase {
			case "increasing":
				target = int(progress * 100)
				if target >= 100 {
					target = 100
					phase = "decreasing"

				}
			case "decreasing":
				target = 100 - int(progress*100)
				if target <= 0 {
					target = 0
				}
			}

			current := int(atomic.LoadInt32(&goroutineCount))
			diff := target - current


			if diff > 0 {
				// Need to add goroutines
				for i := 0; i < diff; i++ {
					go worker(ctx, &goroutineCount)
				}
			}

			// Update current goroutine gauge

			fmt.Printf("Target: %d, Current: %d, Phase: %s\n", target, current, phase)
		}
	}()

	// Metrics reporter
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			count := atomic.LoadInt32(&goroutineCount)
			log.Printf("Active goroutines: %d\n", count)
		}
	}()



	// Wait for operation to complete
	<-ctx.Done()
	log.Println("Operation completed")
}

func worker(ctx context.Context, counter *int32) {
	atomic.AddInt32(counter, 1)
	defer atomic.AddInt32(counter, -1)

	// Simulate work
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Do some work
			_ = 1 + 1
		}
	}
}
