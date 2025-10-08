package ssg

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	hm "github.com/hermesgen/hm"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func (h *WebHandler) ListParams(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List params")

	var response struct {
		Params []feat.Param `json:"params"`
	}
	err := h.apiClient.Get(r, "/ssg/params", &response)
	if err != nil {
		h.Err(w, err, "Cannot get params from API", http.StatusInternalServerError)
		return
	}
	params := response.Params

	page := hm.NewPage(r, ToWebParams(params))
	page.Form.SetAction(ssgPath)

	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(&Param{})

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-params")
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

func (h *WebHandler) NewParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New param form")
	form := NewParamForm(r)
	h.renderParamForm(w, r, form, NewParam("", ""), "", http.StatusOK)
}

func (h *WebHandler) CreateParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create param")

	form, err := ParamFormFromRequest(r)
	if err != nil {
		h.renderParamForm(w, r, form, NewParam("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		param := ToFeatParam(form)
		webParam := ToWebParam(param)
		h.renderParamForm(w, r, form, webParam, "Validation failed", http.StatusBadRequest)
		return
	}

	param := ToFeatParam(form)

	var response struct {
		Param feat.Param `json:"param"`
	}
	err = h.apiClient.Post(r, "/ssg/params", param, &response)
	if err != nil {
		h.Err(w, err, "Failed to create param via API", http.StatusInternalServerError)
		return
	}
	createdParam := ToWebParam(response.Param)

	if hm.IsHTMXRequest(r) {
		redirectURL := hm.EditPath(&createdParam, createdParam.GetID())
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.FlashInfo(w, r, "Param created")
	h.Redir(w, r, hm.EditPath(&createdParam, createdParam.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) EditParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit param")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing param ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Param feat.Param `json:"param"`
	}
	path := fmt.Sprintf("/ssg/params/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get param from API", http.StatusInternalServerError)
		return
	}
	param := response.Param
	webParam := ToWebParam(param)

	form := ToParamForm(r, param)
	h.renderParamForm(w, r, form, webParam, "", http.StatusOK)
}

func (h *WebHandler) UpdateParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update param")

	form, err := ParamFormFromRequest(r)
	if err != nil {
		h.renderParamForm(w, r, form, NewParam("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	form.Validate()
	if form.HasErrors() {
		param := ToFeatParam(form)
		webParam := ToWebParam(param)
		h.renderParamForm(w, r, form, webParam, "Validation failed", http.StatusBadRequest)
		return
	}

	param := ToFeatParam(form)

	path := fmt.Sprintf("/ssg/params/%s", param.GetID())
	err = h.apiClient.Put(r, path, param, nil)
	if err != nil {
		h.Err(w, err, "Failed to update param via API", http.StatusInternalServerError)
		return
	}

	if hm.IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<div id=\"save-status\" data-timestamp=\"" + time.Now().Format(hm.TimeFormat) + "\"></div>"))
		return
	}

	h.FlashInfo(w, r, "Param updated successfully")
	h.Redir(w, r, hm.EditPath(&Param{}, param.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) DeleteParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete param")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing param ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/params/%s", idStr)
	err := h.apiClient.Delete(r, path)
	if err != nil {
		h.Err(w, err, "Failed to delete param via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Param deleted successfully")
	h.Redir(w, r, hm.ListPath(&Param{}), http.StatusSeeOther)
}

func (h *WebHandler) ShowParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show param")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing param ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Param feat.Param `json:"param"`
	}
	path := fmt.Sprintf("/ssg/params/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get param from API", http.StatusInternalServerError)
		return
	}

	param := ToWebParam(response.Param)

	page := hm.NewPage(r, param)
	page.Name = "Show Param"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&param, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-param")
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

func (h *WebHandler) renderParamForm(w http.ResponseWriter, r *http.Request, form ParamForm, param Param, errorMessage string, statusCode int) {
	paramPage := NewParamPage(r, param)
	paramPage.SetForm(&form)

	if param.IsZero() {
		paramPage.Name = "New Param"
		paramPage.IsNew = true
		paramPage.Form.SetAction(hm.CreatePath(&Param{}))
		paramPage.Form.SetSubmitButtonText("Create")
	} else {
		paramPage.Name = "Edit Param"
		paramPage.IsNew = false
		paramPage.Form.SetAction(hm.UpdatePath(&Param{}))
		paramPage.Form.SetSubmitButtonText("Update")
	}

	menu := paramPage.NewMenu(ssgPath)
	menu.AddListItem(&param, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-param")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	paramPage.SetFlash(h.GetFlash(r))

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, paramPage)
	if err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, statusCode)
}
