// File: organs/heart-service/main.go
package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// httpRequestsTotal is the Prometheus Counter for tracking heartbeats.
// This is a vital sign, representing the pulse of our organism.
var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed, representing the organism's pulse.",
	},
	[]string{"service"}, // The 'service' label identifies which organ is emitting the signal.
)

// init registers our Prometheus metric with the default registry.
func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

// beatHandler handles requests to the /beat endpoint.
// Each request signifies a single "lub-dub" or heartbeat.
func beatHandler(w http.ResponseWriter, r *http.Request) {
	// Increment the counter for this specific service.
	httpRequestsTotal.With(prometheus.Labels{"service": "heart-service"}).Inc()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("lub-dub"))
}

func main() {
	// The '/beat' endpoint is the primary function of this organ.
	http.HandleFunc("/beat", beatHandler)

	// The '/metrics' endpoint is the nerve ending, exposing vital signs to Prometheus.
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Heart-Service is alive and beating on port 8080...")
	// The organ begins its life function.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Heart-Service failed to start: %v", err)
	}
}
