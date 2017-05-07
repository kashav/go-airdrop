package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grandcat/zeroconf"
)

func mkServer() *zeroconf.Server {
	server, err := zeroconf.Register(
		name,
		service,
		domain,
		port,
		[]string{"rdrp", op},
		nil)
	if err != nil {
		panic(err)
	}
	return server
}

func startDiscovery() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		fmt.Printf("Failed to initialize resolver: %s\n", err.Error())
		os.Exit(1)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == name || len(entry.AddrIPv4) < 1 || op == entry.Text[1] {
				continue
			}

			// Write to this entry if and only if we haven't seen it yet.
			if _, ok := seen[entry.Instance]; !ok {
				go dial(entry.AddrIPv4[0], entry.Port)
			}
		}
	}(entries)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		fmt.Printf("Failed to browse: %s\n", err.Error())
	}
	<-ctx.Done()
}
