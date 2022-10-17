package rdrp

import (
	"math/rand"
	"net"
	"strings"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
)

// GenerateName sets the seed and generates a random 2-word, hyphen-separated
// name.
func GenerateName() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return petname.Generate(2, "-")
}

// getOpenPort finds an open port on the host machine.
func getOpenPort() (int, error) {
	addr, err := net.ResolveTCPAddr(iptype, ":0")
	if err != nil {
		return -1, err
	}

	listener, err := net.ListenTCP(iptype, addr)
	if err != nil {
		return -1, err
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}

// padRight right-pads a string by n characters.
func padRight(str, pad string, n int) string {
	for {
		str += pad
		if len(str) >= n {
			return str[0:n]
		}
	}
}

// isYes returns the parsed answer to a (Y/n) prompt.
func isYes(response string) bool {
	response = strings.TrimSpace(response)
	return response == "" || response == "Y" || response == "y";
}
