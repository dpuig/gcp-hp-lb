package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

type Server struct {
	requestCounter *prometheus.CounterVec
}

func NewServer() *Server {
	return &Server{
		requestCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "Total number of API requests.",
			},
			[]string{"path", "method", "status"},
		),
	}
}

type target struct {
	url        *url.URL
	activeConn int32
}

func main() {
	// Set up viper to read environment variables.
	viper.AutomaticEnv()

	// Get the FUNCTION_TARGETS environment variable.
	functionTargetsEnv := viper.GetString("FUNCTION_TARGETS")

	// Split the environment variable into a slice.
	functionTargets := strings.Split(functionTargetsEnv, ",")

	server := NewServer()

	proxy := loadBalanceLeastConnections(functionTargets)

	// Instrumentation.
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(server.requestCounter)

	http.HandleFunc("/", server.trackMetrics(proxy.ServeHTTP))

	const readTimeoutSeconds = 5
	const writeTimeoutSeconds = 10

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      nil, // use default mux
		ReadTimeout:  readTimeoutSeconds * time.Second,
		WriteTimeout: writeTimeoutSeconds * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// Middleware for metrics.
func (s *Server) trackMetrics(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ... track start time, defer metrics update ...
		next.ServeHTTP(w, r)
	}
}

func loadBalanceLeastConnections(targets []string) *httputil.ReverseProxy {
	// Convert targets to *url.URL and wrap them in our custom struct.
	targetURLs := make([]*target, len(targets))
	for i, rawurl := range targets {
		u, _ := url.Parse(rawurl)
		targetURLs[i] = &target{url: u}
	}

	director := func(req *http.Request) {
		// Find the target with the least active connections.
		var minTarget *target
		for _, target := range targetURLs {
			if minTarget == nil || target.activeConn < minTarget.activeConn {
				minTarget = target
			}
		}

		// Increment active connections.
		atomic.AddInt32(&minTarget.activeConn, 1)
		defer atomic.AddInt32(&minTarget.activeConn, -1)

		// Rewrite the request to be sent to the selected target.
		req.URL.Scheme = minTarget.url.Scheme
		req.URL.Host = minTarget.url.Host
		req.URL.Path = singleJoiningSlash(minTarget.url.Path, req.URL.Path)
		if minTarget.url.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = minTarget.url.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = minTarget.url.RawQuery + "&" + req.URL.RawQuery
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
