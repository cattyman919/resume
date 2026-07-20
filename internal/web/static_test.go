package web

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed static
var testStaticFiles embed.FS

func TestStaticFilesServed(t *testing.T) {
	staticFS, err := fs.Sub(testStaticFiles, "static")
	if err != nil {
		t.Fatalf("fs.Sub error = %v", err)
	}
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux := http.NewServeMux()
	mux.Handle("/static/", fileServer)

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/static/styles.css", http.StatusOK},
		{"/static/app.js", http.StatusOK},
		{"/static/editor.js", http.StatusOK},
		{"/static/split-panel.js", http.StatusOK},
		{"/static/nonexistent.css", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("GET %s: status = %d, want %d", tt.path, w.Code, tt.wantStatus)
			}
		})
	}
}

func TestStylesCSSDesignTokens(t *testing.T) {
	staticFS, err := fs.Sub(testStaticFiles, "static")
	if err != nil {
		t.Fatalf("fs.Sub error = %v", err)
	}
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux := http.NewServeMux()
	mux.Handle("/static/", fileServer)

	req := httptest.NewRequest(http.MethodGet, "/static/styles.css", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	wantTokens := []string{
		"--color-primary",
		"--color-accent",
		"--color-surface",
		"--color-background",
		".dark {",
		"--color-foreground",
		"--font-family",
		"prefers-reduced-motion",
		"--space-",
		"--color-border",
		"form-input",
		"section-card",
		"modal-backdrop",
		"accordion-icon",
	}

	for _, token := range wantTokens {
		if !strings.Contains(body, token) {
			t.Errorf("styles.css missing token: %q", token)
		}
	}
}

func TestAppJSDarkMode(t *testing.T) {
	staticFS, err := fs.Sub(testStaticFiles, "static")
	if err != nil {
		t.Fatalf("fs.Sub error = %v", err)
	}
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux := http.NewServeMux()
	mux.Handle("/static/", fileServer)

	req := httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	wantTokens := []string{
		"autocv-theme",
		"toggleTheme",
		"dark",
		"localStorage",
		"prefers-color-scheme",
	}

	for _, token := range wantTokens {
		if !strings.Contains(body, token) {
			t.Errorf("app.js missing token: %q", token)
		}
	}
}

func TestEditorJSFunctions(t *testing.T) {
	staticFS, err := fs.Sub(testStaticFiles, "static")
	if err != nil {
		t.Fatalf("fs.Sub error = %v", err)
	}
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux := http.NewServeMux()
	mux.Handle("/static/", fileServer)

	req := httptest.NewRequest(http.MethodGet, "/static/editor.js", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	wantTokens := []string{
		"generatePDF",
		"loadPdf",
		"toggleAccordion",
		"initSortable",
		"initAccordions",
		"openDialog",
		"closeDialog",
		"Escape",
		"zoomIn",
		"zoomOut",
		"zoomFit",
		"currentZoom",
		"updatePageInfo",
	}

	for _, token := range wantTokens {
		if !strings.Contains(body, token) {
			t.Errorf("editor.js missing token: %q", token)
		}
	}
}

func TestSplitPanelJS(t *testing.T) {
	staticFS, err := fs.Sub(testStaticFiles, "static")
	if err != nil {
		t.Fatalf("fs.Sub error = %v", err)
	}
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux := http.NewServeMux()
	mux.Handle("/static/", fileServer)

	req := httptest.NewRequest(http.MethodGet, "/static/split-panel.js", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	wantTokens := []string{
		"autocv-panel-sizes",
		"panel-divider",
		"mobile-tab",
		"editor-panel",
		"pdf-panel",
		"localStorage",
		"ArrowLeft",
		"ArrowRight",
	}

	for _, token := range wantTokens {
		if !strings.Contains(body, token) {
			t.Errorf("split-panel.js missing token: %q", token)
		}
	}
}

func TestDarkModeCSSOverrides(t *testing.T) {
	staticFS, err := fs.Sub(testStaticFiles, "static")
	if err != nil {
		t.Fatalf("fs.Sub error = %v", err)
	}
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux := http.NewServeMux()
	mux.Handle("/static/", fileServer)

	req := httptest.NewRequest(http.MethodGet, "/static/styles.css", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()

	darkTokens := []string{
		"--color-primary: #E2E8F0",
		"--color-accent: #38BDF8",
		"--color-background: #0F172A",
		"--color-surface: #1E293B",
		"--color-foreground: #F8FAFC",
		"panel-divider",
		"mobile-tab",
		"pdf-controls",
	}

	for _, token := range darkTokens {
		if !strings.Contains(body, token) {
			t.Errorf("styles.css .dark block missing: %q", token)
		}
	}
}
