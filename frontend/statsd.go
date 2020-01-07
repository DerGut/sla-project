package main

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"log"
	"net/http"
)

const unixSocketPath = "unix:///var/run/datadog/dsd.socket"

func NewStatsD() (*statsd.Client, error) {
	c, err := statsd.New(unixSocketPath)
	if err != nil {
		return nil, err
	}
	c.Namespace = "frontend."
	return c, nil
}

func CountRequest(r *http.Request) {
	err := sd.Count(
		"request",
		1,
		[]string{
			fmt.Sprintf("path:%s", r.URL.Path),
		},
		1)
	if err != nil {
		log.Fatalf("Couldn't sed count to statsd: %s", err)
	}
	log.Printf("Sent count to statsd")
}
