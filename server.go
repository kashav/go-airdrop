package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grandcat/zeroconf"
)

func makeServer() *zeroconf.Server {
	server, err := zeroconf.Register(
		name,
		service,
		domain,
		port,
		[]string{"rdrp", op, time.Now().String()},
		nil)
	if err != nil {
		panic(err)
	}
	return server
}

func startDiscovery() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize resolver: %s\n", err.Error())
		os.Exit(1)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == name || len(entry.AddrIPv4) < 1 {
				continue
			}

			switch op {
			case "list":
				if *listTypePtr != "all" && entry.Text[1] != *listTypePtr {
					continue
				}

				t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", entry.Text[2])
				if err != nil {
					continue
				}

				fmt.Printf(
					"%s:%d %s %s %s\n",
					entry.AddrIPv4[0],
					entry.Port,
					padRight(entry.Text[1], " ", 9),
					t.Round(time.Second).Format("Jan _2 15:04"),
					entry.Instance)

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
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Millisecond)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to browse: %s\n", err.Error())
	}

	<-ctx.Done()
}
