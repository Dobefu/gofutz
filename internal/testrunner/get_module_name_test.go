package testrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createGoModFile(
	t *testing.T,
	name string,
	content string,
) (string, func()) {
	t.Helper()

	if content == "" {
		return "", func() {}
	}

	modFile := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("get_module_name/%s/go.mod", filepath.Clean(name)),
	)

	err := os.MkdirAll(filepath.Dir(modFile), 0700)

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	err = os.WriteFile(modFile, []byte(content), 0600)

	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	return modFile, func() { _ = os.Remove(modFile) }
}

func TestGetModuleName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		modString string
		expected  string
	}{
		{
			name:      "valid go.mod",
			modString: "module github.com/example/project\n\nrequire github.com/example/dependency v1.0.0\n",
			expected:  "github.com/example/project",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			modFile, cleanup := createGoModFile(t, test.name, test.modString)
			defer cleanup()

			moduleName := GetModuleName(modFile)

			if moduleName != test.expected {
				t.Fatalf(
					"expected module name to be \"%s\", got \"%s\"",
					test.expected,
					moduleName,
				)
			}
		})
	}
}
