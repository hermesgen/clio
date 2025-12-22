package ssg

import (
	"net/http"

	"github.com/hermesgen/hm"
)

type RootRouter struct {
	*hm.WebHandler
	ssgWebHandler *WebHandler
}

func NewRootRouter(ssgWebHandler *WebHandler, params hm.XParams) *RootRouter {
	return &RootRouter{
		WebHandler:    hm.NewWebHandler(nil, nil, params),
		ssgWebHandler: ssgWebHandler,
	}
}

func (rr *RootRouter) SetupRoutes(router *hm.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			rr.ssgWebHandler.RootRedirect(w, r)
		}
	})
}
