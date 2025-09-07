package server

import (
	"context"
	"testing"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		host string
		port int
	}{
		{
			name: "default host and port",
			host: "127.0.0.1",
			port: 7357,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := NewServer(test.host, test.port)

			if server == nil {
				t.Fatal("expected server, got nil")
			}

			if server.httpServer == nil {
				t.Fatal("expected httpServer, got nil")
			}
		})
	}
}

func TestShutdown(t *testing.T) {
	t.Parallel()

	server := NewServer("127.0.0.1", 7357)

	if server == nil {
		t.Fatal("expected server, got nil")
	}

	err := server.Shutdown(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: \"%s\"", err.Error())
	}
}
