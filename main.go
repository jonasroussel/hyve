package main

import (
	"github.com/jonasroussel/hyve/acme"
	"github.com/jonasroussel/hyve/servers"
	"github.com/jonasroussel/hyve/stores"
	"github.com/jonasroussel/hyve/tools"
)

func main() {
	// Load environment variables
	tools.LoadEnv()

	// Load store
	stores.Load()

	// Load or create Let's Encrypt user
	acme.LoadOrCreateUser()

	// Init lego
	acme.InitLego()

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

	// Activate auto renew
	acme.ActivateAutoRenew()

	// Wait for shutdown
	select {}
}
