// Package server provides server functionality.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Start starts the server.
func Start(host string, port int) error {
	url := fmt.Sprintf("http://%s:%d", host, port)
	slog.Info(fmt.Sprintf("Starting server on %s", url))

	mux := http.NewServeMux()

	handleRoutes(mux)

	server := http.Server{ // nolint:exhaustruct
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           mux,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 15,
	}

	err := server.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
