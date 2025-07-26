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
    activeConnections = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "lb_active_connections",
            Help: "Number of active connections per backend.",
        },
        []string{"backend"},
    )
)

type Backend struct {
    URL              *url.URL
    Alive            bool
    ReverseProxy     *httputil.ReverseProxy
    ActiveConnections int // Tracks current number of active connections
    mu               sync.Mutex  // Protects ActiveConnections
}

type LoadBalancer struct {
    backends []*Backend
    mu       sync.Mutex
}

func NewLoadBalancer(backendURLs []string) *LoadBalancer {
    backends := make([]*Backend, len(backendURLs))
    for i, backendURL := range backendURLs {
        url, _ := url.Parse(backendURL)
        backends[i] = &Backend{
            URL:              url,
            Alive:            true,
            ReverseProxy:     httputil.NewSingleHostReverseProxy(url),
            ActiveConnections: 0,
        }
    }
    return &LoadBalancer{backends: backends}
}

func (lb *LoadBalancer) getNextBackend() *Backend {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    minConnections := -1
    var selectedBackend *Backend

    for _, backend := range lb.backends {
        if backend.Alive {
            backend.mu.Lock()
            conn := backend.ActiveConnections
            backend.mu.Unlock()
            if minConnections == -1 || conn < minConnections {
                minConnections = conn
                selectedBackend = backend
            }
        }
    }

    if selectedBackend == nil {
        return nil // No healthy backends available
    }

    return selectedBackend
}

func (lb *LoadBalancer) healthCheck() {
    for {
        for _, backend := range lb.backends {
            resp, err := http.Get(backend.URL.String() + "/")
            backend.mu.Lock()
            backend.Alive = err == nil && resp != nil && resp.StatusCode == http.StatusOK
            if !backend.Alive {
                log.Printf("Backend %s is down", backend.URL.String())
            }
            backend.mu.Unlock()
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

    // Increment active connections
    backend.mu.Lock()
    backend.ActiveConnections++
    activeConnections.WithLabelValues(backend.URL.String()).Inc()
    backend.mu.Unlock()

    start := time.Now()
    requestCounter.WithLabelValues(backend.URL.String()).Inc()

    // Wrap the response writer to detect when the request is complete
    responseWriter := &responseWriterWrapper{
        ResponseWriter: w,
        onClose: func() {
            backend.mu.Lock()
            backend.ActiveConnections--
            activeConnections.WithLabelValues(backend.URL.String()).Dec()
            backend.mu.Unlock()
        },
    }

    backend.ReverseProxy.ServeHTTP(responseWriter, r)
    requestDuration.WithLabelValues(backend.URL.String()).Observe(time.Since(start).Seconds())
}

// responseWriterWrapper wraps http.ResponseWriter to detect when the response is complete
type responseWriterWrapper struct {
    http.ResponseWriter
    onClose func()
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
    defer w.onClose()
    return w.ResponseWriter.Write(data)
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
    defer w.onClose()
    w.ResponseWriter.WriteHeader(statusCode)
}

func init() {
    prometheus.MustRegister(requestCounter, requestDuration, activeConnections)
}

func main() {
    backendURLs := []string{
        "http://backend1:8081",
        "http://backend2:8082",
        "http://backend3:8083",
    }

    lb := NewLoadBalancer(backendURLs)
    go lb.healthCheck()

    http.Handle("/metrics", promhttp.Handler())
    http.Handle("/", lb)

    log.Println("Starting load balancer on port 8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}