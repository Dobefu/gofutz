// Package assets provides the static assets.
package assets

import (
	"embed"
	"net/http"
)

//go:embed css/*.css js/*.js
var assetFS embed.FS

// GetFS gets an asset file.
func GetFS() http.Handler {
	return http.FileServerFS(assetFS)
}
