# Enhanced Goroutine Scaling App with Prometheus Metrics


![image](src/goMonitor.png)

Here's the improved version that adds Prometheus metrics to monitor:
- Current goroutine count
- Goroutines created per second
- Phase transitions (increasing/decreasing)

```go
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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Prometheus metrics
	goroutineGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_goroutines_current",
		Help: "Current number of active goroutines",
	})

	goroutinesCreatedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_goroutines_created_total",
		Help: "Total number of goroutines created",
	})

	phaseGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "app_phase",
		Help: "Current phase of the scaling operation",
	}, []string{"phase"})

	scalingRateGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_scaling_rate",
		Help: "Current goroutines scaling rate (positive=increasing, negative=decreasing)",
	})
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
		lastTarget := 0
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
					phaseGauge.WithLabelValues("increasing").Set(0)
					phaseGauge.WithLabelValues("decreasing").Set(1)
				}
			case "decreasing":
				target = 100 - int(progress*100)
				if target <= 0 {
					target = 0
				}
			}

			current := int(atomic.LoadInt32(&goroutineCount))
			diff := target - current

			// Update scaling rate metric
			scalingRateGauge.Set(float64(target - lastTarget))
			lastTarget = target

			if diff > 0 {
				// Need to add goroutines
				for i := 0; i < diff; i++ {
					go worker(ctx, &goroutineCount)
					goroutinesCreatedCounter.Inc()
				}
			}

			// Update current goroutine gauge
			goroutineGauge.Set(float64(current))

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

	// Initialize phase metrics
	phaseGauge.WithLabelValues("increasing").Set(1)
	phaseGauge.WithLabelValues("decreasing").Set(0)

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
```

## Key Metrics Added

1. **Current Goroutines**:
   - Metric: `app_goroutines_current` (Gauge)
   - Tracks the real-time number of active goroutines

2. **Total Goroutines Created**:
   - Metric: `app_goroutines_created_total` (Counter)
   - Counts all goroutines created during the operation

3. **Phase Tracking**:
   - Metric: `app_phase` (GaugeVec)
   - Labels: `phase="increasing"` or `phase="decreasing"`
   - Value 1 indicates active phase

4. **Scaling Rate**:
   - Metric: `app_scaling_rate` (Gauge)
   - Positive values = increasing goroutines
   - Negative values = decreasing goroutines

## How to Use the Metrics


0. before run once:
```sh
go mod init mygometric
go mod tidy
```
   
1. Run the application:
   ```bash
   go run main.go
   ```

2. Access metrics:
   - Prometheus metrics: http://localhost:2144/metrics
   - pprof: http://localhost:6063/debug/pprof/
   - `go tool pprof -http=:8022 http://localhost:6063/debug/pprof/profile\?seconds\=30`
   - Trace: `go tool trace goroutine_trace.out`

3. Example Prometheus queries:
   ```promql
   # Current goroutine count over time
   app_goroutines_current

   # Goroutine creation rate
   rate(app_goroutines_created_total[1m])

   # Time spent in each phase
   sum by (phase) (rate(app_phase[1m]))
   ```

## Grafana Dashboard Suggestions

Create a dashboard with these panels:
1. **Current Goroutines**: Line graph of `app_goroutines_current`
2. **Phase Timeline**: Stacked graph of `app_phase`
3. **Scaling Rate**: Line graph of `app_scaling_rate`
4. **Creation Rate**: Counter graph of `app_goroutines_created_total`

This enhanced version gives you complete observability into the goroutine scaling behavior with both immediate metrics (Prometheus) and deep profiling (pprof+trace).