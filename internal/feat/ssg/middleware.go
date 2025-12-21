package ssg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

// Context keys
type contextKey string

const (
	siteSlugKey      = contextKey("siteSlug")
	siteRepoKey      = contextKey("siteRepo")
	lastSiteCookie   = "last_site"
	lastSiteMaxAge   = 3600 * 24 * 365 // 1 year
)

// SiteContextMw is middleware that extracts site slug from session and injects appropriate repo.
type SiteContextMw struct {
	hm.Core
	sessionMgr  *auth.SessionManager
	siteRepo    SiteRepo
	repoManager *RepoManager
}

// NewSiteContextMw creates a new site context middleware.
func NewSiteContextMw(sessionMgr *auth.SessionManager, siteRepo SiteRepo, repoManager *RepoManager, params hm.XParams) *SiteContextMw {
	return &SiteContextMw{
		Core:        hm.NewCore("site-context-mw", params),
		sessionMgr:  sessionMgr,
		siteRepo:    siteRepo,
		repoManager: repoManager,
	}
}

func (mw *SiteContextMw) WebHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		_, err := mw.siteRepo.GetSiteBySlug(ctx, siteSlug)
		if err != nil {
			mw.Log().Error("Site not found", "slug", siteSlug)
			next.ServeHTTP(w, r)
			return
		}

		repo, err := mw.repoManager.GetRepoForSite(ctx, siteSlug)
		if err != nil {
			mw.Log().Error("Failed to get repo for site", "slug", siteSlug, "error", err)
			http.Error(w, "Failed to access site database", http.StatusInternalServerError)
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

		// Add site slug and repo to context
		ctx = context.WithValue(ctx, siteSlugKey, siteSlug)
		ctx = context.WithValue(ctx, siteRepoKey, repo)
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

		_, err := mw.siteRepo.GetSiteBySlug(ctx, siteSlug)
		if err != nil {
			mw.Log().Error("Site not found", "slug", siteSlug)
			http.Error(w, "Site not found", http.StatusNotFound)
			return
		}

		repo, err := mw.repoManager.GetRepoForSite(ctx, siteSlug)
		if err != nil {
			mw.Log().Error("Failed to get repo for site", "slug", siteSlug, "error", err)
			http.Error(w, "Failed to access site database", http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, siteSlugKey, siteSlug)
		ctx = context.WithValue(ctx, siteRepoKey, repo)
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

// GetRepoFromContext retrieves the site-specific repository from request context.
func GetRepoFromContext(ctx context.Context) (Repo, bool) {
	repo, ok := ctx.Value(siteRepoKey).(Repo)
	return repo, ok
}

// RequireSiteSlug is a helper to get site slug or return error.
func RequireSiteSlug(ctx context.Context) (string, error) {
	slug, ok := GetSiteSlugFromContext(ctx)
	if !ok || slug == "" {
		return "", fmt.Errorf("no site selected")
	}
	return slug, nil
}

// RequireRepo is a helper to get repo from context or return error.
func RequireRepo(ctx context.Context) (Repo, error) {
	repo, ok := GetRepoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no site repository in context")
	}
	return repo, nil
}
