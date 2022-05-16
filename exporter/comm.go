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

func (e *Exporter) communicate(CONNECT string, reqType comm_request_t ) (int, []byte) {
	// int rc;
	// struct timeval tv;
  
	var sample []byte
	req := e.prepare_request(reqType)
	// log.Info("req: %s", req)

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
				return 0, []byte("failed to read the response.")
			}
			// println("resp: ", string(resp))
			// sam := e.parse_sample(resp)
		
			// println("reply from server: ", string(sam), n)

		  return n, resp;
		}, conn, req,
		NET_OP_MAX_ATTEMPT, NET_OP_RETRY_INTV)

	log.Info("Resp: %s, Length: %d", resp, rc)

	if rc <= 0 {
		log.Fatalf("Failed to communicate with alnair server daemon")
		return rc,  []byte("")
	}
	rc, sample = e.parse_sample(resp)

	return rc, sample
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

  func (e *Exporter) parse_sample(resp []byte) (int, []byte) {
    var unpacked UnpackedSample

	buf := bytes.NewReader(resp)

    err := binary.Read(buf, binary.LittleEndian, &unpacked.id)
	if err != nil {
		log.Errorf("Unable to parse id, Error: %s", err)
		return 0, []byte("reqid is wrong")
	}

    err = binary.Read(buf, binary.LittleEndian, &unpacked.length)
	if err != nil {
		log.Errorf("Unable to parse sample length, Error: %s", err)
		return 0, []byte("response length is wrong")
	}    
    println("sample length: ", unpacked.length)
    unpacked.sample = make([]byte, unpacked.length)
    err = binary.Read(buf, binary.LittleEndian, &unpacked.sample)
	if err != nil {
		log.Errorf("Unable to parse sample, Error: %s", err)
        return 0, []byte("sample contents are wrong.")
	}    
	return int(unpacked.length), unpacked.sample
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
    out, err := json.Marshal(req)
    if err != nil {
        panic (err)
    }
	log.Info("Req: ", string(out))
  
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, req)

	if err != nil {
		log.Errorf("Unable to create the request, Error: %s", err)
		return []byte("")
	}
    // fmt.Println(buf.String())

	return buf.Bytes()
  }