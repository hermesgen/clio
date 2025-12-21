package ssg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hermesgen/hm"
)

const (
	resContentName      = "content"
	resSectionName      = "section"
	resLayoutName       = "layout"
	resTagName          = "tag"
	resParamName        = "param"
	resImageName        = "image"
	resImageVariantName = "image variant"
)

type APIHandler struct {
	*hm.APIHandler
	svc         Service
	siteManager *SiteManager
}

func NewAPIHandler(name string, service Service, siteManager *SiteManager, params hm.XParams) *APIHandler {
	return &APIHandler{
		APIHandler:  hm.NewAPIHandler(name, params),
		svc:         service,
		siteManager: siteManager,
	}
}

func (h *APIHandler) OK(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.OK(w, message, wrappedData)
}

func (h *APIHandler) Created(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.Created(w, message, wrappedData)
}

func (h *APIHandler) wrapData(data interface{}) interface{} {
	switch v := data.(type) {
	// Single entities
	case Site:
		return map[string]interface{}{"site": v}
	case Layout:
		return map[string]interface{}{"layout": v}
	case Section:
		return map[string]interface{}{"section": v}
	case Content:
		return map[string]interface{}{"content": v}
	case Tag:
		return map[string]interface{}{"tag": v}
	case Param:
		return map[string]interface{}{"param": v}
	case Image:
		return map[string]interface{}{"image": v}
	case ImageVariant:
		return map[string]interface{}{"image_variant": v}

	// Slices of entities
	case []Site:
		return map[string]interface{}{"sites": v}
	case []Layout:
		return map[string]interface{}{"layouts": v}
	case []Section:
		return map[string]interface{}{"sections": v}
	case []Content:
		return map[string]interface{}{"contents": v}
	case []Tag:
		return map[string]interface{}{"tags": v}
	case []Param:
		return map[string]interface{}{"params": v}
	case []Image:
		return map[string]interface{}{"images": v}
	case []ImageVariant:
		return map[string]interface{}{"image_variants": v}

	// Default case for nil, maps, or other types
	default:
		return data
	}
}

func (h *APIHandler) Publish(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling Publish", h.Name())

	var err error

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	var data PublishRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	// Run the publish process
	commitURL, err := h.svc.Publish(r.Context(), data.Message)
	if err != nil {
		msg := fmt.Sprintf("Cannot publish: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "Publish process started successfully"
	result := map[string]string{"commitURL": commitURL}
	h.OK(w, msg, result)
}

func (h *APIHandler) GenerateMarkdown(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GenerateMarkdown", h.Name())

	var err error
	err = h.svc.GenerateMarkdown(r.Context())
	if err != nil {
		msg := fmt.Sprintf("Cannot generate markdown: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "Markdown generation process started successfully"
	h.OK(w, msg, nil)
}

func (h *APIHandler) GenerateHTML(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GenerateHTML", h.Name())

	var err error
	err = h.svc.GenerateHTMLFromContent(r.Context())
	if err != nil {
		msg := fmt.Sprintf("Cannot generate HTML: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "HTML generation process started successfully"
	h.OK(w, msg, nil)
}

// PublishRequest represents the data for a publish request.
type PublishRequest struct {
	Message string `json:"message"`
}

// AddTagToContentForm represents the data for adding a tag to content.
type AddTagToContentForm struct {
	Name string `json:"name"`
}
