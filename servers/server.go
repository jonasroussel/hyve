package servers

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"

	"github.com/jonasroussel/proxbee/stores"
)

func NewTLS() (net.Listener, *http.Server, *http.ServeMux) {
	handler := http.NewServeMux()

	listener, err := tls.Listen("tcp", ":443", &tls.Config{
		GetCertificate: getCertificate,
		NextProtos:     []string{"h2", "http/1.1"},
	})
	if err != nil {
		panic(err)
	}

	server := &http.Server{Handler: handler}

	return listener, server, handler
}

func NewHTTP() (net.Listener, *http.Server, *http.ServeMux) {
	handler := http.NewServeMux()

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}

	server := &http.Server{Handler: handler}

	return listener, server, handler
}

// TODO : maybe need to be cached
func getCertificate(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if chi.ServerName == "" {
		return nil, errors.New("server name (sni) is empty")
	}

	cert, err := stores.Active.GetCertificate(chi.ServerName)
	if err != nil {
		return nil, err
	}

	x509Cert, err := tls.X509KeyPair(cert.CertificateData, cert.PrivateKeyData)
	if err != nil {
		return nil, err
	}

	return &x509Cert, nil
}
