package exporter

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)


// isArray simply looks for key details that determine if the JSON response is an array or not.
func isArray(body []byte) bool {

	isArray := false

	for _, c := range body {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		isArray = c == '['
		break
	}

	return isArray

}

  
func (e *Exporter) getGPUMetrics()(int, *GPUMetrics, []byte, []byte)  {

	CONNECT := e.AlnrIP() + ":" + e.AlnrPort()
	rc, sample, podname, uuid:= e.communicate(CONNECT, REQ_SAMPLE)

	if rc <= 0 {
		err := "No sampling data at this moment ..."
		log.Errorf(err)
		return rc, &GPUMetrics{}, []byte(""), []byte("")
	}
	// println("Sample: ", string(sample), rc)

	var data GPUMetrics
	json.Unmarshal(sample, &data)
    // fmt.Println("Struct is:", data)

	return rc, &data, podname, uuid
}