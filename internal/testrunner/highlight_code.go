package testrunner

import (
	"log/slog"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// HighlightCode highlights code for a given language.
func HighlightCode(language string, code string) string {
	lexer := lexers.Get(language)

	if lexer == nil {
		lexer = lexers.Fallback
	}

	baseStyle := styles.Get("monokai")
	style, err := baseStyle.Builder().Add(chroma.Background, "#0000").Build()

	if err != nil {
		slog.Error(err.Error())

		return code
	}

	formatter := html.New(
		html.WithClasses(false),
		html.WithLineNumbers(false),
		html.InlineCode(true),
		html.PreventSurroundingPre(true),
	)

	iterator, err := lexer.Tokenise(nil, code)

	if err != nil {
		slog.Error(err.Error())

		return code
	}

	var highlighted strings.Builder
	err = formatter.Format(&highlighted, style, iterator)

	if err != nil {
		slog.Error(err.Error())

		return code
	}

	return highlighted.String()
}
