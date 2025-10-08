package ssg

import (
	hm "github.com/hermesgen/hm"
)

func NewWebRouter(handler *WebHandler, mw []hm.Middleware, params hm.XParams) *hm.Router {
	core := hm.NewWebRouter("app-ssg-web-router", params)
	core.SetMiddlewares(mw)

	// Content routes
	core.Get("/new-content", handler.NewContent)
	core.Post("/create-content", handler.CreateContent)
	core.Get("/edit-content", handler.EditContent)
	core.Post("/update-content", handler.UpdateContent)
	core.Get("/list-content", handler.ListContent)
	core.Get("/show-content", handler.ShowContent)
	core.Post("/delete-content", handler.DeleteContent)
	// Section routes
	core.Get("/new-section", handler.NewSection)
	core.Post("/create-section", handler.CreateSection)
	core.Get("/edit-section", handler.EditSection)
	core.Post("/update-section", handler.UpdateSection)
	core.Get("/list-sections", handler.ListSections)
	core.Get("/show-section", handler.ShowSection)
	core.Post("/delete-section", handler.DeleteSection)

	// Tag routes
	core.Get("/new-tag", handler.NewTag)
	core.Post("/create-tag", handler.CreateTag)
	core.Get("/edit-tag", handler.EditTag)
	core.Post("/update-tag", handler.UpdateTag)
	core.Get("/list-tags", handler.ListTags)
	core.Get("/show-tag", handler.ShowTag)
	core.Post("/delete-tag", handler.DeleteTag)

	// Layout routes
	core.Get("/new-layout", handler.NewLayout)
	core.Post("/create-layout", handler.CreateLayout)
	core.Get("/edit-layout", handler.EditLayout)
	core.Post("/update-layout", handler.UpdateLayout)
	core.Get("/list-layouts", handler.ListLayouts)
	core.Get("/show-layout", handler.ShowLayout)
	core.Post("/delete-layout", handler.DeleteLayout)

	// Param routes
	core.Get("/new-param", handler.NewParam)
	core.Post("/create-param", handler.CreateParam)
	core.Get("/edit-param", handler.EditParam)
	core.Post("/update-param", handler.UpdateParam)
	core.Get("/list-params", handler.ListParams)
	core.Get("/show-param", handler.ShowParam)
	core.Post("/delete-param", handler.DeleteParam)

	// Image routes
	core.Get("/new-image", handler.NewImage)
	core.Post("/create-image", handler.CreateImage)
	core.Get("/edit-image", handler.EditImage)
	core.Post("/update-image", handler.UpdateImage)
	core.Get("/list-images", handler.ListImages)
	core.Get("/show-image", handler.ShowImage)
	core.Post("/delete-image", handler.DeleteImage)

	// Image Variant routes
	core.Get("/images/:imageID/variants/new", handler.NewImageVariant)
	core.Post("/images/:imageID/variants", handler.CreateImageVariant)
	core.Get("/images/:imageID/variants/:id/edit", handler.EditImageVariant)
	core.Post("/images/:imageID/variants/:id", handler.UpdateImageVariant)
	core.Get("/images/:imageID/variants", handler.ListImageVariants)
	core.Get("/images/:imageID/variants/:id", handler.ShowImageVariant)
	core.Post("/images/:imageID/variants/:id/delete", handler.DeleteImageVariant)

	return core
}
