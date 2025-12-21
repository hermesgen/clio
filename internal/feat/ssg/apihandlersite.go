package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

func (h *APIHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateSite", h.Name())

	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
		Mode string `json:"mode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	if req.Name == "" || req.Slug == "" || req.Mode == "" {
		h.Err(w, http.StatusBadRequest, "name, slug, and mode are required", nil)
		return
	}

	userID := uuid.New()
	site, err := h.siteManager.CreateSite(r.Context(), req.Name, req.Slug, req.Mode, userID)
	if err != nil {
		msg := fmt.Sprintf("Cannot create site: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "Site created successfully"
	h.Created(w, msg, site)
}

func (h *APIHandler) ListSites(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListSites", h.Name())

	sites, err := h.siteManager.ListSites(r.Context(), true)
	if err != nil {
		msg := "Cannot list sites"
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	response := map[string]interface{}{
		"sites": sites,
	}
	h.OK(w, "Sites retrieved successfully", response)
}
