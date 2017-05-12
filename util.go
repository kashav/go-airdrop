package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dustinkirkland/golang-petname"
)

type stringList []string

func (s *stringList) String() string { return fmt.Sprintf("%v", *s) }
func (s *stringList) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}
func (s *stringList) Contains(client string) bool {
	for _, v := range *s {
		if v == client {
			return true
		}
	}
	return false
}

func parseFlags() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, errNoCmd.Error())
		os.Exit(1)
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
			fmt.Println(errNoFile.Error())
			os.Exit(1)
		}

		name = *sendNamePtr
		seen = make(map[string]bool, 0)
	default:
		fmt.Fprintln(os.Stderr, errNoCmd.Error())
		os.Exit(1)
	}

	for name == "" {
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
