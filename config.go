package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	// common parameters
	aspk_node   *string
	asinfo_port *uint
	interval    *uint

	stdout *bool

	statsd_prefix *string
	metric strList

	statsd_host *string
	statsd_port *int
}

var cfg = Config{
	aspk_node:   flag.String("aspk_node", "127.0.0.1", "Aerospike node address"),
	asinfo_port: flag.Uint("asinfo_port", 3003, "Aerospike info port"),
	interval:    flag.Uint("interval", 10, "Stat send interval interval (sec)"),

	stdout: flag.Bool("stdout", false, "Output metrics to STDOUT (all of them regardless of 'metric' parameters)"),

	statsd_prefix: flag.String("statsd_prefix", "stats.aerospike.", "StatsD metrics name prefix"),
	statsd_tags:   flag.Bool("statsd_tags", false, "StatsD metrics name with tags (,host=)"),

	statsd_host: flag.String("statsd_host", "127.0.0.1", "StatsD host"),
	statsd_port: flag.Int("statsd_port", 8125, "StatsD port"),
}

func (*Config) Parse() {
	flag.Var(&(cfg.metric), "metric", "Metric to report (more than one can be provided), default report all")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()

	}

	flag.Parse()
}

type strList []string

func (s *strList) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strList) Set(v string) error {
	*s = append(*s, v)
	return nil
}
