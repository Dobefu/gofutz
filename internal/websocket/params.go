package websocket

import (
	"github.com/Dobefu/gofutz/internal/testrunner"
)

// InitParams defines a set of init message parameters.
type InitParams struct {
	Files     map[string]testrunner.File `json:"files"`
	Coverage  float64                    `json:"coverage"`
	IsRunning bool                       `json:"isRunning"`
	Output    []string                   `json:"output"`
}

// UpdateParams defines a set of update message parameters.
type UpdateParams struct {
	Files     map[string]testrunner.File `json:"files"`
	Coverage  float64                    `json:"coverage"`
	IsRunning bool                       `json:"isRunning"`
}

// OutputParams defines a set of output message parameters.
type OutputParams struct {
	Output []string `json:"output"`
}
