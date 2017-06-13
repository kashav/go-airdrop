package rdrp

import (
	"context"
	"fmt"
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
	port int
	seen map[string]bool

	server *zeroconf.Server
)

// Runner represents a connected client.
type Runner interface {
	Work() error
}

// Start finds an open port and calls the runner's work function.
func Start(r Runner) (err error) {
	if port, err = getOpenPort(); err != nil {
		return err
	}
	return r.Work()
}

// Client represents common config. for each connected client.
type Client struct {
	Command string
	Name    string
}

// printName outputs the client's name.
func (c *Client) printName() {
	fmt.Fprintf(os.Stderr, "Connected as %s.\n", c.Name)
}

// makeServer registers the current Client.
func (c *Client) makeServer() (*zeroconf.Server, error) {
	return zeroconf.Register(
		c.Name,
		service,
		domain,
		port,
		[]string{"rdrp", c.Command, time.Now().Format(time.Stamp)},
		nil,
	)
}

// startDiscovery initiates the discovery process for the current Client.
// Utilized by any client searching for other clients on their respective
// service.
func (c *Client) startDiscovery(ctx context.Context, discoverFunc func(*zeroconf.ServiceEntry)) (err error) {
	var resolver *zeroconf.Resolver
	if resolver, err = zeroconf.NewResolver(nil); err != nil {
		return fmt.Errorf("failed to initialize resolver: %s", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == c.Name || len(entry.AddrIPv4) == 0 {
				continue
			}
			discoverFunc(entry)
		}
	}(entries)

	if err = resolver.Browse(ctx, service, domain, entries); err != nil {
		return fmt.Errorf("failed to browse: %s", err.Error())
	}

	<-ctx.Done()
	return nil
}
