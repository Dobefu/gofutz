package server

import (
	"net/http"
	"testing"
)

func TestHandleRoutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "handle routes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			handleRoutes(mux)

			if mux == nil {
				t.Fatalf("expected mux, got nil")
			}
		})
	}
}
