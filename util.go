package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/dustinkirkland/golang-petname"
)

func parseFlags() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, errNoCmd.Error())
		os.Exit(1)
	}

	broadcastCmd := flag.NewFlagSet("broadcast", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	broadcastNamePtr := broadcastCmd.String("name", "", "Display name")

	sendNamePtr := sendCmd.String("name", "", "Display name")
	sendFilePtr := sendCmd.String("file", "", "File to transfer")

	op = os.Args[1]
	switch op {
	case "broadcast":
		broadcastCmd.Parse(os.Args[2:])
		name = *broadcastNamePtr
	case "send":
		sendCmd.Parse(os.Args[2:])
		if *sendFilePtr == "" {
			fmt.Fprintln(os.Stderr, errNoFile.Error())
			os.Exit(1)
		}

		name = *sendNamePtr
		file = *sendFilePtr
		seen = make(map[string]bool, 0)
	default:
		fmt.Fprintln(os.Stderr, errNoCmd.Error())
		os.Exit(1)
	}

	if name == "" {
		name = petname.Generate(2, "-")
	}
}

func getOpenPort() (int, error) {
	addr, err := net.ResolveTCPAddr(iptype, ":0")
	if err != nil {
		return -1, err
	}

	l, err := net.ListenTCP(iptype, addr)
	if err != nil {
		return -1, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) >= length {
			return str[0:length]
		}
	}
}
