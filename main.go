package main

import (
	"bufio"
	"github.com/davecgh/go-spew/spew"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	cfg.Parse()
	spew.Fprintf(os.Stderr, "started: %#v\n", cfg)

	statsdInit(*cfg.statsd_host, *cfg.statsd_port, *cfg.statsd_prefix, *cfg.statsd_tags)

	getStat()

	for *cfg.interval > 0 {
		time.Sleep(time.Duration(*cfg.interval) * time.Second)
		getStat()
	}
}

func getStat() {
	sTime := time.Now()

	defer func() {
		if e := recover(); e != nil {
			statsdIncrement("error")
			logInfo(0, sTime, "Error getting stat: %s", e.(error).Error())
		}
		statsdTiming(sTime, "done")
	}()

	nsInfo, err := asInfo(*cfg.aspk_node, *cfg.asinfo_port)
	if err != nil {
		logInfo(0, sTime, "Error getting node stat: %s", err.Error())
		statsdIncrement("error")
	} else {
		for nsName, nsStat := range *nsInfo {
			for key, val := range *nsStat {
				if *cfg.stdout {
					spew.Printf("%s.%s = %s\n", nsName, key, val)
				}
				if isIn(key, cfg.metric) {
					num, ok := a2i(val)
					if ok {
						statsdGauge(num, "%s.%s", nsName, key)
					}
				}
			}
		}
	}
}

func splitStatistics(str string) *map[string]string {
	elem := strings.Split(strings.Trim(str, "\n"), ";")
	result := make(map[string]string, len(elem))

	for _, param := range elem {
		keyVal := strings.SplitN(param, "=", 2)
		result[keyVal[0]] = keyVal[1]
	}

	return &result
}

func a2i(str string) (int, bool) {
	n, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		return int(n), true
	}

	b, err := strconv.ParseBool(str)
	if err == nil {
		if b {
			return 1, true
		}
		return 0, true
	}

	switch {
	case str == "on":
		return 1, true
	case str == "off":
		return 0, true
	}

	return 0, false
}

func isIn(v string, m []string) bool {
	if len(m) < 1 {
		return true
	}

	for _, s := range m {
		if v == s {
			return true
		}
	}

	return false
}

type nsStatType *map[string]string

func asInfo(srv string, prt uint) (*map[string]nsStatType, error) {
	conn, err := net.DialTimeout("tcp", spew.Sprintf("%s:%d", srv, prt), time.Second)

	if err != nil {
		return nil, err
	}

	defer func() {
		conn.Close()
	}()

	err = conn.SetDeadline(time.Now().Add(time.Second))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte("namespaces\n"))
	if err != nil {
		return nil, err
	}

	connbuf := bufio.NewReader(conn)
	str, err := connbuf.ReadString('\n')
	if err != nil {
		return nil, err
	}

	nsList := strings.Split(strings.Trim(str, "\n"), ";")

	info := make(map[string]nsStatType, len(nsList)+1)

	for _, nsName := range nsList {
		_, err = conn.Write([]byte("namespace/" + nsName + "\n"))
		if err != nil {
			return nil, err
		}
		str, err = connbuf.ReadString('\n')
		if err != nil {
			return nil, err
		}
		info["namespace."+nsName] = splitStatistics(str)
	}

	_, err = conn.Write([]byte("statistics\n"))
	if err != nil {
		return nil, err
	}

	str, err = connbuf.ReadString('\n')
	if err != nil {
		return nil, err
	}

	info["statistics"] = splitStatistics(str)

	return &info, nil
}
