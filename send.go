package rdrp

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/grandcat/zeroconf"
)

// Sender represents a send client.
type Sender struct {
	Client Client

	FileName string
	Clients  []string
}

// NewSender returns a new Sender and zeroes the seen map.
func NewSender(client Client, fileName string, clients []string) *Sender {
	seen = map[string]bool{}

	return &Sender{
		Client:   client,
		FileName: fileName,
		Clients:  clients,
	}
}

// Work creates a server and starts discovering Broadcasters.
func (s *Sender) Work() (err error) {
	if server, err = s.Client.makeServer(); err != nil {
		return err
	}
	defer server.Shutdown()

	s.Client.printName()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return s.Client.startDiscovery(ctx, s.discover)
}

// discover starts dialing connected Broadcaster clients.
func (s *Sender) discover(entry *zeroconf.ServiceEntry) {
	if entry.Text[1] != "broadcast" {
		return
	}

	// If we've already seen this client, we don't want to resend the
	// request, so continue.
	if _, ok := seen[entry.Instance]; ok {
		return
	}

	// If clients are specified (--send-to) and the list _doesn't_ include
	// this entry, continue.
	if len(s.Clients) > 0 && !s.hasClient(entry.Instance) {
		return
	}

	go s.dial(entry.AddrIPv4[0], entry.Port)
}

// hasClient checks if this Sender's client list contains the provided client.
func (s *Sender) hasClient(client string) bool {
	for _, c := range s.Clients {
		if c == client {
			return true
		}
	}
	return false
}

// getSrcFile returns a pointer to the source file with an associated status
// string. If FileName is empty, return stdin; otherwise, open the file and
// return it.
func (s *Sender) getSrcFile(client string) (src *os.File, status string, err error) {
	if strings.TrimSpace(s.FileName) == "" {
		return os.Stdin, fmt.Sprintf("Sending file to %s...", client), nil
	}
	if src, err = os.Open(s.FileName); err != nil {
		return nil, "", err
	}
	return src, fmt.Sprintf("Sending %s to %s...", s.FileName, client), nil
}

// write initiates a conversation with the connected client and transfers the
// associated source file.
func (s *Sender) write(conn net.Conn) error {
	payload := s.FileName + separator + s.Client.Name
	if _, err := conn.Write([]byte(padRight(payload, padder, 100))); err != nil {
		return err
	}

	buf := make([]byte, 100)
	if _, err := conn.Read(buf); err != nil {
		return err
	}

	// response[0] -> client name, response[1] -> client response
	response := strings.Split(strings.Trim(string(buf), padder), separator)
	if len(response) < 2 {
		return nil
	}

	seen[response[0]] = true
	if response[1] != "Y" {
		return nil
	}

	src, status, err := s.getSrcFile(response[0])
	if err != nil {
		return err
	}

	fmt.Print(status)
	if _, err := io.Copy(conn, src); err != nil {
		fmt.Println("failed :(\a")
		return err
	}
	fmt.Println("done!")

	// Reset src offset for subsequent transfers.
	if _, err := src.Seek(0, 0); err != nil {
		return err
	}
	return nil
}

// dial creates a connection with client at addr:port.
func (s *Sender) dial(addr net.IP, port int) {
	conn, err := net.Dial(iptype, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if conn == nil {
		fmt.Printf("Failed to dial on %v.\n", addr)
		return
	}
	defer conn.Close()

	if err = s.write(conn); err != nil && err.Error() != "EOF" {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
