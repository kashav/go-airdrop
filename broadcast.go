package rdrp

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Broadcaster represents a broadcast client.
type Broadcaster struct {
	Client Client
}

// NewBroadcaster returns a new Broadcaster.
func NewBroadcaster(client Client) *Broadcaster {
	return &Broadcaster{client}
}

// Work creates a server and starts listening for file transfer requests.
func (b Broadcaster) Work() (err error) {
	if server, err = b.Client.makeServer(); err != nil {
		return err
	}
	defer server.Shutdown()
	b.Client.printName()
	return b.listen()
}

// listen listens for transfer requests and initiates a conversation upon
// new connections.
func (b *Broadcaster) listen() error {
	laddr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen(iptype, laddr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		if conn == nil {
			fmt.Fprintf(os.Stderr, "Failed to listen on %v.\n", laddr)
			break
		}
		defer conn.Close()

		if ok, err := b.read(conn); err != nil {
			return err
		} else if !ok {
			continue
		}
		break
	}

	return nil
}

// read interacts with Sender at conn and copies the associated file to stdout.
func (b *Broadcaster) read(conn io.ReadWriter) (ok bool, err error) {
	buf := make([]byte, 100)
	if _, err = conn.Read(buf); err != nil {
		return false, err
	}

	// request[0] -> file name, request[1] -> connection name
	request := strings.Split(strings.Trim(string(buf), padder), separator)
	if len(request) < 2 {
		return false, nil
	}

	if strings.TrimSpace(request[0]) == "" {
		request[0] = "transfer request"
	}

	var response string
	fmt.Fprintf(os.Stderr, "Accept %s from %s? (Y/n) ", request[0], request[1])
	fmt.Fscanf(os.Stderr, "%s", &response)

	payload := fmt.Sprintf("%s%s%s", b.Client.Name, separator, response)
	if _, err = conn.Write([]byte(padRight(payload, padder, 100))); err != nil {
		return false, err
	}

	// Must come after writing the payload!
	if !isYes(response) {
		return false, nil
	}

	if _, err := io.Copy(os.Stdout, conn); err != nil {
		return false, err
	}

	return true, nil
}
