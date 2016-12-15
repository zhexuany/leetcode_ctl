package main

import (
	"strings"
	"time"
)

const (
	DefaultPort = "345678"
)

type HTTPConfig struct {
	// Addr should be of the form "http://host:port"
	// or "http://[ipv6-host%zone]:port".
	Addr string

	// Username is the influxdb username, optional
	Username string

	// Password is the influxdb password, optional
	Password string

	// Timeout for influxdb writes, defaults to no timeout
	Timeout time.Duration
}

func NewDemoHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Addr:    "http://localhost:40000",
		Timeout: time.Duration(1000),
	}
}

func ParseCommandName(args []string) (string, []string) {
	// Retrieve command name as first argument
	var name string
	if len(args) > 0 {
		if !strings.HasPrefix(args[0], "-") {
			name = args[0]
		} else if args[0] == "-h" || args[0] == "-help" || args[0] == "--help" {
			name = "help"
		}
	}

	// If command is "help" and has an argument then rewrite args to use "-h"
	if name == "help" && len(args) > 2 && !strings.HasPrefix(args[1], "-") {
		return args[1], []string{"-h"}
	}

	// If a named command is specified then return it with its arguments.
	if name != "" {
		return name, args[1:]
	}
	return "", args
}
