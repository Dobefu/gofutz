package websocket

import (
	"github.com/Dobefu/gofutz/internal/testrunner"
)

// Params defines a set of message parameters.
type Params struct {
	Files    map[string]testrunner.File `json:"files"`
	Coverage float64                    `json:"coverage"`
}
