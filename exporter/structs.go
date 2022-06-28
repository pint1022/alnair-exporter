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
	"net/http"
	"sync"

	"github.com/pint1022/alnair-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter is used to store Metrics data and embeds the config struct.
// This is done so that the relevant functions have easy access to the
// user defined runtime configuration when the Collect method is called.
type Exporter struct {
	APIMetrics map[string]*prometheus.Desc
	config.Config
	mu sync.Mutex

}

type Asset struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Downloads int32  `json:"download_count"`
	CreatedAt string `json:"created_at"`
}

// RateLimits is used to store rate limit data into a struct
// This data is later represented as a metric, captured at the end of a scrape
type RateLimits struct {
	Limit     float64
	Remaining float64
	Reset     float64
}

// Response struct is used to store http.Response and associated data
type Response struct {
	url      string
	response *http.Response
	body     []byte
	err      error
}

// RateLimits is used to store rate limit data into a struct
// This data is later represented as a metric, captured at the end of a scrape
type GPUMetrics struct {
	Ts      int64   // time stamp
	Bs    	int32   // Burst size
	Ou   	int32   // Over use
	Ws 		int32   // Window size
	Hd    	int32   // Host2Device duration 
	Dh    	int32   // Device2Host duration
	Rm      int32   // Remain quota
	Um      int32   // Used memory
	Mm      int32   // Memory limit
}
type comm_request_t int32

const (
   REQ_QUOTA comm_request_t = 0
   REQ_MEM_LIMIT = 1
   REQ_MEM_UPDATE = 2
   REQ_SAMPLE = 3
   SAM_MSG_LEN = 128
   POD_NAME_LEN = 4
   UUID_LEN = 4
   END_READ = "NONE"
) 

type UnpackedSample struct {
    id  uint32
	length int32
    sample []byte
}


type Packet struct {
    Len     uint64
    Msg     [NAME_LEN]byte
    ReqId   int32
    ReqType comm_request_t
}
