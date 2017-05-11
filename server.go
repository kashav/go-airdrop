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

				// FIXME: is 30 an appropriate amount?
				fmt.Printf("%s%s:%d\n",
					padRight(entry.Instance, " ", 30),
					entry.AddrIPv4[0],
					entry.Port)
			case "send":
				if entry.Text[1] != "broadcast" {
					continue
				}

				// If we've already seen this client, we don't want to resend the
				// request, so continue.
				if _, ok := seen[entry.Instance]; ok {
					continue
				}

				// If clients are specified (-send-to) and the list _doesn't_ include
				// this entry, continue.
				if len(sendClientList) > 0 && !sendClientList.Contains(entry.Instance) {
					continue
				}

				go dial(entry.AddrIPv4[0], entry.Port)
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
