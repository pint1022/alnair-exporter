package exporter

import (
	"bufio"
	"os"
	"sync"
	"fmt"
	"time"
	"net"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)
const NET_OP_MAX_ATTEMPT = 5  // maximum time retrying failed network operations
const NET_OP_RETRY_INTV = 10  // seconds between two retries

func multiple_attempt(f func() int, int max_attempt, int interval) int {
	rc := 0;
	for attempt = 1; attempt <= max_attempt; attempt++ {
	  rc = f();
	  if (rc == 0) {
		  break
	  }
	  log.Errorf("attempt %d", attempt)
	  if (interval > 0) {
		time.Sleep(interval * time.milliseconds)
	}
	return rc
  }

func (e *Exporter) communicate(conn net.Conn, sbuf string, rbuf *string ) int {
	// int rc;
	// struct timeval tv;
  
	mu.Lock()
	defer mu.Unlock()

	// perform communication
	rc := multiple_attempt(
		func () int {
			_, err = conn.Write(sbuf)
			if err != nil return -1

			*rbuf, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil return -1
		  return 0;
		},
		NET_OP_MAX_ATTEMPT, NET_OP_RETRY_INTV);
  
	return rc;
  }
