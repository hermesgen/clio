package ssg

import (
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
	handler := hm.NewWebHandler(tm, flash, params)
	apiClient := hm.NewAPIClient("web-api-client", func() string { return "" }, defaultAPIBaseURL, params)
	return &WebHandler{
		WebHandler: handler,
		apiClient:  apiClient,
	}
}
