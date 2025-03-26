package core

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/jonasroussel/hyve/tools"
)

type Server struct {
	Listener net.Listener
	Handler  *tools.HTTPHandler
}

func (server *Server) Serve() error {
	err := http.Serve(server.Listener, server.Handler)
	if err != nil {
		return err
	}

	return nil
}

func NewHTTPSServer(app *HyveApp) (*Server, error) {
	handler := tools.NewHTTPHandler()

	http.NewServeMux()

	httpsPort := os.Getenv("HYVE_HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = "443"
	}

	listener, err := tls.Listen("tcp", "0.0.0.0:"+httpsPort, &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if info.ServerName == "" {
				return nil, errors.New("server name (sni) is empty")
			}

			cert, err := app.Store.GetCertificate(info.ServerName)
			if err != nil {
				return nil, err
			}

			x509Cert, err := tls.X509KeyPair(cert.CertificateData, cert.PrivateKeyData)
			if err != nil {
				return nil, err
			}

			return &x509Cert, nil
		},
		NextProtos: []string{"h2", "h2c", "http/1.1", "http/1.0", "spdy/3.1", "spdy/3"},
	})
	if err != nil {
		return nil, err
	}

	return &Server{Listener: listener, Handler: handler}, nil
}

func NewHTTPServer() (*Server, error) {
	handler := tools.NewHTTPHandler()

	httpPort := os.Getenv("HYVE_HTTP_PORT")
	if httpPort == "" {
		httpPort = "80"
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+httpPort)
	if err != nil {
		return nil, err
	}

	return &Server{Listener: listener, Handler: handler}, nil
}
