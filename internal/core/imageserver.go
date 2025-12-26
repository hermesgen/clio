package core

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

type AdminFileServer struct {
	hm.Core
}

func NewAdminFileServer(params hm.XParams) *AdminFileServer {
	core := hm.NewCore("admin-file-server", params)
	return &AdminFileServer{
		Core: core,
	}
}

func (s *AdminFileServer) Setup(ctx context.Context) error {
	return nil
}

func (s *AdminFileServer) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteSlug, ok := ssg.GetSiteSlugFromContext(r.Context())
		if !ok || siteSlug == "" {
			siteSlug = "structured"
		}

		sitesBasePath := s.Cfg().StrValOrDef(ssg.SSGKey.SitesBasePath, "_workspace/sites")
		imagesPath := ssg.GetSiteImagesPath(sitesBasePath, siteSlug)

		requestPath := strings.TrimPrefix(r.URL.Path, "/static/images/")
		fullPath := filepath.Join(imagesPath, requestPath)

		http.ServeFile(w, r, fullPath)
	}
}
