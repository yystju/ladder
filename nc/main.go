package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var (
	mode    string
	address string
)

func init() {
	flag.StringVar(&mode, "m", "c", "The mode")
	flag.StringVar(&address, "a", ":1234", "The address of the server (as source or target).")
	flag.Parse()

	//log.Printf("mode : %s\n", mode)
}

func main() {
	if "c" == mode {
		doAsClient()
	} else {
		doAsServer()
	}
}

func doAsClient() {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Panic(err)
	}

	handler(conn)
}

func doAsServer() {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Panic(err)
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	go func() {
		io.Copy(conn, os.Stdin)
		conn.Close()
	}()

	io.Copy(os.Stdout, conn)
	conn.Close()
}
