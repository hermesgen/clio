package ssg

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func (h *WebHandler) NewSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New section form")
	form := NewSectionForm(r)
	h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "", http.StatusOK)
}

func (h *WebHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create section")

	form, err := SectionFormFromRequest(r)
	if err != nil {
		h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		section := ToFeatSection(form)
		webSection := ToWebSection(section)
		h.renderSectionForm(w, r, form, webSection, "Validation failed", http.StatusBadRequest)
		return
	}

	featSection := ToFeatSection(form)

	var response struct {
		Section feat.Section `json:"section"`
	}
	err = h.apiClient.Post(r, "/ssg/sections", featSection, &response)
	if err != nil {
		h.Err(w, err, "Failed to create section via API", http.StatusInternalServerError)
		return
	}
	createdSection := ToWebSection(response.Section)

	h.FlashInfo(w, r, "Section created")
	h.Redir(w, r, hm.EditPath(&Section{}, createdSection.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) EditSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit section")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing section ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Section feat.Section `json:"section"`
	}
	path := fmt.Sprintf("/ssg/sections/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get section from API", http.StatusInternalServerError)
		return
	}
	section := response.Section

	form := ToSectionForm(r, section)
	h.renderSectionForm(w, r, form, ToWebSection(section), "", http.StatusOK)
}

func (h *WebHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update section")

	form, err := SectionFormFromRequest(r)
	if err != nil {
		h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		section := ToFeatSection(form)
		webSection := ToWebSection(section)
		h.renderSectionForm(w, r, form, webSection, "Validation failed", http.StatusBadRequest)
		return
	}

	featSection := ToFeatSection(form)

	path := fmt.Sprintf("/ssg/sections/%s", featSection.GetID())
	err = h.apiClient.Put(r, path, featSection, nil)
	if err != nil {
		h.Err(w, err, "Failed to update section via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Section updated successfully")
	h.Redir(w, r, hm.ListPath(&Section{}), http.StatusSeeOther)
}

func (h *WebHandler) ListSections(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List sections")

	var response struct {
		Sections []feat.Section `json:"sections"`
	}
	err := h.apiClient.Get(r, "/ssg/sections", &response)
	if err != nil {
		h.Err(w, err, "Cannot get sections from API", http.StatusInternalServerError)
		return
	}
	sections := ToWebSections(response.Sections)

	page := hm.NewPage(r, sections)
	page.Form.SetAction(ssgPath)

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-sections")
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

func (h *WebHandler) ShowSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show section")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing section ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Section feat.Section `json:"section"`
	}
	path := fmt.Sprintf("/ssg/sections/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get section from API", http.StatusInternalServerError)
		return
	}
	section := ToWebSection(response.Section)

	page := hm.NewPage(r, section)
	page.Name = "Show Section"

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-section")
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

func (h *WebHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete section")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing section ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/sections/%s", idStr)
	err := h.apiClient.Delete(r, path)
	if err != nil {
		h.Err(w, err, "Failed to delete section via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Section deleted successfully")
	h.Redir(w, r, hm.ListPath(&Section{}), http.StatusSeeOther)
}

func (h *WebHandler) renderSectionForm(w http.ResponseWriter, r *http.Request, form SectionForm, section Section, errorMessage string, statusCode int) {
	var response struct {
		Layouts []Layout `json:"layouts"`
	}
	err := h.apiClient.Get(r, "/ssg/layouts", &response)
	if err != nil {
		h.Err(w, err, "Cannot get layouts from API", http.StatusInternalServerError)
		return
	}
	layouts := response.Layouts

	page := hm.NewPage(r, section)
	page.SetForm(&form)
	page.AddSelect("layouts", hm.ToSelectOpt(hm.ToPtrSlice(layouts)))

	if section.IsZero() {
		page.Name = "New Section"
		page.IsNew = true
		page.Form.SetAction(hm.CreatePath(&Section{}))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Section"
		page.IsNew = false
		page.Form.SetAction(hm.UpdatePath(&Section{}))
		page.Form.SetSubmitButtonText("Update")
	}


	tmpl, err := h.Tmpl().Get(ssgFeat, "new-section")
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
