package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	// mode   string
	target string
	local  string
)

func init() {
	log.SetPrefix("TEST>")
	log.Println("[INIT]")

	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
	// flag.StringVar(&mode, "m", "server", "The running mode: client or server...")
	flag.StringVar(&target, "d", ":1080", "The target...")
	flag.StringVar(&local, "l", ":2000", "The local...")
	flag.Parse()
}

func main() {
	fmt.Printf("TLS wrapper... target : %s\n", target)

	cert, err := tls.LoadX509KeyPair("testdata/example-cert.pem", "testdata/example-key.pem")

	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		KeyLogWriter: os.Stdout,
	}

	listener, err := tls.Listen("tcp", local, cfg)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Panic(err)
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	log.Printf("[handled]\n")

	remote, err := net.Dial("tcp", target)

	if err != nil {
		log.Panic(err)
	}

	cWriter := NewWriterWrapper(conn)
	rReader, err := remote, nil

	if err != nil {
		log.Panic(err)
	}

	rWriter := remote
	cReader, err := NewReaderWrapper(conn), nil

	if err != nil {
		log.Panic(err)
	}

	go func() {
		io.Copy(cWriter, rReader)

		if remote != nil {
			remote.Close()
		}
		if conn != nil {
			conn.Close()
		}
	}()

	io.Copy(rWriter, cReader)

	if remote != nil {
		remote.Close()
	}
	if conn != nil {
		conn.Close()
	}

	log.Println("-DISCONNECTED-")
}
