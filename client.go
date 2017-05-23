package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func readFile(conn net.Conn) error {
	buf := make([]byte, 100)
	_, err := conn.Read(buf)
	if err != nil {
		return err
	}

	rawResponse := strings.Split(strings.Trim(string(buf), padder), separator)
	if len(rawResponse) < 2 {
		return fmt.Errorf("")
	}

	file, connectionName := rawResponse[0], rawResponse[1]
	if strings.TrimSpace(file) == "" {
		file = "transfer request"
	}

	var input string
	fmt.Fprintf(os.Stderr, "Accept %s from %s? (Y/n) ", file, connectionName)
	fmt.Fscanf(os.Stderr, "%s", &input)

	rawPayload := fmt.Sprintf("%s%s%s", name, separator, input)
	_, err = conn.Write([]byte(padRight(rawPayload, padder, 100)))
	if err != nil {
		return err
	}

	if input != "Y" {
		return fmt.Errorf("")
	}

	io.Copy(os.Stdout, conn)
	conn.Close()
	return nil
}

func listen() {
	laddr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen(iptype, laddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		if conn == nil {
			fmt.Fprintf(os.Stderr, "Failed to listen on %v.\n", laddr)
		} else if err := readFile(conn); err != nil {
			if !(err.Error() == "" || err.Error() == "EOF") {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			conn.Close()
			continue
		}
		break
	}
}

func getSourceFile(client string) (*os.File, string, error) {
	if strings.TrimSpace(file) == "" {
		return os.Stdin, fmt.Sprintf("Sending file to %s...", client), nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, "", err
	}

	return f, fmt.Sprintf("Sending %s to %s...", file, client), nil
}

func writeFile(conn net.Conn) error {
	rawPayload := fmt.Sprintf("%s%s%s", file, separator, name)
	_, err := conn.Write([]byte(padRight(rawPayload, padder, 100)))
	if err != nil {
		return err
	}

	buf := make([]byte, 100)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	rawResponse := strings.Split(strings.Trim(string(buf), padder), separator)
	if len(rawResponse) < 2 {
		return nil
	}
	clientName, response := rawResponse[0], rawResponse[1]

	// FIXME: if the response is not `Y`, we set this client's `seen` value,
	// since we don't want to resend requests to clients who've declined. The
	// caveat here is that if they reconnect with the same name (after accepting),
	// they'll get a request again. Is this something to address?
	if response != "Y" {
		seen[clientName] = true
		return nil
	}

	src, status, err := getSourceFile(clientName)
	if err != nil {
		return err
	}

	fmt.Print(status)
	io.Copy(conn, src)
	fmt.Println("done!")

	return nil
}

func dial(addr net.IP, port int) {
	conn, err := net.Dial(iptype, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}
	if conn == nil {
		fmt.Printf("Failed to dial on %v.\n", addr)
	} else {
		if err = writeFile(conn); err != nil && err.Error() != "EOF" {
			fmt.Println(err.Error())
		}
		conn.Close()
	}
}
