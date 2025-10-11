package ssg

import (
	"fmt"
	"html/template"
	"strings"
	
	"github.com/hermesgen/hm"
)

const (
	// WIP: This will be obtained from configuration.
	defaultAPIBaseURL = "http://localhost:8081/api/v1"
)

const (
	ssgFeat = "ssg"
	ssgPath = "/ssg"
)

type WebHandler struct {
	*hm.WebHandler
	apiClient *hm.APIClient
}

func NewWebHandler(tm *hm.TemplateManager, flash *hm.FlashManager, params hm.XParams) *WebHandler {
	// Register SSG-specific template functions
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
		WebHandler: handler,
		apiClient:  apiClient,
	}
}
