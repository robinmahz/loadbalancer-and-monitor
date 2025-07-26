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
}

type LoadBalancer struct {
	backends []*Backend
	current  int
	mu       sync.Mutex
}

func NewLoadBalancer(backendURLs []string) *LoadBalancer {
	backends := make([]*Backend, len(backendURLs))
	for i, backendURL := range backendURLs {
		url, _ := url.Parse(backendURL)
		backends[i] = &Backend{
			URL:          url,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		}
	}
	return &LoadBalancer{backends: backends}
}

func (lb *LoadBalancer) getNextBackend() *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := 0; i < len(lb.backends); i++ {
		lb.current = (lb.current + 1) % len(lb.backends)
		if lb.backends[lb.current].Alive {
			return lb.backends[lb.current]
		}
	}
	return nil
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

	lb := NewLoadBalancer(backendURLs)
	go lb.healthCheck()

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", lb)

	log.Println("Starting load balancer on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
