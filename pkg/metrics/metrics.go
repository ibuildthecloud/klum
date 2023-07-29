package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	ErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "klum_errors_total",
		Help: "The total number of errors found",
	})
)

func StartMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics erver on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
