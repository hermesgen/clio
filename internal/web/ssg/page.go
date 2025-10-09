package ssg

import (
	"net/http"

	"github.com/hermesgen/hm"
)

// ParamPage extends hm.Page to include a specific Param for templates.
type ParamPage struct {
	hm.Page
	Param Param
}

// NewParamPage creates a new ParamPage.
func NewParamPage(r *http.Request, param Param) *ParamPage {
	page := hm.NewPage(r, nil) // Pass nil for Data, as we'll use Param field
	return &ParamPage{
		Page:  *page,
		Param: param,
	}
}
