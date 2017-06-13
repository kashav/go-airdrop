package rdrp

import (
	"context"
	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
)

// Lister represents a list client.
type Lister struct {
	Client Client

	Watch      bool
	ClientType string
}

// NewLister returns a new Lister.
func NewLister(client Client, watch bool, clientType string) *Lister {
	return &Lister{
		Client:     client,
		Watch:      watch,
		ClientType: clientType,
	}
}

// Work creates a server and starts finding and listing connected clients.
func (l *Lister) Work() (err error) {
	if server, err = l.Client.makeServer(); err != nil {
		return err
	}
	defer server.Shutdown()

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	// Use timeout-based context iff we're NOT watching, otherwise cancel-based.
	if l.Watch {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*2)
	}
	defer cancel()

	return l.Client.startDiscovery(ctx, l.discover)
}

func (l *Lister) discover(entry *zeroconf.ServiceEntry) {
	if l.ClientType != "all" && l.ClientType != entry.Text[1] {
		return
	}

	t, err := time.Parse(time.Stamp, entry.Text[2])
	if err != nil {
		return
	}

	fmt.Printf(
		"%s:%d %s %s %s\n",
		entry.AddrIPv4[0],
		entry.Port,
		padRight(entry.Text[1], " ", 9),
		t.Round(time.Second).Format("Jan _2 15:04"),
		entry.Instance,
	)
}
