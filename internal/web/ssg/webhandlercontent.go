package ssg

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hermesgen/clio/internal/feat/auth"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
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
		h.Log().Infof("Error parsing form: %v", err)
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	h.Log().Infof("Form parsed - Heading: %s, Body: %s, SectionID: %s, UserID: %s", form.Heading, form.Body, form.SectionID, form.UserID)

	form.Validate()
	if form.HasErrors() {
		h.Log().Infof("Form validation failed - Body empty: %v", form.Body == "")
		content := ToFeatContent(form)
		webContent := ToWebContent(content)
		h.renderContentForm(w, r, form, webContent, "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToFeatContent(form)
	h.Log().Infof("Converted to content - Heading: %s, SectionID: %s, UserID: %s", content.Heading, content.SectionID, content.UserID)

	var response struct {
		Content feat.Content `json:"content"`
	}
	h.Log().Info("Calling API to create content...")
	err = h.apiClient.Post(h.addSiteSlugHeader(r), "/ssg/contents", content, &response)
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
	err := h.apiClient.Get(h.addSiteSlugHeader(r), path, &response)
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
	err = h.apiClient.Put(h.addSiteSlugHeader(r), path, content, nil)
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

	query := r.URL.Query()
	searchQuery := query.Get("search")
	pageStr := query.Get("page")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var response struct {
		Contents   []feat.Content `json:"contents"`
		Page       int            `json:"page"`
		TotalPages int            `json:"total_pages"`
		TotalCount int            `json:"total_count"`
		Search     string         `json:"search"`
	}

	url := fmt.Sprintf("/ssg/contents/search?page=%d", page)
	if searchQuery != "" {
		url += "&search=" + searchQuery
	}

	h.Log().Info("Calling apiClient.Get", "url", url)
	err := h.apiClient.Get(h.addSiteSlugHeader(r), url, &response)
	if err != nil {
		h.Err(w, err, "Cannot get contents from API", http.StatusInternalServerError)
		return
	}
	contents := response.Contents

	h.Log().Debugf("Contents received: %+v", contents)

	// Create page data with pagination info
	// Calculate pagination values
	prevPage := response.Page - 1
	nextPage := response.Page + 1
	showingFrom := (response.Page-1)*25 + 1
	showingTo := response.Page * 25
	if showingTo > response.TotalCount {
		showingTo = response.TotalCount
	}

	siteSlug, _ := feat.GetSiteSlugFromContext(r.Context())

	pageData := struct {
		hm.Page
		CurrentPage int    `json:"current_page"`
		TotalPages  int    `json:"total_pages"`
		TotalCount  int    `json:"total_count"`
		SearchQuery string `json:"search_query"`
		PageNumbers []int  `json:"page_numbers"`
		PrevPage    int    `json:"prev_page"`
		NextPage    int    `json:"next_page"`
		ShowingFrom int    `json:"showing_from"`
		ShowingTo   int    `json:"showing_to"`
		SiteSlug    string `json:"site_slug"`
	}{
		Page:        *hm.NewPage(r, contents),
		CurrentPage: response.Page,
		TotalPages:  response.TotalPages,
		TotalCount:  response.TotalCount,
		SearchQuery: searchQuery,
		PageNumbers: generatePageNumbers(response.Page, response.TotalPages),
		PrevPage:    prevPage,
		NextPage:    nextPage,
		ShowingFrom: showingFrom,
		ShowingTo:   showingTo,
		SiteSlug:    siteSlug,
	}

	pageData.Form.SetAction(ssgPath)
	pageData.SetFlash(h.GetFlash(r))

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-content")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pageData)
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
	err := h.apiClient.Get(h.addSiteSlugHeader(r), path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get content from API", http.StatusInternalServerError)
		return
	}
	content := response.Content

	page := hm.NewPage(r, content)
	page.Name = "Show Content"

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
	err := h.apiClient.Delete(h.addSiteSlugHeader(r), path)
	if err != nil {
		h.Err(w, err, "Failed to delete content via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Content deleted successfully")
	h.Redir(w, r, hm.ListPath(&Content{}), http.StatusSeeOther)
}

func (h *WebHandler) SearchContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Search content HTMX request")

	query := r.URL.Query()
	searchQuery := query.Get("search")
	pageStr := query.Get("page")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var response struct {
		Contents   []feat.Content `json:"contents"`
		Page       int            `json:"page"`
		TotalPages int            `json:"total_pages"`
		TotalCount int            `json:"total_count"`
		Search     string         `json:"search"`
	}

	url := fmt.Sprintf("/ssg/contents/search?page=%d", page)
	if searchQuery != "" {
		url += "&search=" + searchQuery
	}

	err := h.apiClient.Get(h.addSiteSlugHeader(r), url, &response)
	if err != nil {
		h.Err(w, err, "Cannot search contents from API", http.StatusInternalServerError)
		return
	}

	contents := response.Contents

	// Calculate pagination values
	prevPage := response.Page - 1
	nextPage := response.Page + 1
	showingFrom := (response.Page-1)*25 + 1
	showingTo := response.Page * 25
	if showingTo > response.TotalCount {
		showingTo = response.TotalCount
	}

	pageData := struct {
		Data        []feat.Content `json:"data"`
		CurrentPage int            `json:"current_page"`
		TotalPages  int            `json:"total_pages"`
		TotalCount  int            `json:"total_count"`
		SearchQuery string         `json:"search_query"`
		PageNumbers []int          `json:"page_numbers"`
		PrevPage    int            `json:"prev_page"`
		NextPage    int            `json:"next_page"`
		ShowingFrom int            `json:"showing_from"`
		ShowingTo   int            `json:"showing_to"`
	}{
		Data:        contents,
		CurrentPage: response.Page,
		TotalPages:  response.TotalPages,
		TotalCount:  response.TotalCount,
		SearchQuery: searchQuery,
		PageNumbers: generatePageNumbers(response.Page, response.TotalPages),
		PrevPage:    prevPage,
		NextPage:    nextPage,
		ShowingFrom: showingFrom,
		ShowingTo:   showingTo,
	}

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-content")
	if err != nil {
		h.Err(w, err, "Template not found", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "list-content-table", pageData)
	if err != nil {
		h.Err(w, err, "Cannot render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func generatePageNumbers(currentPage, totalPages int) []int {
	if totalPages <= 7 {
		pages := make([]int, totalPages)
		for i := 0; i < totalPages; i++ {
			pages[i] = i + 1
		}
		return pages
	}

	var pages []int
	if currentPage <= 4 {
		for i := 1; i <= 5; i++ {
			pages = append(pages, i)
		}
		pages = append(pages, -1)
		pages = append(pages, totalPages)
	} else if currentPage >= totalPages-3 {
		pages = append(pages, 1)
		pages = append(pages, -1)
		for i := totalPages - 4; i <= totalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		pages = append(pages, 1)
		pages = append(pages, -1)
		for i := currentPage - 1; i <= currentPage+1; i++ {
			pages = append(pages, i)
		}
		pages = append(pages, -1)
		pages = append(pages, totalPages)
	}
	return pages
}

func (h *WebHandler) renderContentForm(w http.ResponseWriter, r *http.Request, form ContentForm, content Content, errorMessage string, statusCode int) {
	h.Log().Debugf("h.apiClient: %+v", h.apiClient)

	// Get site mode to determine UI behavior
	siteMode := h.paramManager.GetSiteMode(r.Context())
	h.Log().Infof("Site mode: %s", siteMode)

	var sectionsResponse struct {
		Sections []Section `json:"sections"`
	}
	h.Log().Debug("Calling API to get sections")
	err := h.apiClient.Get(h.addSiteSlugHeader(r), "/ssg/sections", &sectionsResponse)
	if err != nil {
		h.Log().Errorf("Cannot get sections from API: %v", err)
		h.Err(w, err, "Cannot get sections from API", http.StatusInternalServerError)
		return
	}
	sections := sectionsResponse.Sections
	h.Log().Debugf("Sections received: %+v", sections)

	// In blog mode, only show root section
	if siteMode == "blog" {
		var rootSections []Section
		for _, s := range sections {
			if s.Name == "root" {
				rootSections = append(rootSections, s)
			}
		}
		sections = rootSections
	}

	var usersResponse struct {
		Users []auth.User `json:"users"`
	}
	h.Log().Debug("Calling API to get users")
	err = h.apiClient.Get(h.addSiteSlugHeader(r), "/auth/users", &usersResponse)
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
	err = h.apiClient.Get(h.addSiteSlugHeader(r), "/ssg/tags", &tagsResponse)
	if err != nil {
		h.Log().Errorf("Cannot get tags from API: %v", err)
		h.Err(w, err, "Cannot get tags from API", http.StatusInternalServerError)
		return
	}
	tags := tagsResponse.Tags
	h.Log().Debugf("Tags received: %+v", tags)

	// In blog mode, only "blog" content type is allowed
	var kinds []hm.SelectOpt
	if siteMode == "blog" {
		kinds = []hm.SelectOpt{
			{Value: "blog", Label: "Blog"},
		}
	} else {
		kinds = []hm.SelectOpt{
			{Value: "article", Label: "Article"},
			{Value: "page", Label: "Page"},
			{Value: "blog", Label: "Blog"},
			{Value: "series", Label: "Series"},
		}
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

func (h *WebHandler) GenerateHTML(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Generate HTML")

	siteSlug, ok := feat.GetSiteSlugFromContext(r.Context())
	if !ok {
		siteSlug = "structured"
	}

	err := h.apiClient.Post(h.addSiteSlugHeader(r), "/ssg/generate-html", nil, nil)
	if err != nil {
		h.FlashError(w, r, fmt.Sprintf("Failed to generate HTML: %v", err))
		h.Redir(w, r, "/ssg/list-content", http.StatusSeeOther)
		return
	}

	previewPort := h.Cfg().StrValOrDef("server.preview.port", "8082")
	previewURL := fmt.Sprintf("http://%s.localhost:%s/", siteSlug, previewPort)

	h.FlashSuccess(w, r, fmt.Sprintf("HTML generated successfully! Preview available at: %s", previewURL))
	h.Redir(w, r, "/ssg/list-content", http.StatusSeeOther)
}

