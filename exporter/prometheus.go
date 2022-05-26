package exporter

import (
	"github.com/prometheus/client_golang/prometheus"

    // "encoding/json"
	log "github.com/sirupsen/logrus"
)

// Describe - loops through the API metrics and passes them to prometheus.Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.APIMetrics {
		ch <- m
	}

}

// Collect function, called on by Prometheus Client library
// This function is called when a scrape is peformed on the /metrics page
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	// data := []*Datum{}
	var err error

	// var GPUdata *GPUMetrics
	// var rc int
	// var podname []byte
	// var uuid string

	rc, GPUdata, podname, uuid := e.getGPUMetrics()
	if rc <= 0 {
		log.Errorf("Error gathering GPU metrics from remote API: ", rc)
		return
	}
    // out, err := json.Marshal(GPUdata)
    // if err != nil {
    //     panic (err)
    // }

    // fmt.Println("struct is: ", *GPUdata)

	// Set prometheus gauge metrics using the data gathered
	err = e.processGPUMetrics(GPUdata, ch, string(podname), string(uuid))

	if err != nil {
		log.Error("Error Processing Metrics", err)
		return
	}
	log.Info("All Metrics successfully collected")

}
