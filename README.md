## rdrp [![Build Status](https://travis-ci.org/kashav/rdrp.svg?branch=master)](https://travis-ci.org/kashav/rdrp) [![Go Report Card](https://goreportcard.com/badge/github.com/kashav/rdrp)](https://goreportcard.com/report/github.com/kashav/rdrp)

> A cross-platform command line tool for sending and receiving files over your local network, inspired by [AirDrop](https://support.apple.com/en-ca/HT204144).

### Contents

  - [Demo](#demo)
  - [Design](#design)
  - [Installation](#installation--setup)
  - [Usage](#usage)
    + [Send](#sender)
    + [Broadcast](#broadcaster)
    + [Monitor](#list)
  - [Docker](#use-with-docker)
  - [Contribute](#contribute)
  - [Related](#related)
  - [License](#license)

### Demo

<a href="https://asciinema.org/a/120148"><img src="./media/rdrp.gif"></a>

### Design

rdrp uses [Multicast DNS](https://en.wikipedia.org/wiki/Multicast_DNS) to enable peer-to-peer discovery between clients. This means that rdrp will likely **not** work in most cloud/virtual environments.

When a client first connects, they're registered as a new instance on the `_rdrp._tcp` service. Each sender continuously browses this service for newly connected broadcasters with whom they'll establish a connection and attempt to send their respective file.

This program implements mDNS with [grandcat/zeroconf](https://github.com/grandcat/zeroconf).

Read more about mDNS: [RFC 6762](https://tools.ietf.org/html/rfc6762) and DNS-SD: [RFC 6763](https://tools.ietf.org/html/rfc6763).

### Installation / setup

Go should be [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).

Install with Go:

  ```sh
  $ go get -v github.com/kashav/rdrp/...
  $ which rdrp
  $GOPATH/bin/rdrp
  ```

Or, install directly via source:

  ```sh
  $ git clone https://github.com/kashav/rdrp.git $GOPATH/src/github.com/kashav/rdrp
  $ cd $_ # $GOPATH/src/github.com/kashav/rdrp
  $ make install all
  $ ./rdrp
  ```

### Usage

Run rdrp with the `--help` flag to view the usage dialogue.

  ```sh
  $ rdrp --help
  usage: rdrp [<flags>] <command> [<args> ...]

  Send and receive files over your local network.

  Flags:
        --help       Show context-sensitive help (also try --help-long and --help-man).
    -n, --name=NAME  Set your connection name.
    -d, --debug      Enable debug mode.
        --version    Show application version.

  Commands:
    help [<command>...]
      Show help.

    broadcast
      Receive a file.

    list [<flags>]
      View active clients.

    send [<flags>]
      Send a file.

  ```

There's two parties involved in a single transaction: the [sender](#send) and the [receiver](#broadcast).

#### Send

To send a file, use the `send` command. Provide the file path with the `--file` flag or pass the file's contents via stdin. 

Every broadcaster will receive a request to transfer the file (unless names are specified with the `--to` flag). This process continues until aborted (Ctrl+C).

  ```sh
  $ rdrp help send
  usage: rdrp send [<flags>]

  Send a file.

  Flags:
        --help       Show context-sensitive help (also try --help-long and --help-man).
    -n, --name=NAME  Set your connection name.
    -d, --debug      Enable debug mode.
        --version    Show application version.
    -f, --file=FILE  Specify the transfer file (you may optionally pass your file via stdin).
        --to=TO ...  Comma-separated list of client names.

  ```

##### Examples

```sh
$ rdrp send --file=README.md
```

```sh
$ rdrp send --name sender < README.md
```

```sh
$ tar -cvzf archive.tar.gz /path/to/directory/
$ rdrp send --file=archive.tar.gz --to=a
```

```sh
$ echo "hello" | rdrp send --to=b,c
```

#### Broadcast

To broadcast yourself as a receiver (i.e. someone receiving a file), use the `broadcast` command.

You'll be listening for incoming `send` requests. Upon a new connection, you'll be prompted on whether you'd like to accept the file or not, just like AirDrop. The incoming file is copied to stdout.

  ```sh
  $ rdrp broadcast -help
  usage: rdrp broadcast

  Receive a file.

  Flags:
        --help       Show context-sensitive help (also try --help-long and --help-man).
    -n, --name=NAME  Set your connection name.
    -d, --debug      Enable debug mode.
        --version    Show application version.

  ```

##### Examples

  ```sh
  $ rdrp broadcast # output is copied to stdout
  ...
  ```

  ```sh
  $ rdrp broadcast --name b > archive.tar.gz
  ```

Note that each of the above roles has an **optional** name flag, a name is chosen at random if not provided (which is what happened in [the demo above](#demo)).

#### List

You can view all connected clients with `list`. Use `--type` to specify the type of clients to list and `--watch` to listen for new connections.

  ```
  $ rdrp list -help

  View active clients.

  Flags:
        --help        Show context-sensitive help (also try --help-long and --help-man).
    -n, --name=NAME   Set your connection name.
    -d, --debug       Enable debug mode.
        --version     Show application version.
    -w, --watch       Watch for new connections.
    -t, --type="all"  Specify which type of client to listen for.

  ```

### Use with Docker

Start off by cloning the repository (if you've already cloned, navigate to the project root):

  ```sh
  $ git clone https://github.com/kashav/rdrp
  $ cd rdrp
  ```

Build the Docker image:

  ```sh
  $ docker build -t kashav/rdrp .
  ```

And run it! The `--rm` flag automatically removes the container when the program exits.

  ```sh
  $ docker run --rm kashav/rdrp [broadcast|list|send] ...
  ```

### Contribute

This project is completely open source, feel free to [open an issue](https://github.com/kashav/rdrp/issues) or [submit a pull request](https://github.com/kashav/rdrp/pulls).

Before submitting a PR, please ensure that _tests are passing_ and that the linter is happy. The following commands may be of use.

```sh
$ make install \
       get-tools
$ make fmt \
       vet \
       lint
$ make test \
       coverage
```

The demo GIF was generated with [asciinema](https://asciinema.org/), with [tmux](https://tmux.github.io/).

### Related

- [taku-k/airdrop](https://github.com/taku-k/airdrop)
- [tlehman/zerocat](https://github.com/tlehman/zerocat)

### License

rdrp source code is released under the [MIT License](./LICENSE).
