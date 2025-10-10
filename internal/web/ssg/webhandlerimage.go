package ssg

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func (h *WebHandler) NewImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New image form")
	form := NewImageForm(r) // Pass r
	h.renderImageForm(w, r, form, NewImage("", "", "", "", "", "", 0, 0, 0), "", http.StatusOK)
}

func (h *WebHandler) CreateImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create image")

	err := r.ParseMultipartForm(32 << 20) // 32MB max upload size
	if err != nil {
		h.Err(w, err, "Cannot parse multipart form", http.StatusBadRequest)
		return
	}

	form := NewImageForm(r) // Pass r
	form.Name = r.FormValue("name")
	form.Description = r.FormValue("description")
	form.AltText = r.FormValue("altText")

	file, header, err := r.FormFile("file")
	if err != nil {
		form.BaseForm.Validation().AddFieldError("file", "", "Image file is required") // Use BaseForm.Validation()
		h.renderImageForm(w, r, form, NewImage("", "", "", "", "", "", 0, 0, 0), "Validation failed", http.StatusBadRequest)
		return
	}
	defer file.Close()

	form.File = header
	form.Validate()

	if form.HasErrors() {
		h.renderImageForm(w, r, form, NewImage("", "", "", "", "", "", 0, 0, 0), "Validation failed", http.StatusBadRequest)
		return
	}

	// Determine upload directory
	uploadDir, _ := h.Cfg().StrVal(feat.SSGKey.ImagesPath)
	if err = os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		h.Err(w, err, "Cannot create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(header.Filename))
	destPath := filepath.Join(uploadDir, filename)

	// Save the file
	dst, err := os.Create(destPath)
	if err != nil {
		h.Err(w, err, "Cannot create file on disk", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		h.Err(w, err, "Cannot save file to disk", http.StatusInternalServerError)
		return
	}

	// Get image dimensions and mime type
	imgFile, err := os.Open(destPath)
	if err != nil {
		h.Err(w, err, "Cannot open saved image file", http.StatusInternalServerError)
		return
	}
	defer imgFile.Close()

	imgConfig, format, err := image.DecodeConfig(imgFile)
	if err != nil {
		h.Err(w, err, "Cannot decode image config", http.StatusInternalServerError)
		return
	}

	// Reset file pointer for subsequent reads if needed, though not strictly necessary here
	imgFile.Seek(0, io.SeekStart)

	// Construct feat.Image
	featImage := ToFeatImage(form)
	featImage.Mime = "image/" + format
	featImage.FilesizeByte = header.Size
	featImage.Width = imgConfig.Width
	featImage.Height = imgConfig.Height

	var response struct {
		Image feat.Image `json:"image"`
	}
	err = h.apiClient.Post(r, "/ssg/images", featImage, &response)
	if err != nil {
		h.Err(w, err, "Failed to create image via API", http.StatusInternalServerError)
		return
	}
	createdImage := ToWebImage(response.Image)

	// Create a default variant for the image
	featImageVariant := feat.NewImageVariant()
	featImageVariant.ImageID = createdImage.ID
	featImageVariant.Kind = "original"
	featImageVariant.BlobRef = "/static/images/" + filename // Use BlobRef
	featImageVariant.Mime = "image/" + format
	featImageVariant.FilesizeByte = header.Size
	featImageVariant.Width = imgConfig.Width
	featImageVariant.Height = imgConfig.Height

	var variantResponse struct {
		ImageVariant feat.ImageVariant `json:"imageVariant"`
	}
	err = h.apiClient.Post(r, fmt.Sprintf("/ssg/images/%s/variants", createdImage.ID), featImageVariant, &variantResponse)
	if err != nil {
		h.Err(w, err, "Failed to create original image variant via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Image created successfully")
	h.Redir(w, r, hm.EditPath(&Image{}, createdImage.GetID()), http.StatusSeeOther)
}

func (h *WebHandler) EditImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit image")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing image ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Image feat.Image `json:"image"`
	}
	path := fmt.Sprintf("/ssg/images/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get image from API", http.StatusInternalServerError)
		return
	}
	image := response.Image

	form := ToImageForm(r, ToWebImage(image)) // Pass r
	h.renderImageForm(w, r, form, ToWebImage(image), "", http.StatusOK)
}

