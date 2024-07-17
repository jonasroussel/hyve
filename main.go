package main

import (
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
	acme.LoadOrCreateUser()

	// Create TLS server
	tlsListener, tlsServer, tlsHandler := servers.NewTLS()

	// Create HTTP server
	httpListener, httpServer, httpHandler := servers.NewHTTP()

	// Add the admin api to the TLS handler
	servers.AdminAPI(tlsHandler)

	// Add the reverse proxy to the TLS handler
	servers.ReverseProxy(tlsHandler)

	// Add the HTTP-01 challenge solver to the HTTP handler
	servers.HTTP01ChallengeSolver(httpHandler)

	// Start TLS server
	go tlsServer.Serve(tlsListener)

	// Start HTTP server
	go httpServer.Serve(httpListener)

	// Register admin domain if needed
	acme.RegisterAdminDomain()

	// Wait for shutdown
	select {}
}
