package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	iptype  = "tcp"
	service = "_rdrp._tcp"
	domain  = "local."

	padder    = ":::"
	separator = ";;;"
)

var (
	port  int
	debug bool
	op    string
	name  string
	file  string
	seen  map[string]bool

	listTypePtr  *string
	listWatchPtr *bool

	sendClientList stringList

	server *zeroconf.Server

	errNoCmd  = fmt.Errorf("Usage: %s [broadcast|list|send] [options]", os.Args[0])
	errNoFile = errors.New("expected either file path or data through stdin")
)

func main() {
	var err error

	// Set seed for golang-petname (to generate instance names).
	rand.Seed(time.Now().UTC().UnixNano())

	if err = parseFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if !debug {
		if err = redirectLogOutput(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	if port, err = getOpenPort(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if op != "list" {
		fmt.Fprintf(os.Stderr, "Connected as %s.\n", name)
	}

	server = makeServer()
	defer server.Shutdown()

	switch op {
	case "broadcast":
		listen()
	case "list":
		fallthrough
	case "send":
		startDiscovery()
	}
}