func (h *WebHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update image")

	err := r.ParseForm()
	if err != nil {
		h.Err(w, err, "Cannot parse form", http.StatusBadRequest)
		return
	}

	form := NewImageForm(r) // Pass r
	form.ID = r.FormValue("id")
	form.Name = r.FormValue("name")
	form.Description = r.FormValue("description")
	form.AltText = r.FormValue("altText")
	// Note: File upload is not handled in update for simplicity, assuming image content is immutable after creation

	form.Validate()
	if form.HasErrors() {
		// Re-fetch image to pass to render function, as form doesn't contain full image data
		id, _ := uuid.Parse(form.ID)
		var currentImageResponse struct {
			Image feat.Image `json:"image"`
		}
		err = h.apiClient.Get(r, fmt.Sprintf("/ssg/images/%s", id), &currentImageResponse) // Declare err
		if err != nil {
			h.Err(w, err, "Cannot get image from API for validation", http.StatusInternalServerError)
			return
		}
		h.renderImageForm(w, r, form, ToWebImage(currentImageResponse.Image), "Validation failed", http.StatusBadRequest)
		return
	}

	featImage := ToFeatImage(form)
	// Mime, FilesizeByte, Width, Height are not updated here, as file upload is not handled

	path := fmt.Sprintf("/ssg/images/%s", featImage.GetID())
	err = h.apiClient.Put(r, path, featImage, nil)
	if err != nil {
		h.Err(w, err, "Failed to update image via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Image updated successfully")
	h.Redir(w, r, hm.ListPath(&Image{}), http.StatusSeeOther)
}

func (h *WebHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List images")

	var response struct {
		Images []feat.Image `json:"images"`
	}
	err := h.apiClient.Get(r, "/ssg/images", &response)
	if err != nil {
		h.Err(w, err, "Cannot get images from API", http.StatusInternalServerError)
		return
	}
	images := ToWebImages(response.Images)

	page := hm.NewPage(r, images)
	page.Form.SetAction(ssgPath)
	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(&Image{})

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-images")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *WebHandler) ShowImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show image")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing image ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Image feat.Image `json:"image"`
	}
	path := fmt.Sprintf("/ssg/images/%s", idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get image from API", http.StatusInternalServerError)
		return
	}

	image := ToWebImage(response.Image)

	page := hm.NewPage(r, image)
	page.Name = "Show Image"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&image, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-image")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, http.StatusOK)
}

func (h *WebHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete image")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing image ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/images/%s", idStr)
	err := h.apiClient.Delete(r, path)
	if err != nil {
		h.Err(w, err, "Failed to delete image via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Image deleted successfully")
	h.Redir(w, r, hm.ListPath(&Image{}), http.StatusSeeOther)
}

func (h *WebHandler) renderImageForm(w http.ResponseWriter, r *http.Request, form ImageForm, image Image, errorMessage string, statusCode int) {
	page := hm.NewPage(r, image)
	page.SetForm(&form)

	if image.IsZero() {
		page.Name = "New Image"
		page.IsNew = true
		page.Form.SetAction(hm.CreatePath(&Image{}))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Image"
		page.IsNew = false
		page.Form.SetAction(hm.UpdatePath(&Image{}))
		page.Form.SetSubmitButtonText("Update")
	}

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&image)

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-image")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	page.SetFlash(h.GetFlash(r))

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, statusCode)
}

// NewImageVariant displays the form for creating a new image variant.
func (h *WebHandler) NewImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New image variant form")

	imageIDStr := r.URL.Query().Get("imageID")
	if imageIDStr == "" {
		h.Err(w, nil, "Missing image ID for variant", http.StatusBadRequest)
		return
	}

	form := NewImageVariantForm(r) // Pass r
	form.ImageID = imageIDStr
	h.renderImageVariantForm(w, r, form, NewImageVariant(uuid.Nil, "", "", "", "", 0, 0, 0), "", http.StatusOK)
}

