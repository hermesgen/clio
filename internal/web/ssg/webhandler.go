package ssg

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

const (
	defaultAPIBaseURL = "http://localhost:8081/api/v1"
)

const (
	ssgFeat = "ssg"
	ssgPath = "/ssg"
)

type WebHandler struct {
	*hm.WebHandler
	apiClient      *hm.APIClient
	paramManager   *feat.ParamManager
	siteManager    *feat.SiteManager
	sessionManager interface {
		SetUserSession(w http.ResponseWriter, userID uuid.UUID, siteSlug string) error
		GetUserSession(r *http.Request) (userID uuid.UUID, siteSlug string, err error)
	}
}

func (wh *WebHandler) addSiteSlugHeader(r *http.Request) *http.Request {
	ctx := r.Context()
	if siteSlug, ok := feat.GetSiteSlugFromContext(ctx); ok && siteSlug != "" {
		r.Header.Set("X-Site-Slug", siteSlug)
	}
	return r
}

func NewWebHandler(tm *hm.TemplateManager, flash *hm.FlashManager, paramManager *feat.ParamManager, siteManager *feat.SiteManager, sessionManager interface {
	SetUserSession(w http.ResponseWriter, userID uuid.UUID, siteSlug string) error
	GetUserSession(r *http.Request) (userID uuid.UUID, siteSlug string, err error)
}, params hm.XParams) *WebHandler {
	ssgFunctions := template.FuncMap{
		"newPath": func(entityType string) string {
			return fmt.Sprintf("/ssg/new-%s", strings.ToLower(entityType))
		},
		"listPath": func(entityType string) string {
			return fmt.Sprintf("/ssg/list-%s", strings.ToLower(hm.Plural(entityType)))
		},
		"editPath": func(entityType, id string) string {
			return fmt.Sprintf("/ssg/edit-%s?id=%s", strings.ToLower(entityType), id)
		},
	}
	
	tm.RegisterFunctions(ssgFunctions)
	
	handler := hm.NewWebHandler(tm, flash, params)
	apiClient := hm.NewAPIClient("web-api-client", func() string { return "" }, defaultAPIBaseURL, params)
	return &WebHandler{
		WebHandler:     handler,
		apiClient:      apiClient,
		paramManager:   paramManager,
		siteManager:    siteManager,
		sessionManager: sessionManager,
	}
}
