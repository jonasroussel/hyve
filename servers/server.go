package servers

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jonasroussel/hyve/caching"
)

type LoggerWriter struct{}

func (w LoggerWriter) Write(p []byte) (n int, err error) {
	if string(p) == "EOF" || string(p) == "certificate not found" {
		return 0, nil
	} else {
		return os.Stdout.Write(p)
	}
}

func NewTLS() (net.Listener, *http.Server, *http.ServeMux) {
	handler := http.NewServeMux()

	listener, err := tls.Listen("tcp", ":443", &tls.Config{
		GetCertificate: caching.CertificateRetriever,
		NextProtos:     []string{"h2", "h2c", "http/1.1", "http/1.0", "spdy/2", "spdy/3"},
	})
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler:  handler,
		ErrorLog: log.New(LoggerWriter{}, "[TLS] ", log.LstdFlags),
	}

	return listener, server, handler
}

func NewHTTP() (net.Listener, *http.Server, *http.ServeMux) {
	handler := http.NewServeMux()

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{Handler: handler}

	return listener, server, handler
}
