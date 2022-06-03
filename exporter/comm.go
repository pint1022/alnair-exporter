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
	// "bufio"
	"time"
	"context"
	"net"
	"bytes"
	// "fmt"
    "encoding/binary"
    "encoding/json"
	log "github.com/sirupsen/logrus"
)

const NET_OP_MAX_ATTEMPT = 5  // maximum time retrying failed network operations
const NET_OP_RETRY_INTV = 10  // seconds between two retries
const NAME_LEN = 21  // seconds between two retries


func multiple_attempt(f func(net.Conn, []byte) (int, []byte), conn net.Conn, req []byte, max_attempt int, interval int) (int,  []byte) {
	var rc int
	var resp []byte

	for attempt := 1; attempt <= max_attempt; attempt++ {
	  rc, resp = f(conn, req);
	  if (rc > 0) {
		  break
	  }
	  log.Errorf("attempt %d", attempt)
	  if (interval > 0) {
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
  }
  return rc, []byte(resp)

}

func (e *Exporter) communicate(CONNECT string, reqType comm_request_t ) (int, []byte, []byte, []byte) {
	req := e.prepare_request(reqType)

	e.mu.Lock()
	defer e.mu.Unlock()
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", CONNECT)
	if err != nil {
		log.Fatalf("Failed to dial: %v, %s", err, CONNECT)
	}
	defer conn.Close()

	// perform communication
	rc, resp:= multiple_attempt(
		func ( c net.Conn, r []byte) (int, []byte) {
			_, err = c.Write(r)
			if err != nil {
				return  0, []byte("failed to send the request.")
			}

			resp := make([]byte, SAM_MSG_LEN)

			// var n int32
		
			n, err1 := conn.Read(resp)
			if err1 != nil {
				log.Info("Failed to read response from alnair server")
				return 0, []byte("failed to read the response.")
			}

		  return n, resp;
		}, conn, req,
		NET_OP_MAX_ATTEMPT, NET_OP_RETRY_INTV)

	// log.Info("Resp:", resp, ", length: ", rc)

	if rc <= 0 {
		log.Fatalf("Failed to communicate with alnair server daemon")
		return rc,  []byte(""), []byte(""), []byte("")
	}
	length, sample, podname, uuid := e.parse_sample(resp)

	return int(length), sample, podname, uuid
  }

  func (e *Exporter) parse_response(resp []byte) []byte {
    var unpacked UnpackedSample

	buf := bytes.NewReader(resp)

    err := binary.Read(buf, binary.LittleEndian, &unpacked)
	if err != nil {
		log.Errorf("Unable to parse the response, Error: %s", err)
		return []byte("")
	}    
	return unpacked.sample
  }

  func (e *Exporter) parse_sample(resp []byte) (int32, []byte, []byte, []byte) {
    var unpacked UnpackedSample

	buf := bytes.NewReader(resp)

    err := binary.Read(buf, binary.LittleEndian, &unpacked.id)
	if err != nil {
		log.Errorf("Unable to parse id, Error: ", err)
		return 0, []byte("reqid is wrong"), []byte(""), []byte("")
	}
	var _len int32
    err = binary.Read(buf, binary.LittleEndian, &_len)
	if err != nil {
		log.Errorf("Unable to parse pod name length, Error: ", err)
		return 0, []byte("response length is wrong"), []byte(""), []byte("")
	}    
    podname := make([]byte, _len)
    err = binary.Read(buf, binary.LittleEndian, &podname)
	// log.Info("podname: ", string(podname), ", len: ", _len)

	if err != nil {
		log.Errorf("Unable to parse podname, Error: ", err)
        return 0, []byte("podname are wrong."), []byte(""), []byte("")
	}
    err = binary.Read(buf, binary.LittleEndian, &_len)
	if err != nil {
		log.Errorf("Unable to parse uuid length, Error: ", err)
		return 0, []byte("response uuid length is wrong"), []byte(""), []byte("")
	}    
    uuid := make([]byte, _len)
    err = binary.Read(buf, binary.LittleEndian, &uuid)
	// log.Info("uuid: ", string(uuid), ", len: ", _len)

	if err != nil {
		log.Errorf("Unable to parse uuid, Error: ", err)
        return 0, []byte("uuid are wrong."), []byte(""), []byte("")
	}

    err = binary.Read(buf, binary.LittleEndian, &unpacked.length)
	if err != nil {
		log.Errorf("Unable to parse sample length, Error: ", err)
		return 0, []byte("response length is wrong"), []byte(""), []byte("")
	}    
    unpacked.sample = make([]byte, unpacked.length)

    err = binary.Read(buf, binary.LittleEndian, &unpacked.sample)

	if err != nil {
		log.Errorf("Unable to parse sample, Error: ", err)
        return 0, []byte("sample contents are wrong."), []byte(""), []byte("")
	}

	return unpacked.length,  unpacked.sample, podname, uuid
  }

// Attempt a function several times. Non-zero return of func is treated as an error. If func return
// -1, errno will be returned.

func (e *Exporter) prepare_request(r comm_request_t) []byte {

	var id int32

	pod_name := "alnair-client\000"
	id = 0
	var a [NAME_LEN]byte

	copy(a[:], pod_name)	
	req := Packet{Len:uint64(len(a) - 1), Msg:a, ReqId:id, ReqType:r}
    _, err := json.Marshal(req)
    if err != nil {
        panic (err)
    }
	// log.Info("Req: ", string(out))
  
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, req)

	if err != nil {
		log.Errorf("Unable to create the request, Error: %s", err)
		return []byte("")
	}
    // fmt.Println(buf.String())

	return buf.Bytes()
  }