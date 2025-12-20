package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hermesgen/hm"

	"github.com/google/uuid"
)

func (h *APIHandler) GetAllContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllContentWithMeta", h.Name())

	var contents []Content
	var err error
	contents, err = h.svc.GetAllContentWithMeta(r.Context())
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetAllItems, hm.Cap(resContentName))
	h.OK(w, msg, contents)
}

func (h *APIHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var content Content
	content, err = h.svc.GetContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// The new GetContent service method should already include tags.
	// If not, this logic is a fallback.
	if len(content.Tags) == 0 {
		var tags []Tag
		tags, err = h.svc.GetTagsForContent(r.Context(), id)
		if err != nil {
			msg := fmt.Sprintf("Cannot get tags for content %s", id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
		content.Tags = tags
	}

	msg := fmt.Sprintf(hm.MsgGetItem, hm.Cap(resContentName))
	h.OK(w, msg, content)
}

func (h *APIHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("API: Handling CreateContent")

	var content Content
	var err error
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Log().Infof("API: Failed to decode content body: %v", err)
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	h.Log().Infof("API: Content decoded - Heading: %s, Body length: %d, SectionID: %s", content.Heading, len(content.Body), content.SectionID)

	content.GenCreateValues()
	h.Log().Infof("API: Generated create values - ID: %s", content.ID)

	err = h.svc.CreateContent(r.Context(), &content)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotCreateResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	for _, tag := range content.Tags {
		err = h.svc.AddTagToContent(r.Context(), content.ID, tag.Name)
		if err != nil {
			msg := fmt.Sprintf("Cannot add tag %s to content %s", tag.Name, content.ID)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
	}

	msg := fmt.Sprintf(hm.MsgCreateItem, hm.Cap(resContentName))
	h.Created(w, msg, content)
}

func (h *APIHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var content Content
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	content.SetID(id, true)
	content.GenUpdateValues()

	err = h.svc.UpdateContent(r.Context(), &content)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotUpdateResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	var existingTags []Tag
	existingTags, err = h.svc.GetTagsForContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf("Cannot get existing tags for content %s", id)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}
	for _, tag := range existingTags {
		err = h.svc.RemoveTagFromContent(r.Context(), id, tag.ID)
		if err != nil {
			msg := fmt.Sprintf("Cannot remove tag %s from content %s", tag.ID, id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
	}

	for _, tag := range content.Tags {
		err = h.svc.AddTagToContent(r.Context(), id, tag.Name)
		if err != nil {
			msg := fmt.Sprintf("Cannot add tag %s to content %s", tag.Name, id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
	}

	msg := fmt.Sprintf(hm.MsgUpdateItem, hm.Cap(resContentName))
	h.OK(w, msg, content)
}

func (h *APIHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotDeleteResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgDeleteItem, hm.Cap(resContentName))
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) AddTagToContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTagToContent", h.Name())

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var form AddTagToContentForm
	err = json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	err = h.svc.AddTagToContent(r.Context(), contentID, form.Name)
	if err != nil {
		msg := fmt.Sprintf("Cannot add tag %s to content %s", form.Name, contentID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Tag %s added to content %s", form.Name, contentID)
	h.Created(w, msg, nil)
}

func (h *APIHandler) RemoveTagFromContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling RemoveTagFromContent", h.Name())

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var tagIDStr string
	tagIDStr, err = h.Param(w, r, "tag_id")
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resTagName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var tagID uuid.UUID
	tagID, err = uuid.Parse(tagIDStr)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resTagName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.RemoveTagFromContent(r.Context(), contentID, tagID)
	if err != nil {
		msg := fmt.Sprintf("Cannot remove tag %s from content %s", tagID, contentID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Tag %s removed from content %s", tagID, contentID)
	h.OK(w, msg, json.RawMessage("null"))
}

// UploadContentImage handles image upload for content (header or content images)
func (h *APIHandler) UploadContentImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UploadContentImage", h.Name())

	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID format", err)
		return
	}

	imageTypeStr := r.FormValue("image_type")
	if imageTypeStr == "" {
		h.Err(w, http.StatusBadRequest, "Missing image_type parameter", nil)
		return
	}

	imageType := ImageType(imageTypeStr)
	if imageType != ImageTypeContent && imageType != ImageTypeHeader {
		h.Err(w, http.StatusBadRequest, "Invalid image_type for content", nil)
		return
	}

	altText := r.FormValue("alt_text")
	caption := r.FormValue("caption")

	file, header, err := r.FormFile("image")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Failed to parse uploaded file", err)
		return
	}
	defer file.Close()

	result, err := h.svc.UploadContentImage(r.Context(), contentID, file, header, imageType, altText, caption)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "Failed to upload image", err)
		return
	}

	msg := fmt.Sprintf("Image uploaded successfully: %s", result.Filename)
	h.OK(w, msg, map[string]interface{}{
		"filename":      result.Filename,
		"relative_path": result.RelativePath,
		"metadata":      result.Metadata,
	})
}

// GetContentImages returns all images for a specific content
func (h *APIHandler) GetContentImages(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContentImages", h.Name())

	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID format", err)
		return
	}

	images, err := h.svc.GetContentImages(r.Context(), contentID)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "Failed to get content images", err)
		return
	}

	msg := fmt.Sprintf("Retrieved %d images for content", len(images))
	h.OK(w, msg, map[string]interface{}{
		"images": images,
	})
}

// DeleteContentImageRequest represents the request body for deleting content images
type DeleteContentImageRequest struct {
	ImagePath string `json:"image_path"`
}

// DeleteContentImage handles deletion of content images by path
func (h *APIHandler) DeleteContentImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Infof("%s: Handling DeleteContentImage", h.Name())

	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		h.Log().Infof("Failed to parse content_id: %v", err)
		h.Err(w, http.StatusBadRequest, "Invalid content ID", err)
		return
	}
	h.Log().Infof("Content ID: %s", contentIDStr)

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		h.Log().Infof("Failed to parse UUID: %v", err)
		h.Err(w, http.StatusBadRequest, "Invalid content ID format", err)
		return
	}

	var req DeleteContentImageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Log().Infof("Failed to parse request body: %v", err)
		h.Err(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	h.Log().Infof("Image path: %s", req.ImagePath)

	err = h.svc.DeleteContentImage(r.Context(), contentID, req.ImagePath)
	if err != nil {
		h.Log().Infof("Service delete failed: %v", err)
		h.Err(w, http.StatusInternalServerError, "Failed to delete image", err)
		return
	}

	h.Log().Infof("Image deleted successfully: %s", req.ImagePath)
	msg := fmt.Sprintf("Content image deleted successfully")
	h.OK(w, msg, nil)
}
