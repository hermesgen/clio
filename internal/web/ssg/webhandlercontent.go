package ssg

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	hm "github.com/hermesgen/hm"
	"github.com/hermesgen/clio/internal/feat/auth"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func (h *WebHandler) NewContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New content form")
	form := NewContentForm(r)
	h.renderContentForm(w, r, form, NewContent("", ""), "", http.StatusOK)
}

func (h *WebHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create content")

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		content := ToFeatContent(form)
		webContent := ToWebContent(content)
		h.renderContentForm(w, r, form, webContent, "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToFeatContent(form)

	var response struct {
		Content feat.Content `json:"content"`
	}
	err = h.apiClient.Post(r, "/ssg/contents", content, &response)
	if err != nil {
		h.Err(w, err, "Failed to create content via API", http.StatusInternalServerError)
		return
	}
	createdContent := ToWebContent(response.Content)

	if hm.IsHTMXRequest(r) {
		redirectURL := hm.EditPath(&createdContent, createdContent.GetID())
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.FlashInfo(w, r, "Content created")
	h.Redir(w, r, hm.EditPath(&createdContent, createdContent.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) EditContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit content")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing content ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Content feat.Content `json:"content"`
	}
	path := fmt.Sprintf("/ssg/contents/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get content from API", http.StatusInternalServerError)
		return
	}
	content := response.Content
	webContent := ToWebContent(content)

	form := ToContentForm(r, content)
	h.renderContentForm(w, r, form, webContent, "", http.StatusOK)
}

func (h *WebHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update content")

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		content := ToFeatContent(form)
		webContent := ToWebContent(content)
		h.renderContentForm(w, r, form, webContent, "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToFeatContent(form)

	path := fmt.Sprintf("/ssg/contents/%s", content.GetID())
	err = h.apiClient.Put(r, path, content, nil)
	if err != nil {
		h.Err(w, err, "Failed to update content via API", http.StatusInternalServerError)
		return
	}

	if hm.IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<div id=\"save-status\" data-timestamp=\"" + time.Now().Format(hm.TimeFormat) + "\"></div>"))
		return
	}

	h.FlashInfo(w, r, "Content updated successfully")
	h.Redir(w, r, hm.EditPath(&Content{}, content.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) ListContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List content")

	var response struct {
		Contents []feat.Content `json:"contents"`
	}
	h.Log().Info("Calling apiClient.Get /contents")
	h.Log().Infof("h.apiClient: %+v", h.apiClient)
	err := h.apiClient.Get(r, "/ssg/contents", &response)
	if err != nil {
		h.Err(w, err, "Cannot get contents from API", http.StatusInternalServerError)
		return
	}
	contents := response.Contents

	h.Log().Infof("Contents received: %+v", contents)
	page := hm.NewPage(r, contents)
	h.Log().Info("Page created")
	page.Form.SetAction(ssgPath)
	h.Log().Info("Form action set")

	menu := page.NewMenu(ssgPath)
	h.Log().Info("Menu created")
	menu.AddNewItem(&Content{})
	h.Log().Info("Menu item added")

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-content")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}
	h.Log().Info("Template retrieved")

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}
	h.Log().Info("Template executed")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *WebHandler) ShowContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show content")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing content ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Content feat.Content `json:"content"`
	}
	path := fmt.Sprintf("/ssg/contents/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get content from API", http.StatusInternalServerError)
		return
	}
	content := response.Content

	page := hm.NewPage(r, content)
	page.Name = "Show Content"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&content, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-content")
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

func (h *WebHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete content")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing content ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/contents/%s", idStr)
	err := h.apiClient.Delete(r, path)
	if err != nil {
		h.Err(w, err, "Failed to delete content via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Content deleted successfully")
	h.Redir(w, r, hm.ListPath(&Content{}), http.StatusSeeOther)
}

func (h *WebHandler) renderContentForm(w http.ResponseWriter, r *http.Request, form ContentForm, content Content, errorMessage string, statusCode int) {
	h.Log().Debugf("h.apiClient: %+v", h.apiClient)

	var sectionsResponse struct {
		Sections []Section `json:"sections"`
	}
	h.Log().Debug("Calling API to get sections")
	err := h.apiClient.Get(r, "/ssg/sections", &sectionsResponse)
	if err != nil {
		h.Log().Errorf("Cannot get sections from API: %v", err)
		h.Err(w, err, "Cannot get sections from API", http.StatusInternalServerError)
		return
	}
	sections := sectionsResponse.Sections
	h.Log().Debugf("Sections received: %+v", sections)

	var usersResponse struct {
		Users []auth.User `json:"users"`
	}
	h.Log().Debug("Calling API to get users")
	err = h.apiClient.Get(r, "/auth/users", &usersResponse)
	if err != nil {
		h.Log().Errorf("Cannot get users from API: %v", err)
		h.Err(w, err, "Cannot get users from API", http.StatusInternalServerError)
		return
	}
	users := usersResponse.Users
	h.Log().Debugf("Users received: %+v", users)

	var tagsResponse struct {
		Tags []Tag `json:"tags"`
	}
	h.Log().Debug("Calling API to get tags")
	err = h.apiClient.Get(r, "/ssg/tags", &tagsResponse)
	if err != nil {
		h.Log().Errorf("Cannot get tags from API: %v", err)
		h.Err(w, err, "Cannot get tags from API", http.StatusInternalServerError)
		return
	}
	tags := tagsResponse.Tags
	h.Log().Debugf("Tags received: %+v", tags)

	kinds := []hm.SelectOpt{
		{Value: "article", Label: "Article"},
		{Value: "page", Label: "Page"},
		{Value: "blog", Label: "Blog"},
		{Value: "series", Label: "Series"},
	}

	page := hm.NewPage(r, content)
	page.SetForm(&form)
	page.AddSelect("sections", hm.ToSelectOpt(hm.ToPtrSlice(sections)))
	page.AddSelect("users", hm.ToSelectOpt(hm.ToPtrSlice(users)))
	page.AddSelect("tags", hm.ToSelectOpt(hm.ToPtrSlice(tags)))
	page.AddSelect("kinds", kinds)

	if content.IsZero() {
		page.Name = "New Content"
		page.IsNew = true
		page.Form.SetAction(hm.CreatePath(&Content{}))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Content"
		page.IsNew = false
		page.Form.SetAction(hm.UpdatePath(&Content{}))
		page.Form.SetSubmitButtonText("Update")
	}

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&content, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-content")
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