// CreateImageVariant handles the creation of a new image variant.
func (h *WebHandler) CreateImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create image variant")

	parseFormErr := r.ParseForm()
	if parseFormErr != nil {
		h.Err(w, parseFormErr, "Cannot parse form", http.StatusBadRequest)
		return
	}

	form := NewImageVariantForm(r) // Pass r
	form.ID = r.FormValue("id")    // Add this line to get ID from form
	form.ImageID = r.FormValue("imageID")
	form.Name = r.FormValue("name")

	form.Validate()
	if form.HasErrors() {
		// Re-fetch variant to pass to render function
		imageID, _ := uuid.Parse(form.ImageID)
		id, _ := uuid.Parse(form.ID)
		var currentVariantResponse struct {
			ImageVariant feat.ImageVariant `json:"imageVariant"`
		}
		apiGetErr := h.apiClient.Get(r, fmt.Sprintf("/ssg/images/%s/variants/%s", imageID, id), &currentVariantResponse)
		if apiGetErr != nil {
			h.Err(w, apiGetErr, "Cannot get image variant from API for validation", http.StatusInternalServerError)
			return
		}
		h.renderImageVariantForm(w, r, form, ToWebImageVariant(currentVariantResponse.ImageVariant), "Validation failed", http.StatusBadRequest)
		return
	}

	imageID, parseUUIDErr := uuid.Parse(form.ImageID)
	if parseUUIDErr != nil {
		h.Err(w, parseUUIDErr, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Fetch parent image to get path, url, mimetype, size, width, height
	var imageResponse struct {
		Image feat.Image `json:"image"`
	}
	imagePath := fmt.Sprintf("/ssg/images/%s", imageID)
	apiGetErr := h.apiClient.Get(r, imagePath, &imageResponse)
	if apiGetErr != nil {
		h.Err(w, apiGetErr, "Cannot get parent image from API", http.StatusInternalServerError)
		return
	}
	parentImage := imageResponse.Image

	featImageVariant := feat.NewImageVariant()
	featImageVariant.ImageID = imageID
	featImageVariant.Kind = form.Name
	featImageVariant.BlobRef = "/static/images/" + uuid.New().String() + "." + parentImage.Mime // Use BlobRef and generate a new unique name
	featImageVariant.Mime = parentImage.Mime
	featImageVariant.FilesizeByte = parentImage.FilesizeByte
	featImageVariant.Width = parentImage.Width
	featImageVariant.Height = parentImage.Height

	var response struct {
		ImageVariant feat.ImageVariant `json:"imageVariant"`
	}
	apiPath := fmt.Sprintf("/ssg/images/%s/variants", imageID)
	apiPostErr := h.apiClient.Post(r, apiPath, featImageVariant, &response)
	if apiPostErr != nil {
		h.Err(w, apiPostErr, "Failed to create image variant via API", http.StatusInternalServerError)
		return
	}
	createdVariant := ToWebImageVariant(response.ImageVariant)

	h.FlashInfo(w, r, "Image variant created successfully")
	h.Redir(w, r, hm.EditPath(&ImageVariant{}, createdVariant.GetID()), http.StatusSeeOther)
}

// EditImageVariant displays the form for editing an existing image variant.
func (h *WebHandler) EditImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit image variant")

	imageIDStr := r.URL.Query().Get("imageID")
	idStr := r.URL.Query().Get("id")
	if imageIDStr == "" || idStr == "" {
		h.Err(w, nil, "Missing image ID or variant ID", http.StatusBadRequest)
		return
	}

	var response struct {
		ImageVariant feat.ImageVariant `json:"imageVariant"`
	}
	path := fmt.Sprintf("/ssg/images/%s/variants/%s", imageIDStr, idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get image variant from API", http.StatusInternalServerError)
		return
	}
	variant := response.ImageVariant

	form := ToImageVariantForm(r, ToWebImageVariant(variant)) // Pass r
	h.renderImageVariantForm(w, r, form, ToWebImageVariant(variant), "", http.StatusOK)
}

