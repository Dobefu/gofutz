package testrunner

import (
	"fmt"
	"os"
)

// Test defines a single test.
type Test struct {
	Name string `json:"name"`
}

// Run runs a test.
func (t *Test) Run() {
	_, _ = fmt.Fprintf(os.Stdout, "Running test: %s\n", t.Name)
}
