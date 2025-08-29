package server

import (
	"net/http"

	"github.com/Dobefu/gofutz/internal/server/routes/index"
)

func handleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", index.Handle)
}
