package core

import (
	"github.com/jonasroussel/hyve/stores"
	"github.com/jonasroussel/hyve/ui"
	"gorm.io/gorm"
)

type App interface {
	Bootstrap() error
	ServeAll() error

	GetStore() (*DBStore, error)
	SetStore(DBStore) error
}

var _ App = (*HyveApp)(nil)

type HyveApp struct {
	DB    *gorm.DB
	Store stores.Store

	https *Server
	http  *Server
}

func NewApp() *HyveApp {
	return &HyveApp{}
}

func (app *HyveApp) Bootstrap() error {
	// Database
	db, err := OpenDatabase()
	if err != nil {
		return err
	}
	app.DB = db

	// Store
	store, err := app.GetStore()
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if store != nil {
		err := app.SetStore(*store)
		if err != nil {
			return err
		}
	}

	// HTTPS server
	app.https, err = NewHTTPSServer(app)
	if err != nil {
		return err
	}

	// HTTP server
	app.http, err = NewHTTPServer()
	if err != nil {
		return err
	}

	// Serve the API
	ServeAPI(app.http.Handler, app)

	// Serve the UI static files
	ui.ServeFiles(app.http.Handler)

	return nil
}

func (app *HyveApp) ServeAll() error {
	servers := []*Server{app.https, app.http}

	errChan := make(chan error, len(servers))

	for _, server := range servers {
		go func() {
			errChan <- server.Serve()
		}()
	}

	for range servers {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}
