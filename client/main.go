package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const DEFAULT_CREDENTIAL = `
-----BEGIN CERTIFICATE-----
MIIDVTCCAj2gAwIBAgIJAMZ1iGRzPYBZMA0GCSqGSIb3DQEBCwUAMEExPzA9BgNV
BAMMNmVjMi0xMy0xMTItODItMTQ0LmFwLW5vcnRoZWFzdC0xLmNvbXB1dGUuYW1h
em9uYXdzLmNvbTAeFw0xOTA1MTIwNDEzMjZaFw0yMjA1MTEwNDEzMjZaMEExPzA9
BgNVBAMMNmVjMi0xMy0xMTItODItMTQ0LmFwLW5vcnRoZWFzdC0xLmNvbXB1dGUu
YW1hem9uYXdzLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKJi
4073pvfz/vlsRACDSS8pHcJ+iJqEQeYZNBStqU/ukVo3R4qo/97ozKzyDQyCbRGk
OzHU2mqzxfh248ysJQSnXURIHPpRBpxE5OtdexbtrdRv7H4LkyjG2ISonTyS1ww5
D31ZKcGuDE3AlRMf6S0jPQu7/Y778V4VhsZPm8tpsjloBuV3sobal53oq1+evLhB
3DiGFEYPbXQLVD/7QXiuL42HDjJ5sXFBr80yJJ7MpNY24AhWXas0/K4egLFAKJey
YMl18f81twsavW3X08nsXZBl+wLwGhdk5gwEUMDYbE4NEwr8VhQ9Amw3jNWdfc39
AP0YBl7NEHNF7eL17hUCAwEAAaNQME4wHQYDVR0OBBYEFHCZz44tmci7xTtEGkUi
0OBi2f3pMB8GA1UdIwQYMBaAFHCZz44tmci7xTtEGkUi0OBi2f3pMAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQELBQADggEBAGyJx4fDxyw5nL7oZRbgLUK3DUBbFq67
stA/pJ6/QZIKm7cph603spWqZ+1OjNsDiEyp+pH1zDP4QPujwCvqZ4kZl65WT8V6
9EvazmiDPJIQ0kjtEvEjwEG7he00VryXsXeYqlbp/rGJ/covtgTaYR6ToGWEE1/W
NX0AYx1zi0E9zW7x5KgDTViUadyGRpwjARwlEhNRPu8EEWpRLFvtjhisC/572PZ9
MBBad6mxvmd3ky8OtD473Uw9m5uNBePOKGAsg3+BoZl5k96YtQ6c8BoXMnOB1BDN
O08sW3NjuUaBZyq5dB3+USUkbhGVlSyP0W83FOMtQjXaKR+UovbDg9s=
-----END CERTIFICATE-----
`

var (
	server     string
	listen     string
	credential string
)

func init() {
	flag.StringVar(&server, "s", "localhost:2080", "Server address")
	flag.StringVar(&listen, "l", "localhost:1080", "Local listener address")
	flag.StringVar(&credential, "c", "", "Credential File Path Name")
	flag.Parse()

	log.Printf("server : %s, listen : %s\n", server, listen)
}

func main() {
	rootPEM := DEFAULT_CREDENTIAL

	if "" != credential {
		b, err := ioutil.ReadFile(credential)

		if err != nil {
			log.Panic(err)
		}

		rootPEM = string(b)
	}

	// log.Println(rootPEM)

	roots := x509.NewCertPool()

	ok := roots.AppendCertsFromPEM([]byte(rootPEM))

	if !ok {
		log.Panic("failed to parse root certificate")
	}

	ln, err := net.Listen("tcp", listen)

	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Panic(err)
		}

		go handler(server, conn, roots)
	}
}

func handler(server string, conn net.Conn, roots *x509.CertPool) {
	remote, err := tls.Dial("tcp", server, &tls.Config{
		RootCAs: roots,
		// KeyLogWriter: os.Stdout,
	})

	if err != nil {
		log.Printf("ERR : %s\n", err.Error())
		return
	}

	cWriter := conn
	rReader, err := NewReaderWrapper(remote), nil

	if err != nil {
		log.Panic(err)
	}

	rWriter := NewWriterWrapper(remote)
	cReader, err := conn, nil

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
}
