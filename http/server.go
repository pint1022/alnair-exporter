package http

import (
	"github.com/pint1022/alnair-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Server struct {
	Handler  http.Handler
	exporter exporter.Exporter
}

func NewServer(exporter exporter.Exporter) *Server {
	r := http.NewServeMux()

	// Register Metrics from each of the endpoints
	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(&exporter)

	r.Handle(exporter.MetricsPath(), promhttp.Handler())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
		                <head><title>Alnair Exporter</title></head>
		                <body>
		                   <h1>GitHub Prometheus Metrics Exporter</h1>
						   <p>For more information, visit <a href=https://github.com/pint1022/alnair-exporter>GitHub</a></p>
		                   <p><a href='` + exporter.MetricsPath() + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})

	return &Server{Handler: r, exporter: exporter}
}

func (s *Server) Start() {
	log.Fatal(http.ListenAndServe(":"+s.exporter.ListenPort(), s.Handler))
}
