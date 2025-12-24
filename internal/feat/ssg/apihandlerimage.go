package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hermesgen/hm"

	"github.com/google/uuid"
)

func (h *APIHandler) CreateImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateImage", h.Name())

	var image Image
	var err error
	err = json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	newImage := NewImage() // Call constructor without arguments
	// Assign fields from the decoded JSON to the newImage instance
	newImage.SiteID = image.SiteID
	newImage.FileName = image.FileName
	newImage.FilePath = image.FilePath
	newImage.Width = image.Width
	newImage.Height = image.Height
	newImage.Title = image.Title
	newImage.AltText = image.AltText

	newImage.GenCreateValues()

	err = h.svc.CreateImage(r.Context(), &newImage)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotCreateResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgCreateItem, hm.Cap(resImageName))
	h.Created(w, msg, newImage)
}

func (h *APIHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetImage", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var image Image
	image, err = h.svc.GetImage(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetItem, hm.Cap(resImageName))
	h.OK(w, msg, image)
}

func (h *APIHandler) GetImageByShortID(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetImageByShortID", h.Name())

	var err error
	var shortID string
	shortID, err = h.Param(w, r, "short_id")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", hm.ErrInvalidParam, "short_id")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var image Image
	image, err = h.svc.GetImageByShortID(r.Context(), shortID)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetItem, hm.Cap(resImageName))
	h.OK(w, msg, image)
}

func (h *APIHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListImages", h.Name())

	var images []Image
	var err error
	images, err = h.svc.ListImages(r.Context())
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetAllItems, hm.Cap(resImageName))
	h.OK(w, msg, images)
}

func (h *APIHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateImage", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var image Image
	err = json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	updatedImage := NewImage()   // Call constructor without arguments
	updatedImage.SetID(id, true) // Set the ID from the URL on the decoded content

	// Assign fields from the decoded JSON to the updatedImage instance
	updatedImage.SiteID = image.SiteID
	updatedImage.FileName = image.FileName
	updatedImage.FilePath = image.FilePath
	updatedImage.Width = image.Width
	updatedImage.Height = image.Height
	updatedImage.Title = image.Title
	updatedImage.AltText = image.AltText

	updatedImage.GenUpdateValues()

	err = h.svc.UpdateImage(r.Context(), &updatedImage)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotUpdateResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgUpdateItem, hm.Cap(resImageName))
	h.OK(w, msg, updatedImage)
}

func (h *APIHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteImage", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteImage(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotDeleteResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgDeleteItem, hm.Cap(resImageName))
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) CreateImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateImageVariant", h.Name())

	var variant ImageVariant
	var err error
	err = json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	newVariant := NewImageVariant() // Call constructor without arguments
	// Assign fields from the decoded JSON to the newVariant instance
	newVariant.ImageID = variant.ImageID
	newVariant.Kind = variant.Kind
	newVariant.Width = variant.Width
	newVariant.Height = variant.Height
	newVariant.FilesizeByte = variant.FilesizeByte
	newVariant.Mime = variant.Mime
	newVariant.BlobRef = variant.BlobRef

	newVariant.GenCreateValues()

	err = h.svc.CreateImageVariant(r.Context(), &newVariant)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotCreateResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgCreateItem, hm.Cap(resImageVariantName))
	h.Created(w, msg, newVariant)
}

func (h *APIHandler) GetImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetImageVariant", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageVariantName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var variant ImageVariant
	variant, err = h.svc.GetImageVariant(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetItem, hm.Cap(resImageVariantName))
	h.OK(w, msg, variant)
}

func (h *APIHandler) ListImageVariantsByImageID(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListImageVariantsByImageID", h.Name())

	var err error
	var imageIDStr string
	imageIDStr, err = h.Param(w, r, "image_id")
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var imageID uuid.UUID
	imageID, err = uuid.Parse(imageIDStr)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var variants []ImageVariant
	variants, err = h.svc.ListImageVariantsByImageID(r.Context(), imageID)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetAllItems, hm.Cap(resImageVariantName))
	h.OK(w, msg, variants)
}

func (h *APIHandler) UpdateImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateImageVariant", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageVariantName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var variant ImageVariant
	err = json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	updatedVariant := NewImageVariant() // Call constructor without arguments
	updatedVariant.SetID(id, true)      // Set the ID from the URL on the decoded content

	// Assign fields from the decoded JSON to the updatedVariant instance
	updatedVariant.ImageID = variant.ImageID
	updatedVariant.Kind = variant.Kind
	updatedVariant.Width = variant.Width
	updatedVariant.Height = variant.Height
	updatedVariant.FilesizeByte = variant.FilesizeByte
	updatedVariant.Mime = variant.Mime
	updatedVariant.BlobRef = variant.BlobRef

	updatedVariant.GenUpdateValues()

	err = h.svc.UpdateImageVariant(r.Context(), &updatedVariant)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotUpdateResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgUpdateItem, hm.Cap(resImageVariantName))
	h.OK(w, msg, updatedVariant)
}

func (h *APIHandler) DeleteImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteImageVariant", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resImageVariantName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteImageVariant(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotDeleteResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgDeleteItem, hm.Cap(resImageVariantName))
	h.OK(w, msg, json.RawMessage("null"))
}
