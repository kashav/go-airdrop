package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
			if entry.Instance == name || len(entry.AddrIPv4) < 1 || entry.Text[1] == op {
				continue
			}

			switch op {
			case "list":
				if *listTypePtr != "all" && entry.Text[1] != *listTypePtr {
					continue
				}

				fmt.Printf("%s%s:%d\n",
					padRight(entry.Instance, " ", 30),
					entry.AddrIPv4[0],
					entry.Port)
			case "send":
				if entry.Text[1] != "broadcast" {
					continue
				}

				// Write to this entry if and only if we haven't seen it yet.
				if _, ok := seen[entry.Instance]; !ok {
					go dial(entry.AddrIPv4[0], entry.Port)
				}
			}
		}
	}(entries)

	var ctx context.Context
	var cancel context.CancelFunc

	if op == "list" && !*listWatchPtr {
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		fmt.Printf("Failed to browse: %s\n", err.Error())
	}

	<-ctx.Done()
}
