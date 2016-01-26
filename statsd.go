package main

import (
    "github.com/onokonem/statsd-go"
	"time"
    "os"
    "fmt"
    "bytes"
    "strings"
)

var hostname string
var statsd_client *statsd.StatsdClient

func statsdInit(host string, port int) {
    hostname, _ = os.Hostname()
    hostname = strings.Replace(hostname, ".", "_", -1)
    statsd_client = statsd.New(host, port)
}

func statsdIncrement(s string, p ...interface{}) {
    statsd_client.Increment(concatName(s, p...))
}

func statsdGauge(v int, s string, p ...interface{}) {
    statsd_client.Gauge(concatName(s, p...), v)
}

func statsdIncrementByValue(v int, s string, p ...interface{}) {
    statsd_client.IncrementByValue(concatName(s, p...), v)
}

func statsdTiming(sTime time.Time, s string, p ...interface{}) {
    statsd_client.Timing(concatName(s, p...), int64(time.Now().Sub(sTime) / time.Millisecond))
}

func concatName(s string, p ...interface{}) string {
    if len(p) > 0 {
        s = fmt.Sprintf(s, p...)
    }

    var buffer bytes.Buffer
    buffer.WriteString(s)
    buffer.WriteString(".")
    buffer.WriteString(hostname)
    
    return buffer.String()
}