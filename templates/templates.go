// Package templates provides the HTML templates.
package templates

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed layouts/*.gohtml pages/*.gohtml components/**/*.gohtml
var templateFS embed.FS

// GetPageTemplate gets a page template, along with its dependencies.
func GetPageTemplate(page string) (*template.Template, error) {
	tmpl, err := template.ParseFS(
		templateFS,
		fmt.Sprintf("pages/%s.gohtml", page),
		"layouts/default.gohtml",
		"components/**/*.gohtml",
	)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
