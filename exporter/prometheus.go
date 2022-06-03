/**
 * Copyright 2022 Steven Wang, Futurewei Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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

	var GPUdata *GPUMetrics
	var rc int
	var podname []byte
	var uuid []byte

	// for {
		rc, GPUdata, podname, uuid = e.getGPUMetrics()
		if rc <= 0 {
			log.Errorf("Error gathering GPU metrics from remote API: ", rc)
			return
		}

		// Set prometheus gauge metrics using the data gathered
		err = e.processGPUMetrics(GPUdata, ch, string(podname), string(uuid))
	
		if err != nil {
			log.Error("Error Processing Metrics", err)
			return
		}
	// 	if string(uuid) == END_READ	{
	// 		break
	// 	}
	// }
	log.Info("All Metrics successfully collected")

}
