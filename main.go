package main

import (
	"log"

	"github.com/jonasroussel/proxbee/acme"
	"github.com/jonasroussel/proxbee/servers"
	"github.com/jonasroussel/proxbee/stores"
	"github.com/jonasroussel/proxbee/tools"
)

func main() {
	// Load environment variables
	tools.LoadEnv()

	// Load store
	stores.Load()

	// Load or create Let's Encrypt user
	err := acme.LoadOrCreateUser()
	if err != nil {
		log.Fatal(err)
	}

	// Create TLS server
	tlsListener, tlsServer, tlsHandler := servers.NewTLS()

	// Create HTTP server
	httpListener, httpServer, httpHandler := servers.NewHTTP()

	// Add the admin api to the TLS handler
	servers.AdminAPI(httpHandler)

	// Add the reverse proxy to the TLS handler
	servers.ReverseProxy(tlsHandler)

	// Add the HTTP-01 challenge solver to the HTTP handler
	servers.HTTP01ChallengeSolver(httpHandler)

	// Start TLS server
	go tlsServer.Serve(tlsListener)

	// Start HTTP server
	go httpServer.Serve(httpListener)

	// Wait for shutdown
	select {}
}
