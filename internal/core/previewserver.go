package core

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

// MultiSitePreviewHandler serves static HTML files for multiple sites.
// URL format: /{slug}/path/to/file.html
type MultiSitePreviewHandler struct {
	cfg *hm.Config
	log hm.Logger
}

func NewMultiSitePreviewHandler(params hm.XParams) *MultiSitePreviewHandler {
	return &MultiSitePreviewHandler{
		cfg: params.Cfg,
		log: params.Log,
	}
}

func (h *MultiSitePreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract site slug from subdomain: {slug}.localhost:8082
	host := r.Host
	var slug string

	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Check for subdomain pattern
	if strings.HasSuffix(host, ".localhost") {
		slug = strings.TrimSuffix(host, ".localhost")
	} else if host == "localhost" {
		slug = "default"
	} else {
		http.Error(w, "Invalid host. Use {slug}.localhost or localhost for default site", http.StatusBadRequest)
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "index.html"
	}

	sitesBasePath := h.cfg.StrValOrDef(ssg.SSGKey.SitesBasePath, "_workspace/sites")
	htmlPath := ssg.GetSiteHTMLPath(sitesBasePath, slug)
	fullPath := filepath.Join(htmlPath, filePath)

	cleanPath := filepath.Clean(fullPath)
	if !strings.HasPrefix(cleanPath, filepath.Clean(htmlPath)) {
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		cleanPath = filepath.Join(cleanPath, "index.html")
	}

	h.log.Debug("Serving file", "slug", slug, "path", filePath, "fullPath", cleanPath)
	http.ServeFile(w, r, cleanPath)
}
