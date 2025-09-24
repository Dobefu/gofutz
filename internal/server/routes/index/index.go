// Package index provides an index route handler.
package index

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dobefu/gofutz/internal/server/routes"
	"github.com/Dobefu/gofutz/internal/testrunner"
	"github.com/Dobefu/gofutz/templates"
)

// Handle handles the route.
func Handle(w http.ResponseWriter, r *http.Request) {
	const tmplName = "index"

	tmpl, err := templates.GetPageTemplate(tmplName)

	if err != nil {
		slog.Error(fmt.Sprintf("Could not load template: %s", err.Error()))
		http.Error(w, "Internal server error", 500)

		return
	}

	title := "GoFutz"
	cwd, err := os.Getwd()

	if err == nil {
		moduleName := testrunner.GetModuleName(filepath.Join(cwd, "go.mod"))
		moduleParts := strings.Split(moduleName, "/")

		title = fmt.Sprintf("%s | GoFutz", moduleParts[len(moduleParts)-1])
	}

	sortOptions := []routes.SortOption{
		{Value: "name-asc", Label: "Name (A-Z)"},
		{Value: "name-desc", Label: "Name (Z-A)"},
		{Value: "coverage-asc", Label: "Coverage (Low-High)"},
		{Value: "coverage-desc", Label: "Coverage (High-Low)"},
	}

	selectedSort := r.URL.Query().Get("sort")

	if selectedSort == "" {
		selectedSort = "name-asc"
	}

	err = tmpl.ExecuteTemplate(
		w,
		fmt.Sprintf("pages/%s", tmplName),
		routes.PageVars{
			Title:              title,
			SortOptions:        sortOptions,
			SelectedSortOption: selectedSort,
		},
	)

	if err != nil {
		slog.Error(fmt.Sprintf("Could not load template: %s", err.Error()))
		http.Error(w, "Internal server error", 500)

		return
	}
}
