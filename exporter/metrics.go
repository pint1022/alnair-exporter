package exporter

import "github.com/prometheus/client_golang/prometheus"

// AddMetrics - Add's all of the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {

	APIMetrics := make(map[string]*prometheus.Desc)



	APIMetrics["BurstSize"] = prometheus.NewDesc(
		prometheus.BuildFQName("GPU", "pod", "burst"),
		"The elapse time (milliseconds)at which the current kernel runs on GPU.",
		[]string{"pod", "gpu_uuid"}, nil,
	)	
	APIMetrics["Overuse"] = prometheus.NewDesc(
		prometheus.BuildFQName("GPU", "pod", "overuse"),
		"The elapse time (milliseconds)at which the current kernel runs overtime on GPU.",
		[]string{"pod", "gpu_uuid"}, nil,
	)	
	APIMetrics["MemH2D"] = prometheus.NewDesc(
		prometheus.BuildFQName("GPU", "pod", "H2D"),
		"The elapse time (milliseconds)at which the current pod memory copy (H2D) on GPU.",
		[]string{"pod", "gpu_uuid"}, nil,
	)	
	APIMetrics["MemD2H"] = prometheus.NewDesc(
		prometheus.BuildFQName("GPU", "pod", "D2H"),
		"The elapse time (milliseconds)at which the current pod memory copy(D2H) on GPU.",
		[]string{"pod", "gpu_uuid"}, nil,
	)	
	return APIMetrics
}


// processGPUmetrics - processes the response GPU metrics using it as a source
func (e *Exporter) processGPUMetrics( data *GPUMetrics, ch chan<- prometheus.Metric, podname string, uuid string) error {

	// Set Rate limit stats
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["BurstSize"], prometheus.GaugeValue, float64(data.Bs), podname, uuid)
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["Overuse"], prometheus.GaugeValue, float64(data.Ou), podname, uuid)
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["MemH2D"], prometheus.GaugeValue, float64(data.Hd), podname, uuid)
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["MemD2H"], prometheus.GaugeValue, float64(data.Dh), podname, uuid)

	return nil
}
