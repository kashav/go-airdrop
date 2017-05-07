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

	msg := strings.Split(strings.Trim(string(buf), padder), separator)
	if len(msg) < 2 {
		return fmt.Errorf("")
	}
	file, connectionName := msg[0], msg[1]

	var input string
	fmt.Fprintf(os.Stderr, "Accept %s from %s? (Y/n) ", file, connectionName)
	fmt.Fscanf(os.Stderr, "%s", &input)

	response := fmt.Sprintf("%s%s%s", name, separator, input)

	if input != "Y" {
		_, err = conn.Write([]byte(padRight(response, padder, 100)))
		if err != nil {
			return err
		}
		return fmt.Errorf("")
	}

	_, err = conn.Write([]byte(padRight(response, padder, 100)))
	if err != nil {
		return err
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
		if conn != nil {
			if err = readFile(conn); err != nil {
				if err.Error() != "" && err.Error() != "EOF" {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				conn.Close()
				continue
			}
		} else {
			fmt.Fprintf(os.Stderr, "Failed to listen on %v\n", laddr)
		}
		break
	}
}

func writeFile(conn net.Conn) error {
	_, err := conn.Write([]byte(padRight(fmt.Sprintf("%s%s%s", file, separator, name), padder, 100)))
	if err != nil {
		return err
	}

	buf := make([]byte, 100)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	msg := strings.Split(strings.Trim(string(buf), padder), separator)
	if len(msg) < 2 {
		return nil
	}
	clientName, response := msg[0], msg[1]
	seen[clientName] = true

	if response != "Y" {
		return nil
	}

	fmt.Printf("Sending %s to %s...", file, clientName)
	if f, err := os.Open(file); err == nil {
		io.Copy(conn, f)
		fmt.Println(" done!")
	}

	return nil
}

func dial(addr net.IP, port int) {
	conn, err := net.Dial(iptype, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}
	if conn != nil {
		if err = writeFile(conn); err != nil && err.Error() != "EOF" {
			fmt.Println(err.Error())
		}
		conn.Close()
	} else {
		fmt.Printf("Failed to dial on %v\n", addr)
	}
}
