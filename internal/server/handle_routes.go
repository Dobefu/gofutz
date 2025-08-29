package server

import (
	"net/http"

	"github.com/Dobefu/gofutz/assets"
	"github.com/Dobefu/gofutz/internal/server/routes/index"
)

func handleRoutes(mux *http.ServeMux) {
	// Pages.
	mux.HandleFunc("GET /{$}", index.Handle)

	// Assets.
	mux.Handle("GET /", assets.GetFS())
}
