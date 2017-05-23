package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/dustinkirkland/golang-petname"
)

// stringList is used to handle a comma-separated input list.
type stringList []string

func (s *stringList) String() string         { return fmt.Sprintf("%v", *s) }
func (s *stringList) Set(value string) error { *s = strings.Split(value, ","); return nil }
func (s *stringList) Contains(client string) bool {
	for _, v := range *s {
		if v == client {
			return true
		}
	}
	return false
}

// Redirect all log output to /dev/null.
func redirectLogOutput() error {
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return err
	}
	log.SetOutput(devNull)
	return nil
}

// Parse the command line flags.
func parseFlags() error {
	if len(os.Args) < 2 {
		return errNoCmd
	}

	broadcastCmd := flag.NewFlagSet("broadcast", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	broadcastNamePtr := broadcastCmd.String("name", "", "Connection name")

	listTypePtr = listCmd.String("type", "all", "Type of client (\"broadcast\", \"send\", or \"all\")")
	listWatchPtr = listCmd.Bool("watch", false, "Listen for new connections (use Ctrl+C to exit)")

	sendNamePtr := sendCmd.String("name", "", "Connection name")
	sendFilePtr := sendCmd.String("file", "", "File to transfer")
	sendCmd.Var(&sendClientList, "send-to", "Comma-separated list of client names")

	op = os.Args[1]
	switch op {
	case "broadcast":
		broadcastCmd.Parse(os.Args[2:])
		name = *broadcastNamePtr
	case "list":
		listCmd.Parse(os.Args[2:])
	case "send":
		sendCmd.Parse(os.Args[2:])

		stat, _ := os.Stdin.Stat()
		// If `-file` is empty and stdin is empty, exit with errNoFile.
		if *sendFilePtr != "" {
			file = *sendFilePtr
		} else if (stat.Mode() & os.ModeCharDevice) != 0 {
			return errNoFile
		}

		name = *sendNamePtr
		seen = make(map[string]bool, 0)
	default:
		return errNoCmd
	}

	for name == "" {
		name = petname.Generate(2, "-")
	}

	return nil
}

// Get the next open port.
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

// Right pad a string length characters.
func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) >= length {
			return str[0:length]
		}
	}
}
