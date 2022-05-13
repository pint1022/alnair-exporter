package exporter

import (
	"bufio"
	"time"
	"context"
	"net"
	"bytes"
	"fmt"
    "encoding/binary"
    "encoding/json"
	log "github.com/sirupsen/logrus"
)

const NET_OP_MAX_ATTEMPT = 5  // maximum time retrying failed network operations
const NET_OP_RETRY_INTV = 10  // seconds between two retries
const NAME_LEN = 21  // seconds between two retries


func multiple_attempt(f func(net.Conn, []byte) (string, int), conn net.Conn, req []byte, max_attempt int, interval int) ([]byte, int) {
	var rc int
	var resp string

	for attempt := 1; attempt <= max_attempt; attempt++ {
	  resp, rc = f(conn, req);
	  if (rc == 0) {
		  break
	  }
	  log.Errorf("attempt %d", attempt)
	  if (interval > 0) {
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
  }
  return []byte(resp), rc

}

func (e *Exporter) communicate(CONNECT string, reqType comm_request_t ) ([]byte, int) {
	// int rc;
	// struct timeval tv;
  
	var sample []byte
	req := e.prepare_request(reqType)
	log.Info("req: %s", req)

	e.mu.Lock()
	defer e.mu.Unlock()
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", CONNECT)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// perform communication
	resp, rc := multiple_attempt(
		func ( c net.Conn, r []byte) (string, int) {
			_, err = c.Write(r)
			if err != nil {
				return  "", -1
			}

			resp, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				return "", -1
			}
		  return resp, 0;
		},conn, req,
		NET_OP_MAX_ATTEMPT, NET_OP_RETRY_INTV)
  
	if rc == 0 {
		sample = e.parse_response(resp)
	}

	return sample, rc
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
// Attempt a function several times. Non-zero return of func is treated as an error. If func return
// -1, errno will be returned.

func (e *Exporter) prepare_request(r comm_request_t) []byte {

	var id int32

	pod_name := "alnair-exporter-nod\000"
	id = 0
	var a [NAME_LEN]byte

	copy(a[:], pod_name)	
	req := Packet{Len:uint64(len(a) - 1), Msg:a, ReqId:id, ReqType:r}
    out, err := json.Marshal(req)
    if err != nil {
        panic (err)
    }

    fmt.Println(string(out))
  
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, req)

	if err != nil {
		log.Errorf("Unable to create the request, Error: %s", err)
		return []byte("")
	}
    fmt.Println(buf.String())

	return buf.Bytes()
  }