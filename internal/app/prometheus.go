package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheusExporter() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PROMETHEUS_PORT")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	go func() {
		log.Info("Starting Prometheus exporter", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Prometheus exporter failed", "err", err)
		}
	}()
}
