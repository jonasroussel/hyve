package servers

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/jonasroussel/hyve/caching"
)

func NewTLS() (net.Listener, *http.Server, *http.ServeMux) {
	handler := http.NewServeMux()

	listener, err := tls.Listen("tcp", ":443", &tls.Config{GetCertificate: caching.CertificateRetriever})
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{Handler: handler}

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
