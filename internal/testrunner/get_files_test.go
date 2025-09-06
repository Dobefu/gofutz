package testrunner

import (
	"testing"
)

func TestGetFiles(t *testing.T) {
	t.Parallel()

	runner := &TestRunner{ // nolint:exhaustruct
		files: map[string]File{
			"test.go": { // nolint:exhaustruct
				Name: "test.go",
				Functions: []Function{
					{ // nolint:exhaustruct
						Name: "SomeFunction",
					},
				},
			},
		},
	}

	files := runner.GetFiles()

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got: %d", len(files))
	}

	if files["test.go"].Name != "test.go" {
		t.Fatalf("expected file name to be test.go, got: %s", files["test.go"].Name)
	}

	if len(files["test.go"].Functions) != 1 {
		t.Fatalf("expected 1 function, got: %d", len(files["test.go"].Functions))
	}

	if files["test.go"].Functions[0].Name != "SomeFunction" {
		t.Fatalf(
			"expected function name to be \"SomeFunction\", got: \"%s\"",
			files["test.go"].Functions[0].Name,
		)
	}
}
