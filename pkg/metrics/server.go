package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Start runs a lightweight HTTP server exposing /metrics for Prometheus scraping.
// It runs in a background goroutine and does not block.
func Start(port string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Printf("Metrics server listening on :%s/metrics", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Printf("metrics server error: %v", err)
		}
	}()
}
