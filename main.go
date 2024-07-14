package main

import (
	"github.com/jonasroussel/proxbee/config"
	"github.com/jonasroussel/proxbee/server"
)

func main() {
	// Load config
	config.Load()

	// Load store
	config.STORE.Load()

	// Create TLS server
	tlsListener, tlsServer, tlsHandler := server.NewTLS()

	// Create HTTP server
	httpListener, httpServer, httpHandler := server.NewHTTP()

	// Add the admin api to the TLS handler
	server.AdminAPI(tlsHandler)

	// Add the foward proxy to the TLS handler
	server.ForwardProxy(tlsHandler)

	// Add the HTTP-01 challenge to the HTTP handler
	server.HTTP01Challenge(httpHandler)

	// Start TLS server
	go tlsServer.Serve(tlsListener)

	// Start HTTP server
	go httpServer.Serve(httpListener)

	// Wait for shutdown
	select {}
}
