package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hermesgen/hm"

	"github.com/google/uuid"
)

func (h *APIHandler) GetAllSections(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllSections", h.Name())

	var sections []Section
	var err error
	sections, err = h.svc.GetSections(r.Context())
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResources, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetAllItems, hm.Cap(resSectionName))
	h.OK(w, msg, sections)
}

func (h *APIHandler) GetSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetSection", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resSectionName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var section Section
	section, err = h.svc.GetSection(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotGetResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgGetItem, hm.Cap(resSectionName))
	h.OK(w, msg, section)
}

func (h *APIHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateSection", h.Name())

	var section Section
	var err error
	err = json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	newSection := NewSection(section.Name, section.Description, section.Path, section.LayoutID)
	newSection.GenCreateValues()

	err = h.svc.CreateSection(r.Context(), newSection)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotCreateResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgCreateItem, hm.Cap(resSectionName))
	h.Created(w, msg, newSection)
}

func (h *APIHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateSection", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resSectionName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var section Section
	err = json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		h.Err(w, http.StatusBadRequest, hm.ErrInvalidBody, err)
		return
	}

	updatedSection := NewSection(section.Name, section.Description, section.Path, section.LayoutID)
	updatedSection.SetID(id, true)
	updatedSection.GenUpdateValues()

	err = h.svc.UpdateSection(r.Context(), updatedSection)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotUpdateResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgUpdateItem, hm.Cap(resSectionName))
	h.OK(w, msg, updatedSection)
}

func (h *APIHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteSection", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrInvalidID, hm.Cap(resSectionName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteSection(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(hm.ErrCannotDeleteResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(hm.MsgDeleteItem, hm.Cap(resSectionName))
	h.OK(w, msg, json.RawMessage("null"))
}

// UploadSectionImage handles image upload for sections (section header or blog header)
func (h *APIHandler) UploadSectionImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UploadSectionImage", h.Name())

	sectionIDStr, err := h.Param(w, r, "section_id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid section ID", err)
		return
	}

	sectionID, err := uuid.Parse(sectionIDStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid section ID format", err)
		return
	}

	imageTypeStr := r.FormValue("image_type")
	if imageTypeStr == "" {
		h.Err(w, http.StatusBadRequest, "Missing image_type parameter", nil)
		return
	}

	imageType := ImageType(imageTypeStr)
	if imageType != ImageTypeSectionHeader && imageType != ImageTypeBlogHeader {
		h.Err(w, http.StatusBadRequest, "Invalid image_type for section", nil)
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

	result, err := h.svc.UploadSectionImage(r.Context(), sectionID, file, header, imageType, altText, caption)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "Failed to upload image", err)
		return
	}

	msg := fmt.Sprintf("Section image uploaded successfully: %s", result.Filename)
	h.OK(w, msg, map[string]interface{}{
		"filename":      result.Filename,
		"relative_path": result.RelativePath,
		"metadata":      result.Metadata,
	})
}

// DeleteSectionImage handles deletion of section images (section header or blog header)
func (h *APIHandler) DeleteSectionImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteSectionImage", h.Name())

	sectionIDStr, err := h.Param(w, r, "section_id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid section ID", err)
		return
	}

	sectionID, err := uuid.Parse(sectionIDStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid section ID format", err)
		return
	}

	imageTypeStr, err := h.Param(w, r, "image_type")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid image type", err)
		return
	}

	imageType := ImageType(imageTypeStr)
	if imageType != ImageTypeSectionHeader && imageType != ImageTypeBlogHeader {
		h.Err(w, http.StatusBadRequest, "Invalid image_type for section", nil)
		return
	}

	err = h.svc.DeleteSectionImage(r.Context(), sectionID, imageType)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "Failed to delete image", err)
		return
	}

	msg := fmt.Sprintf("Section image deleted successfully")
	h.OK(w, msg, nil)
}
