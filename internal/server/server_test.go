package server

import (
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
