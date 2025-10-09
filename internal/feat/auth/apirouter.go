package auth

import (
	"github.com/hermesgen/hm"
)

func NewAPIRouter(handler *APIHandler, mw []hm.Middleware, params hm.XParams) *hm.Router {
	core := hm.NewAPIRouter("api-router", params)
	core.SetMiddlewares(mw)

	// User API routes
	core.Get("/users", handler.GetAllUsers)
	core.Get("/users/{id}", handler.GetUser)
	core.Post("/users", handler.CreateUser)
	core.Put("/users/{id}", handler.UpdateUser)
	core.Delete("/users/{id}", handler.DeleteUser)

	return core
}
