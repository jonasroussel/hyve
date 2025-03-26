package hyve

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jonasroussel/hyve/core"
)

var _ core.App = (*Hyve)(nil)

type Hyve struct {
	core.App
}

func New() *Hyve {
	hyve := &Hyve{App: core.NewApp()}

	return hyve
}

func (hyve *Hyve) Start() error {
	if err := hyve.Bootstrap(); err != nil {
		return err
	}

	done := make(chan error, 1)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch

		done <- nil
	}()

	// serve all servers
	go func() {
		if err := hyve.ServeAll(); err != nil {
			done <- err
			return
		}

		done <- nil
	}()

	if err := <-done; err != nil {
		return err
	}

	return nil
}
