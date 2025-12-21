package ssg

import (
	"bytes"
	"net/http"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

// NOTE: Site handlers use SiteManager directly instead of API client.
// TODO: Add API endpoints for sites to support API-only mode (e.g., Neovim plugin).

func (wh *WebHandler) ListSites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sites, err := wh.siteManager.ListSites(ctx, true)
	if err != nil {
		wh.Log().Error("Failed to get sites", "error", err)
		wh.Err(w, err, "Cannot get sites", http.StatusInternalServerError)
		return
	}

	page := hm.NewPage(r, sites)
	page.Form.SetAction(ssgPath)

	tmpl, err := wh.Tmpl().Get(ssgFeat, "list-sites")
	if err != nil {
		wh.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		wh.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (wh *WebHandler) NewSite(w http.ResponseWriter, r *http.Request) {
	page := hm.NewPage(r, nil)
	page.Form.SetAction(ssgPath)

	tmpl, err := wh.Tmpl().Get(ssgFeat, "new-site")
	if err != nil {
		wh.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		wh.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (wh *WebHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		wh.FlashError(w, r, "Invalid form data")
		http.Redirect(w, r, "/ssg/sites/new", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	slug := r.FormValue("slug")
	mode := r.FormValue("mode")

	if name == "" || slug == "" || mode == "" {
		wh.FlashError(w, r, "Name, slug, and mode are required")
		http.Redirect(w, r, "/ssg/sites/new", http.StatusSeeOther)
		return
	}

	userID := uuid.New()

	_, err := wh.siteManager.CreateSite(ctx, name, slug, mode, userID)
	if err != nil {
		wh.Log().Error("Failed to create site", "error", err)
		wh.FlashError(w, r, "Failed to create site: "+err.Error())
		http.Redirect(w, r, "/ssg/sites/new", http.StatusSeeOther)
		return
	}

	wh.FlashInfo(w, r, "Site created successfully")
	http.Redirect(w, r, "/ssg/sites", http.StatusSeeOther)
}

func (wh *WebHandler) SwitchSite(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		wh.FlashError(w, r, "Site slug is required")
		http.Redirect(w, r, "/ssg/sites", http.StatusSeeOther)
		return
	}

	userID := uuid.New()
	if err := wh.sessionManager.SetUserSession(w, userID, slug); err != nil {
		wh.Log().Error("Failed to set session", "error", err)
		wh.FlashError(w, r, "Failed to switch site")
		http.Redirect(w, r, "/ssg/sites", http.StatusSeeOther)
		return
	}

	wh.FlashInfo(w, r, "Switched to site: "+slug)
	http.Redirect(w, r, "/ssg/list-content?site="+slug, http.StatusSeeOther)
}

func (wh *WebHandler) RootRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var lastSite string

	if cookie, err := r.Cookie("last_site"); err == nil {
		lastSite = cookie.Value
	}

	if lastSite != "" {
		if _, err := wh.siteManager.GetSiteBySlug(ctx, lastSite); err == nil {
			http.Redirect(w, r, "/ssg/list-content?site="+lastSite, http.StatusFound)
			return
		}
	}

	http.Redirect(w, r, "/ssg/sites", http.StatusFound)
}
