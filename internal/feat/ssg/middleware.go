package ssg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

// Context keys
type contextKey string

const (
	siteSlugKey    = contextKey("siteSlug")
	siteIDKey      = contextKey("siteID")
	lastSiteCookie = "last_site"
	lastSiteMaxAge = 3600 * 24 * 365 // 1 year
)

type SiteRepoProvider interface {
	GetSiteBySlug(ctx context.Context, slug string) (Site, error)
}

// SiteContextMw is middleware that extracts site slug from session and injects site context.
type SiteContextMw struct {
	hm.Core
	sessionMgr       *auth.SessionManager
	siteRepoProvider SiteRepoProvider
}

// NewSiteContextMw creates a new site context middleware.
func NewSiteContextMw(sessionMgr *auth.SessionManager, siteRepoProvider SiteRepoProvider, params hm.XParams) *SiteContextMw {
	return &SiteContextMw{
		Core:             hm.NewCore("site-context-mw", params),
		sessionMgr:       sessionMgr,
		siteRepoProvider: siteRepoProvider,
	}
}

func (mw *SiteContextMw) isExemptPath(path string) bool {
	exemptPaths := []string{
		"/ssg/sites",
		"/ssg/sites/new",
		"/ssg/sites/create",
		"/ssg/sites/switch",
		"/ssg/sites/delete",
	}

	for _, exempt := range exemptPaths {
		if path == exempt {
			return true
		}
	}

	return false
}

func (mw *SiteContextMw) WebHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mw.isExemptPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		var siteSlug string

		siteSlug = r.URL.Query().Get("site")

		if siteSlug == "" {
			if cookie, err := r.Cookie(lastSiteCookie); err == nil {
				siteSlug = cookie.Value
			}
		}

		if siteSlug == "" {
			http.Redirect(w, r, "/ssg/sites", http.StatusFound)
			return
		}

		site, err := mw.siteRepoProvider.GetSiteBySlug(ctx, siteSlug)
		if err != nil {
			mw.Log().Info("Site not found in database, clearing session", "slug", siteSlug)
			http.SetCookie(w, &http.Cookie{
				Name:     lastSiteCookie,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/ssg/sites", http.StatusFound)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     lastSiteCookie,
			Value:    siteSlug,
			Path:     "/",
			MaxAge:   lastSiteMaxAge,
			Expires:  time.Now().Add(time.Duration(lastSiteMaxAge) * time.Second),
			SameSite: http.SameSiteLaxMode,
		})

		// Add site slug and ID to context
		ctx = context.WithValue(ctx, siteSlugKey, siteSlug)
		ctx = context.WithValue(ctx, siteIDKey, site.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *SiteContextMw) APIHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		siteSlug := r.Header.Get("X-Site-Slug")

		if siteSlug == "" {
			http.Error(w, "X-Site-Slug header is required", http.StatusBadRequest)
			return
		}

		site, err := mw.siteRepoProvider.GetSiteBySlug(ctx, siteSlug)
		if err != nil {
			mw.Log().Error("Site not found", "slug", siteSlug)
			http.Error(w, "Site not found", http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, siteSlugKey, siteSlug)
		ctx = context.WithValue(ctx, siteIDKey, site.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *SiteContextMw) Handler(next http.Handler) http.Handler {
	return mw.WebHandler(next)
}

// GetSiteSlugFromContext retrieves site slug from request context.
func GetSiteSlugFromContext(ctx context.Context) (string, bool) {
	slug, ok := ctx.Value(siteSlugKey).(string)
	return slug, ok
}

// GetSiteIDFromContext retrieves site ID from request context.
func GetSiteIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(siteIDKey).(uuid.UUID)
	return id, ok
}

// RequireSiteSlug is a helper to get site slug or return error.
func RequireSiteSlug(ctx context.Context) (string, error) {
	slug, ok := GetSiteSlugFromContext(ctx)
	if !ok || slug == "" {
		return "", fmt.Errorf("no site selected")
	}
	return slug, nil
}

// RequireSiteID is a helper to get site ID from context or return error.
func RequireSiteID(ctx context.Context) (uuid.UUID, error) {
	id, ok := GetSiteIDFromContext(ctx)
	if !ok || id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("no site ID in context")
	}
	return id, nil
}
