package index

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "basic response",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			Handle(w, r)
			res := w.Result()
			defer func() { _ = r.Body.Close() }()

			body, err := io.ReadAll(res.Body)

			if err != nil {
				t.Fatalf("expected no error, got %s", err.Error())
			}

			if len(body) <= 100 {
				t.Errorf("expected a response length of at least 100 characters, got %d", len(body))
				t.Logf("body: %s", string(body))
			}
		})
	}
}
