// Package api gère tout ce qui se rapport à l'api
package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	http.Server
}

func NewServer(port, address string, api *API) *Server {
	handler := api.Routes()
	server := Server{
		Server: http.Server{
			Addr:    fmt.Sprintf("%s:%s", address, port),
			Handler: handler,
		},
	}
	return &server
}

func (s *Server) Start() error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.ListenAndServe()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case err := <-errChan:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Println("Exctinction du server")
		err := s.Shutdown(ctx)
		return err
	}
	return nil
}
