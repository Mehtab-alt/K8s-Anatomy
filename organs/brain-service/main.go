// File: organs/brain-service/main.go
package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// httpRequestDuration measures the latency of "thoughts".
// This is a critical health indicator for our organism's cognitive function.
var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Histogram of request latencies, representing cognitive processing time.",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 0.75, 1.0},
	},
	[]string{"service"}, // The 'service' label identifies this organ.
)

func init() {
	prometheus.MustRegister(httpRequestDuration)
}

// thinkHandler simulates a complex cognitive task.
func thinkHandler(w http.ResponseWriter, r *http.Request) {
	// Start a timer to measure the duration of the thought process.
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues("brain-service"))
	defer timer.ObserveDuration()

	// Simulate work by introducing a random delay.
	minDelay := 50 * time.Millisecond
	maxDelay := 750 * time.Millisecond
	delay := minDelay + time.Duration(rand.Int63n(int64(maxDelay-minDelay)))
	time.Sleep(delay)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("thought processed"))
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	http.HandleFunc("/think", thinkHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Brain-Service is online and thinking on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Brain-Service failed to start: %v", err)
	}
}
