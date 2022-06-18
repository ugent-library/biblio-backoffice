package services

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Service interface {
	Name() string
	Start() error
	Stop(context.Context) error
}

func Start(services ...Service) (err error) {
	var (
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)

	// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
	go func() {
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

		<-stopChan

		timer, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(len(services))
		for _, service := range services {
			service := service
			go func() {
				defer wg.Done()
				if err := service.Stop(timer); err != nil {
					errChan <- err
				}
			}()
		}
		wg.Wait()

		errChan <- nil
	}()

	// Start the services
	for _, service := range services {
		service := service
		go func() {
			if err := service.Start(); err != nil {
				errChan <- err
			}
		}()
	}

	return <-errChan
}
