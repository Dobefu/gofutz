package assets

import (
	"testing"
)

func TestGetFS(t *testing.T) {
	t.Parallel()

	fs := GetFS()

	if fs == nil {
		t.Fatalf("expected FS http.Handler, got nil")
	}
}
