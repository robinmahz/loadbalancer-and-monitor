package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "lb_requests_total",
            Help: "Total number of requests handled by load balancer.",
        },
        []string{"backend"},
    )
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "lb_request_duration_seconds",
            Help:    "Duration of requests handled by load balancer.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"backend"},
    )
)

type Backend struct {
    URL          *url.URL
    Alive        bool
    ReverseProxy *httputil.ReverseProxy
    Weight       int // Added weight field
    CurrentWeight int // Added for weighted round-robin algorithm
}

type LoadBalancer struct {
    backends []*Backend
    mu       sync.Mutex
}

func NewLoadBalancer(backendURLs []string, weights []int) *LoadBalancer {
    if len(backendURLs) != len(weights) {
        log.Fatal("Number of backends and weights must match")
    }
    backends := make([]*Backend, len(backendURLs))
    for i, backendURL := range backendURLs {
        url, _ := url.Parse(backendURL)
        backends[i] = &Backend{
            URL:          url,
            Alive:        true,
            ReverseProxy: httputil.NewSingleHostReverseProxy(url),
            Weight:       weights[i],
            CurrentWeight: weights[i], // Initialize current weight to max weight
        }
    }
    return &LoadBalancer{backends: backends}
}

func (lb *LoadBalancer) getNextBackend() *Backend {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    // Find total weight and maximum current weight
    totalWeight := 0
    maxCurrentWeight := -1
    var selectedBackend *Backend

    for _, backend := range lb.backends {
        if backend.Alive {
            totalWeight += backend.Weight
            backend.CurrentWeight += backend.Weight
            if backend.CurrentWeight > maxCurrentWeight {
                maxCurrentWeight = backend.CurrentWeight
                selectedBackend = backend
            }
        }
    }

    if selectedBackend == nil {
        return nil // No healthy backends available
    }

    // Decrease the current weight of the selected backend
    selectedBackend.CurrentWeight -= totalWeight

    return selectedBackend
}

func (lb *LoadBalancer) healthCheck() {
    for {
        for _, backend := range lb.backends {
            resp, err := http.Get(backend.URL.String() + "/")
            backend.Alive = err == nil && resp.StatusCode == http.StatusOK
            if !backend.Alive {
                log.Printf("Backend %s is down", backend.URL.String())
            }
        }
        time.Sleep(10 * time.Second)
    }
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    backend := lb.getNextBackend()
    if backend == nil {
        http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
        return
    }

    start := time.Now()
    requestCounter.WithLabelValues(backend.URL.String()).Inc()
    backend.ReverseProxy.ServeHTTP(w, r)
    requestDuration.WithLabelValues(backend.URL.String()).Observe(time.Since(start).Seconds())
}

func init() {
    prometheus.MustRegister(requestCounter, requestDuration)
}

func main() {
    backendURLs := []string{
        "http://backend1:8081",
        "http://backend2:8082",
        "http://backend3:8083",
    }
    weights := []int{1, 2, 3} // Example weights: backend1 gets 1/6, backend2 gets 2/6, backend3 gets 3/6 of requests

    lb := NewLoadBalancer(backendURLs, weights)
    go lb.healthCheck()

    http.Handle("/metrics", promhttp.Handler())
    http.Handle("/", lb)

    log.Println("Starting load balancer on port 8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}