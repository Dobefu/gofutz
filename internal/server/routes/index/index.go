// Package index provides an index route handler.
package index

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Dobefu/gofutz/internal/server/routes"
	"github.com/Dobefu/gofutz/templates"
)

// Handle handles the route.
func Handle(w http.ResponseWriter, _ *http.Request) {
	const tmplName = "index"

	tmpl, err := templates.GetPageTemplate(tmplName)

	if err != nil {
		slog.Error(fmt.Sprintf("Could not load template: %s", err.Error()))
		http.Error(w, "Internal server error", 500)

		return
	}

	err = tmpl.ExecuteTemplate(w, fmt.Sprintf("pages/%s", tmplName), routes.PageVars{
		Title: "Home",
	})

	if err != nil {
		slog.Error(fmt.Sprintf("Could not load template: %s", err.Error()))
		http.Error(w, "Internal server error", 500)

		return
	}
}
