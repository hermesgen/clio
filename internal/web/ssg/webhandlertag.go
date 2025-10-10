package ssg

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func (h *WebHandler) NewTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New tag form")
	form := NewTagForm(r)
	h.renderTagForm(w, r, form, NewTag(""), "", http.StatusOK)
}

func (h *WebHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create tag")

	form, err := TagFormFromRequest(r)
	if err != nil {
		h.renderTagForm(w, r, form, NewTag(""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		tag := ToFeatTag(form)
		webTag := ToWebTag(tag)
		h.renderTagForm(w, r, form, webTag, "Validation failed", http.StatusBadRequest)
		return
	}

	featTag := ToFeatTag(form)

	var response struct {
		Tag feat.Tag `json:"tag"`
	}
	err = h.apiClient.Post(r, "/ssg/tags", featTag, &response)
	if err != nil {
		h.Err(w, err, "Failed to create tag via API", http.StatusInternalServerError)
		return
	}
	createdTag := ToWebTag(response.Tag)

	if hm.IsHTMXRequest(r) {
		redirectURL := hm.EditPath(&createdTag, createdTag.GetID())
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.FlashInfo(w, r, "Tag created")
	h.Redir(w, r, hm.EditPath(&createdTag, createdTag.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) EditTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit tag")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing tag ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Tag feat.Tag `json:"tag"`
	}
	path := fmt.Sprintf("/ssg/tags/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get tag from API", http.StatusInternalServerError)
		return
	}
	webTag := ToWebTag(response.Tag)

	form := ToTagForm(r, response.Tag)
	h.renderTagForm(w, r, form, webTag, "", http.StatusOK)
}

func (h *WebHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update tag")

	form, err := TagFormFromRequest(r)
	if err != nil {
		h.renderTagForm(w, r, form, NewTag(""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		tag := ToFeatTag(form)
		webTag := ToWebTag(tag)
		h.renderTagForm(w, r, form, webTag, "Validation failed", http.StatusBadRequest)
		return
	}

	featTag := ToFeatTag(form)

	path := fmt.Sprintf("/ssg/tags/%s", featTag.GetID())
	err = h.apiClient.Put(r, path, featTag, nil)
	if err != nil {
		h.Err(w, err, "Failed to update tag via API", http.StatusInternalServerError)
		return
	}

	if hm.IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<div id=\"save-status\" data-timestamp=\"" + time.Now().Format(hm.TimeFormat) + "\"></div>"))
		return
	}

	h.FlashInfo(w, r, "Tag updated successfully")
	webTag := ToWebTag(featTag)
	h.Redir(w, r, hm.EditPath(&webTag, webTag.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List tags")

	var response struct {
		Tags []feat.Tag `json:"tags"`
	}
	err := h.apiClient.Get(r, "/ssg/tags", &response)
	if err != nil {
		h.Err(w, err, "Cannot get tags from API", http.StatusInternalServerError)
		return
	}
	webTags := ToWebTags(response.Tags)

	page := hm.NewPage(r, webTags)
	page.Form.SetAction(ssgPath)

	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(&Tag{})

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-tags")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *WebHandler) ShowTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show tag")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing tag ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Tag feat.Tag `json:"tag"`
	}
	path := fmt.Sprintf("/ssg/tags/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get tag from API", http.StatusInternalServerError)
		return
	}
	tag := ToWebTag(response.Tag)

	page := hm.NewPage(r, tag)
	page.Name = "Show Tag"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&tag, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-tag")
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

func (h *WebHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete tag")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing tag ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/tags/%s", idStr)
	err := h.apiClient.Delete(r, path)
	if err != nil {
		h.Err(w, err, "Failed to delete tag via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Tag deleted successfully")
	h.Redir(w, r, hm.ListPath(&Tag{}), http.StatusSeeOther)
}

func (h *WebHandler) renderTagForm(w http.ResponseWriter, r *http.Request, form TagForm, tag Tag, errorMessage string, statusCode int) {
	page := hm.NewPage(r, tag)
	page.SetForm(&form)

	if tag.IsZero() {
		page.Name = "New Tag"
		page.IsNew = true
		page.Form.SetAction(hm.CreatePath(&Tag{}))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Tag"
		page.IsNew = false
		page.Form.SetAction(hm.UpdatePath(&Tag{}))
		page.Form.SetSubmitButtonText("Update")
	}

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&tag, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-tag")
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
