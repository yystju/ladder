package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIC/jCCAeYCCQCrqp7PYuzz1jANBgkqhkiG9w0BAQsFADBBMT8wPQYDVQQDDDZl
YzItMTMtMTEyLTgyLTE0NC5hcC1ub3J0aGVhc3QtMS5jb21wdXRlLmFtYXpvbmF3
cy5jb20wHhcNMTkwNTEyMDQxNDE2WhcNMjIwNTExMDQxNDE2WjBBMT8wPQYDVQQD
DDZlYzItMTMtMTEyLTgyLTE0NC5hcC1ub3J0aGVhc3QtMS5jb21wdXRlLmFtYXpv
bmF3cy5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCiYuNO96b3
8/75bEQAg0kvKR3CfoiahEHmGTQUralP7pFaN0eKqP/e6Mys8g0Mgm0RpDsx1Npq
s8X4duPMrCUEp11ESBz6UQacROTrXXsW7a3Ub+x+C5MoxtiEqJ08ktcMOQ99WSnB
rgxNwJUTH+ktIz0Lu/2O+/FeFYbGT5vLabI5aAbld7KG2ped6Ktfnry4Qdw4hhRG
D210C1Q/+0F4ri+Nhw4yebFxQa/NMiSezKTWNuAIVl2rNPyuHoCxQCiXsmDJdfH/
NbcLGr1t19PJ7F2QZfsC8BoXZOYMBFDA2GxODRMK/FYUPQJsN4zVnX3N/QD9GAZe
zRBzRe3i9e4VAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAFIJM5n996u1ZSFdeauk
s51MXAqZaVq7Fp2vRuGwFHCiHWkSLyaMGVJ4r8R98ZHp58xhyV/W5bDOknPIf5vY
DSTLibiSiPjnpVVTK1Tr12NM/GUsWVf98CZy4oPwzSTYRRyF8NXedZlJjIPDJ0Uh
hKNydRj87v5BlzpVjPt4wQCRf/SVmr4rfUNicjFXTkdre9cAAI/EU7dPBdZedYPE
AHpuWUtJugSurjODqElhcbEvFWS0JE2XTFbjeCVYy2rUAyCh0DQDwgTaEzIyYrpf
XS71Q/qyK01vsWDqyOnH5jElM9fwKET69LR0ij4DDZq9TZXUQCjL8rhSoLFbOLfP
LB0=
-----END CERTIFICATE-----`

var (
	server string
	listen string
)

func init() {
	flag.StringVar(&server, "s", "localhost:2080", "...")
	flag.StringVar(&listen, "l", "localhost:1080", "...")
	flag.Parse()

	log.Printf("server : %s, listen : %s\n", server, listen)
}

func main() {
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
		RootCAs:      roots,
		KeyLogWriter: os.Stdout,
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
