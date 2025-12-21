package ssg

import (
	"bytes"
	"fmt"
	"net/http"

	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func (h *WebHandler) NewLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New layout form")
	form := NewLayoutForm(r)
	h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "", http.StatusOK)
}

func (h *WebHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create layout")

	form, err := LayoutFormFromRequest(r)
	if err != nil {
		h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		layout := ToFeatLayout(form)
		webLayout := ToWebLayout(layout)
		h.renderLayoutForm(w, r, form, webLayout, "Validation failed", http.StatusBadRequest)
		return
	}

	featLayout := ToFeatLayout(form)

	var response struct {
		Layout feat.Layout `json:"layout"`
	}
	err = h.apiClient.Post(h.addSiteSlugHeader(r), "/ssg/layouts", featLayout, &response)
	if err != nil {
		h.Err(w, err, "Failed to create layout via API", http.StatusInternalServerError)
		return
	}
	createdLayout := ToWebLayout(response.Layout)

	h.FlashInfo(w, r, "Layout created")
	h.Redir(w, r, hm.EditPath(&Layout{}, createdLayout.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) EditLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit layout")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing layout ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Layout feat.Layout `json:"layout"`
	}
	path := fmt.Sprintf("/ssg/layouts/%s", idStr)
	err := h.apiClient.Get(h.addSiteSlugHeader(r), path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get layout from API", http.StatusInternalServerError)
		return
	}
	layout := response.Layout

	form := ToLayoutForm(r, layout)
	h.renderLayoutForm(w, r, form, ToWebLayout(layout), "", http.StatusOK)
}

func (h *WebHandler) UpdateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update layout")

	form, err := LayoutFormFromRequest(r)
	if err != nil {
		h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		layout := ToFeatLayout(form)
		webLayout := ToWebLayout(layout)
		h.renderLayoutForm(w, r, form, webLayout, "Validation failed", http.StatusBadRequest)
		return
	}

	featLayout := ToFeatLayout(form)

	path := fmt.Sprintf("/ssg/layouts/%s", featLayout.GetID())
	err = h.apiClient.Put(h.addSiteSlugHeader(r), path, featLayout, nil)
	if err != nil {
		h.Err(w, err, "Failed to update layout via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Layout updated successfully")
	h.Redir(w, r, hm.ListPath(&Layout{}), http.StatusSeeOther)
}

func (h *WebHandler) ListLayouts(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List layouts")

	var response struct {
		Layouts []feat.Layout `json:"layouts"`
	}
	err := h.apiClient.Get(h.addSiteSlugHeader(r), "/ssg/layouts", &response)
	if err != nil {
		h.Err(w, err, "Cannot get layouts from API", http.StatusInternalServerError)
		return
	}
	layouts := ToWebLayouts(response.Layouts)

	page := hm.NewPage(r, layouts)
	page.Form.SetAction(ssgPath)

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-layouts")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *WebHandler) ShowLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show layout")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing layout ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Layout feat.Layout `json:"layout"`
	}
	path := fmt.Sprintf("/ssg/layouts/%s", idStr)
	err := h.apiClient.Get(h.addSiteSlugHeader(r), path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get layout from API", http.StatusInternalServerError)
		return
	}

	layout := ToWebLayout(response.Layout)

	page := hm.NewPage(r, layout)
	page.Name = "Show Layout"

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-layout")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, http.StatusOK)
}

func (h *WebHandler) DeleteLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete layout")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing layout ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/layouts/%s", idStr)
	err := h.apiClient.Delete(h.addSiteSlugHeader(r), path)
	if err != nil {
		h.Err(w, err, "Failed to delete layout via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Layout deleted successfully")
	h.Redir(w, r, hm.ListPath(&Layout{}), http.StatusSeeOther)
}

func (h *WebHandler) renderLayoutForm(w http.ResponseWriter, r *http.Request, form LayoutForm, layout Layout, errorMessage string, statusCode int) {
	page := hm.NewPage(r, layout)
	page.SetForm(&form)

	if layout.IsZero() {
		page.Name = "New Layout"
		page.IsNew = true
		page.Form.SetAction(hm.CreatePath(&Layout{}))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Layout"
		page.IsNew = false
		page.Form.SetAction(hm.UpdatePath(&Layout{}))
		page.Form.SetSubmitButtonText("Update")
	}

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-layout")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	page.SetFlash(h.GetFlash(r))

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, statusCode)
}
