package websocket

import (
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "basic handler",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			handler, err := NewHandler()

			if err != nil {
				t.Fatalf("expected no error, got: %s", err.Error())
			}

			if handler == nil {
				t.Fatalf("expected handler, got nil")
			}
		})
	}
}
