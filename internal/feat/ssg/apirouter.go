package ssg

import (
	"github.com/hermesgen/hm"
)

func NewAPIRouter(handler *APIHandler, mw []hm.Middleware, params hm.XParams) *hm.Router {
	core := hm.NewAPIRouter("api-router", params)
	core.SetMiddlewares(mw)

	// SSG API routes
	core.Post("/generate-markdown", handler.GenerateMarkdown)
	core.Post("/generate-html", handler.GenerateHTML)

	// Publish API routes
	core.Post("/publish", handler.Publish)

	// Layout API routes
	core.Get("/layouts", handler.GetAllLayouts)
	core.Get("/layouts/{id}", handler.GetLayout)
	core.Post("/layouts", handler.CreateLayout)
	core.Put("/layouts/{id}", handler.UpdateLayout)
	core.Delete("/layouts/{id}", handler.DeleteLayout)

	// Section API routes
	core.Get("/sections", handler.GetAllSections)
	core.Get("/sections/{id}", handler.GetSection)
	core.Post("/sections", handler.CreateSection)
	core.Put("/sections/{id}", handler.UpdateSection)
	core.Delete("/sections/{id}", handler.DeleteSection)

	// Content API routes
	core.Get("/contents", handler.GetAllContent)
	core.Get("/contents/{id}", handler.GetContent)
	core.Post("/contents", handler.CreateContent)
	core.Put("/contents/{id}", handler.UpdateContent)
	core.Delete("/contents/{id}", handler.DeleteContent)

	// Content-Tag API routes
	core.Post("/contents/{content_id}/tags", handler.AddTagToContent)
	core.Delete("/contents/{content_id}/tags/{tag_id}", handler.RemoveTagFromContent)

	// Content Image Upload API routes
	core.Post("/contents/{content_id}/images", handler.UploadContentImage)
	core.Get("/contents/{content_id}/images", handler.GetContentImages)
	core.Delete("/contents/{content_id}/images/delete", handler.DeleteContentImage)

	// Section Image Upload API routes
	core.Post("/sections/{section_id}/images", handler.UploadSectionImage)
	core.Delete("/sections/{section_id}/images/{image_type}", handler.DeleteSectionImage)

	// Tag API routes
	core.Get("/tags", handler.GetAllTags)
	core.Get("/tags/{id}", handler.GetTag)
	core.Get("/tags/name/{name}", handler.GetTagByName)
	core.Post("/tags", handler.CreateTag)
	core.Put("/tags/{id}", handler.UpdateTag)
	core.Delete("/tags/{id}", handler.DeleteTag)

	// Param API routes
	core.Get("/params", handler.ListParams)
	core.Get("/params/{id}", handler.GetParam)
	core.Get("/params/name/{name}", handler.GetParamByName)
	core.Get("/params/refkey/{ref_key}", handler.GetParamByRefKey)
	core.Post("/params", handler.CreateParam)
	core.Put("/params/{id}", handler.UpdateParam)
	core.Delete("/params/{id}", handler.DeleteParam)

	// Image API routes
	core.Get("/images", handler.ListImages)
	core.Get("/images/{id}", handler.GetImage)
	core.Get("/images/short/{short_id}", handler.GetImageByShortID)
	core.Post("/images", handler.CreateImage)
	core.Put("/images/{id}", handler.UpdateImage)
	core.Delete("/images/{id}", handler.DeleteImage)

	// Image Variant API routes
	core.Get("/images/{image_id}/variants", handler.ListImageVariantsByImageID)
	core.Get("/images/{image_id}/variants/{id}", handler.GetImageVariant)
	core.Post("/images/{image_id}/variants", handler.CreateImageVariant)
	core.Put("/images/{image_id}/variants/{id}", handler.UpdateImageVariant)
	core.Delete("/images/{image_id}/variants/{id}", handler.DeleteImageVariant)

	return core
}
