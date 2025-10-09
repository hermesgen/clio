package auth

import (
	"net/http"

	"github.com/hermesgen/hm"
)

const (
	resUserName    = "user"
	resUserNameCap = "User"
)

type APIHandler struct {
	*hm.APIHandler
	svc Service
}

func NewAPIHandler(name string, service Service, params hm.XParams) *APIHandler {
	h := hm.NewAPIHandler(name, params)
	return &APIHandler{
		APIHandler: h,
		svc:        service,
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
	case User:
		return map[string]interface{}{"user": v}

	// Slices of entities
	case []User:
		return map[string]interface{}{"users": v}

	// Default case for nil, maps, or other types
	default:
		return data
	}
}
