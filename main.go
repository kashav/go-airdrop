package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	iptype  = "tcp"
	service = "_rdrp._tcp"
	domain  = "local."

	padder    = ":"
	separator = ";"
)

var (
	port int
	op   string
	name string
	file string
	seen map[string]bool

	server *zeroconf.Server

	errNoCmd  = errors.New("expected one of: broadcast or send")
	errNoFile = errors.New("broadcast command should include path of file")
)

func main() {
	// Redirect all log output to /dev/null.
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	log.SetOutput(devNull)

	// Set seed for golang-petname (to generate instance names).
	rand.Seed(time.Now().UTC().UnixNano())

	parseFlags()
	port, err = getOpenPort()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Connected as %s.\n", name)

	server = mkServer()
	defer server.Shutdown()

	switch op {
	case "broadcast":
		listen()
	case "send":
		startDiscovery()
	}
}
