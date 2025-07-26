package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests.",
        },
        []string{"endpoint", "server"},
    )
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"endpoint", "server"},
    )
)

func init() {
    prometheus.MustRegister(requestCounter, requestDuration)
}

func main() {
    serverName := os.Getenv("SERVER_NAME")
    if serverName == "" {
        serverName = "unknown"
    }
    capacity := os.Getenv("CAPACITY")
    if capacity == "" {
        capacity = "medium"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // Simulate some CPU work
        for i := 0; i < 1000000; i++ {
            _ = i * i // Simple computation to consume CPU
        }
        requestCounter.WithLabelValues("/", serverName).Inc()
        requestDuration.WithLabelValues("/", serverName).Observe(time.Since(start).Seconds())
        fmt.Fprintf(w, "Hello from %s (Capacity: %s)", serverName, capacity)
    })

    http.Handle("/metrics", promhttp.Handler())

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server %s with capacity %s on port %s", serverName, capacity, port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}