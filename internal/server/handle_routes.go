package server

import (
	"net/http"

	"github.com/Dobefu/gofutz/assets"
	"github.com/Dobefu/gofutz/internal/server/routes/index"
	"github.com/Dobefu/gofutz/internal/server/routes/ws"
)

func handleRoutes(mux *http.ServeMux) {
	// Assets.
	mux.Handle("GET /", assets.GetFS())

	// Pages.
	mux.HandleFunc("GET /{$}", index.Handle)
	mux.HandleFunc("GET /ws", ws.Handle)
}