// UpdateImageVariant handles the update of an existing image variant.
func (h *WebHandler) UpdateImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update image variant")

	parseFormErr := r.ParseForm()
	if parseFormErr != nil {
		h.Err(w, parseFormErr, "Cannot parse form", http.StatusBadRequest)
		return
	}

	form := NewImageVariantForm(r) // Pass r
	form.ID = r.FormValue("id")
	form.ImageID = r.FormValue("imageID")
	form.Name = r.FormValue("name")

	form.Validate()
	if form.HasErrors() {
		// Re-fetch variant to pass to render function
		imageID, _ := uuid.Parse(form.ImageID)
		id, _ := uuid.Parse(form.ID)
		var currentVariantResponse struct {
			ImageVariant feat.ImageVariant `json:"imageVariant"`
		}
		apiGetErr := h.apiClient.Get(r, fmt.Sprintf("/ssg/images/%s/variants/%s", imageID, id), &currentVariantResponse)
		if apiGetErr != nil {
			h.Err(w, apiGetErr, "Cannot get image variant from API for validation", http.StatusInternalServerError)
			return
		}
		h.renderImageVariantForm(w, r, form, ToWebImageVariant(currentVariantResponse.ImageVariant), "Validation failed", http.StatusBadRequest)
		return
	}

	imageID, parseUUIDErr := uuid.Parse(form.ImageID)
	if parseUUIDErr != nil {
		h.Err(w, parseUUIDErr, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Fetch parent image to get path, url, mimetype, size, width, height
	var imageResponse struct {
		Image feat.Image `json:"image"`
	}
	imagePath := fmt.Sprintf("/ssg/images/%s", imageID)
	apiGetErr := h.apiClient.Get(r, imagePath, &imageResponse)
	if apiGetErr != nil {
		h.Err(w, apiGetErr, "Cannot get parent image from API", http.StatusInternalServerError)
		return
	}
	parentImage := imageResponse.Image

	featImageVariant := feat.NewImageVariant()
	featImageVariant.ID = uuid.MustParse(form.ID)
	featImageVariant.ImageID = imageID
	featImageVariant.Kind = form.Name
	featImageVariant.BlobRef = "/static/images/" + uuid.New().String() + "." + parentImage.Mime // Use BlobRef and generate a new unique name
	featImageVariant.Mime = parentImage.Mime
	featImageVariant.FilesizeByte = parentImage.FilesizeByte
	featImageVariant.Width = parentImage.Width
	featImageVariant.Height = parentImage.Height

	path := fmt.Sprintf("/ssg/images/%s/variants/%s", featImageVariant.ImageID, featImageVariant.GetID()) // Use featImageVariant.ImageID
	apiPutErr := h.apiClient.Put(r, path, featImageVariant, nil)
	if apiPutErr != nil {
		h.Err(w, apiPutErr, "Failed to update image variant via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Image variant updated successfully")
	h.Redir(w, r, hm.ListPath(&ImageVariant{}), http.StatusSeeOther)
}

// ListImageVariants lists all image variants for a given image.
func (h *WebHandler) ListImageVariants(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List image variants")

	imageIDStr := r.URL.Query().Get("imageID")
	if imageIDStr == "" {
		h.Err(w, nil, "Missing image ID for variant listing", http.StatusBadRequest)
		return
	}

	var response struct {
		ImageVariants []feat.ImageVariant `json:"imageVariants"`
	}
	path := fmt.Sprintf("/ssg/images/%s/variants", imageIDStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get image variants from API", http.StatusInternalServerError)
		return
	}
	variants := ToWebImageVariants(response.ImageVariants)

	page := hm.NewPage(r, variants)
	page.Form.SetAction(ssgPath)
	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(&ImageVariant{})

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-image-variants")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

// ShowImageVariant displays a specific image variant.
func (h *WebHandler) ShowImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show image variant")

	imageIDStr := r.URL.Query().Get("imageID")
	idStr := r.URL.Query().Get("id")
	if imageIDStr == "" || idStr == "" {
		h.Err(w, nil, "Missing image ID or variant ID", http.StatusBadRequest)
		return
	}

	var response struct {
		ImageVariant feat.ImageVariant `json:"imageVariant"`
	}
	path := fmt.Sprintf("/ssg/images/%s/variants/%s", imageIDStr, idStr)
	err := h.apiClient.Get(r, path, &response)
	if err != nil {
		h.Err(w, err, "Cannot get image variant from API", http.StatusInternalServerError)
		return
	}

	variant := ToWebImageVariant(response.ImageVariant)

	page := hm.NewPage(r, variant)
	page.Name = "Show Image Variant"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&variant, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-image-variant")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, http.StatusOK)
}

// DeleteImageVariant handles the deletion of an image variant.
func (h *WebHandler) DeleteImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete image variant")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	imageIDStr := r.Form.Get("imageID")
	idStr := r.Form.Get("id")
	if imageIDStr == "" || idStr == "" {
		h.Err(w, nil, "Missing image ID or variant ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/ssg/images/%s/variants/%s", imageIDStr, idStr)
	err := h.apiClient.Delete(r, path)
	if err != nil {
		h.Err(w, err, "Failed to delete image variant via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Image variant deleted successfully")
	h.Redir(w, r, hm.ListPath(&ImageVariant{}), http.StatusSeeOther)
}

func (h *WebHandler) renderImageVariantForm(w http.ResponseWriter, r *http.Request, form ImageVariantForm, variant ImageVariant, errorMessage string, statusCode int) {
	page := hm.NewPage(r, variant)
	page.SetForm(&form)

	if variant.IsZero() {
		page.Name = "New Image Variant"
		page.IsNew = true
		page.Form.SetAction(hm.CreatePath(&ImageVariant{}))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Image Variant"
		page.IsNew = false
		page.Form.SetAction(hm.UpdatePath(&ImageVariant{}))
		page.Form.SetSubmitButtonText("Update")
	}

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(&variant)

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-image-variant")
	if err != nil {
		h.Err(w, err, hm.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	page.SetFlash(h.GetFlash(r))

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, hm.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, statusCode)
}
