package main

import (
	"fmt"
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/kshvmdn/rdrp"
	"github.com/kshvmdn/rdrp/version"
)

var (
	app   = kingpin.New("rdrp", "Send and receive files over your local network.")
	name  = app.Flag("name", "Set your connection name.").Short('n').String()
	debug = app.Flag("debug", "Enable debug mode.").Short('d').Bool()

	bc = app.Command("broadcast", "Receive a file.")

	list      = app.Command("list", "View active clients.")
	listWatch = list.Flag("watch", "Watch for new connections.").Short('w').Bool()
	listType  = list.Flag("type", "Specify which type of client to listen for.").Default("all").Short('t').String()

	send     = app.Command("send", "Send a file.")
	sendFile = send.Flag("file", "Specify the transfer file (you may optionally pass your file via stdin).").Short('f').String()
	sendList = send.Flag("to", "Comma-separated list of client names.").Strings()
)

func main() {
	app.Version(version.Version)
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if !*debug {
		devNull, err := os.Open(os.DevNull)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}
		log.SetOutput(devNull)
	}

	if *name == "" {
		*name = rdrp.GenerateName()
	}

	var r rdrp.Runner
	client := rdrp.Client{Command: command, Name: *name}

	switch command {
	case bc.FullCommand():
		r = rdrp.NewBroadcaster(client)
	case list.FullCommand():
		r = rdrp.NewLister(client, *listWatch, *listType)
	case send.FullCommand():
		r = rdrp.NewSender(client, *sendFile, *sendList)
	}

	if err := rdrp.Start(r); err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
	}
}
