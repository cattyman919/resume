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
var latexEscaper = strings.NewReplacer(
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
	`\`, `\textbackslash{}`,
)

func escapeLatex(text string) string {
	// replace all **<text>** with \textbf{text}
	return boldRegex.ReplaceAllString(latexEscaper.Replace(text), `\textbf{$1}`)
}

func GetTemplate() (*template.Template, error) {

	filenames := []string{
		filepath.Join("template_cv", "bw_cv", "sections", "Achivements_Skills.tex"),
		filepath.Join("template_cv", "bw_cv", "sections", "Research_Interests.tex"),
		filepath.Join("template_cv", "bw_cv", "sections", "Awards.tex"),
		filepath.Join("template_cv", "bw_cv", "sections", "Education.tex"),
		filepath.Join("template_cv", "bw_cv", "sections", "Experience.tex"),
		filepath.Join("template_cv", "bw_cv", "sections", "Header.tex"),
		filepath.Join("template_cv", "bw_cv", "sections", "Projects.tex"),
	}

	// --- Prepare Template ---
	// Create a FuncMap to register our helper function
	funcMap := template.FuncMap{
		"latex": escapeLatex,
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
