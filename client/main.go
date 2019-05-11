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
MIIC+TCCAeGgAwIBAgIUOnOKS798wyPowATY10/oMtK9AfgwDQYJKoZIhvcNAQEL
BQAwDDEKMAgGA1UEAwwBKjAeFw0xOTA1MTAxMzEyMzJaFw0yMjA1MDkxMzEyMzJa
MAwxCjAIBgNVBAMMASowggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCi
YuNO96b38/75bEQAg0kvKR3CfoiahEHmGTQUralP7pFaN0eKqP/e6Mys8g0Mgm0R
pDsx1Npqs8X4duPMrCUEp11ESBz6UQacROTrXXsW7a3Ub+x+C5MoxtiEqJ08ktcM
OQ99WSnBrgxNwJUTH+ktIz0Lu/2O+/FeFYbGT5vLabI5aAbld7KG2ped6Ktfnry4
Qdw4hhRGD210C1Q/+0F4ri+Nhw4yebFxQa/NMiSezKTWNuAIVl2rNPyuHoCxQCiX
smDJdfH/NbcLGr1t19PJ7F2QZfsC8BoXZOYMBFDA2GxODRMK/FYUPQJsN4zVnX3N
/QD9GAZezRBzRe3i9e4VAgMBAAGjUzBRMB0GA1UdDgQWBBRwmc+OLZnIu8U7RBpF
ItDgYtn96TAfBgNVHSMEGDAWgBRwmc+OLZnIu8U7RBpFItDgYtn96TAPBgNVHRMB
Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAWf+f7uyiJzUn/SZnuuHh3Egg/
DaNrAhThAiglnoqA3s/cTtlPOabPTkhIeREWD3DkA1d39yvT+fxdd6mrxW1Rbwmi
GmlpChfT0zXCfeUYVsmAwjphSVZ8pM9h4kS/EP7Vyl5jqMz6QiuTgcPSqAAewPQb
UGe4+5gaxuzXb2anOJTK196PSDx3bxzvn6n56FcV36SAdldkGFoSVn4WBjjpNcJo
UfRP0lZZf9jcWRLDBKSxxuGmA8McznipVp3xpIThj13CABA0ZutSyAMUqawiDxg6
/bQe6s6A1NOR0Enc4pX2ZWyxhppq4j7mc07UqBM0uscMS6nbV6uqx9oiadRN
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

	defer remote.Close()
	defer conn.Close()

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

	go io.Copy(cWriter, rReader)
	io.Copy(rWriter, cReader)
}
