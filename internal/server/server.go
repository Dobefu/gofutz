// Package server provides server functionality.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Server defines a server.
type Server struct {
	host       string
	port       int
	httpServer *http.Server
}

// NewServer creates a new server.
func NewServer(host string, port int) *Server {
	mux := http.NewServeMux()

	handleRoutes(mux)

	server := &Server{
		host: host,
		port: port,
		httpServer: &http.Server{ // nolint:exhaustruct
			Addr:              fmt.Sprintf("%s:%d", host, port),
			Handler:           mux,
			ReadTimeout:       time.Second * 15,
			ReadHeaderTimeout: time.Second * 15,
			WriteTimeout:      time.Second * 15,
			IdleTimeout:       time.Second * 15,
		},
	}

	return server
}

// Start starts the server.
func (s *Server) Start() error {
	url := fmt.Sprintf("http://%s:%d", s.host, s.port)
	slog.Info(fmt.Sprintf("Starting server on %s", url))

	return s.httpServer.ListenAndServe()
}
