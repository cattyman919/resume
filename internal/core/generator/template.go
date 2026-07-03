package generator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var boldRegex = regexp.MustCompile(`\*\*(.*?)\*\*`)

// --- Helper Functions ---
var typstEscaper = strings.NewReplacer(
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
)

func escapeTypst(text string) string {
	// replace all **<text>** with \textbf{text}
	return boldRegex.ReplaceAllString(typstEscaper.Replace(text), `\textbf{$1}`)
}

func NewTemplate() (*template.Template, error) {
	filenames := []string{
		filepath.Join("src", "main.typ"),
		filepath.Join("src", "sections.typ"),
	}

	// --- Prepare Template ---
	// Create a FuncMap to register our helper function
	funcMap := template.FuncMap{
		"typst": escapeTypst,
	}

	// Initialize template with:
	// 1. The FuncMap
	// 2. Custom Delimiters << >> (to avoid collision with LaTeX {})
	tmpl, err := template.New("root").
		Funcs(funcMap).
		Delims("<<", ">>").
		ParseFiles(filenames...)

	if err != nil {
		return nil, fmt.Errorf("Error parsing template: %v", err)
	}

	return tmpl, nil
}
