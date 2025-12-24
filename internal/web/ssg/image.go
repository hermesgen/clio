package ssg

import (
	"mime/multipart"
	"net/http" // Import http

	"github.com/google/uuid"

	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

const (
	imageType = "image"
)

// Image model for the web layer.
type Image struct {
	ID          uuid.UUID `json:"id"`
	ShortID     string    `json:"-"`
	Name        string    `json:"name"`        // Maps to feat.Image.Title
	Description string    `json:"description"` // Maps to feat.Image.LongDescription
	Path        string    `json:"path"`        // From ImageVariant
	URL         string    `json:"url"`         // From ImageVariant
	AltText     string    `json:"altText"`
	MimeType    string    `json:"mimeType"`
	Size        int64     `json:"size"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
}

// NewImage creates a new Image for the web layer.
func NewImage(name, description, path, url, altText, mimeType string, size int64, width, height int) Image {
	return Image{
		Name:        name,
		Description: description,
		Path:        path,
		URL:         url,
		AltText:     altText,
		MimeType:    mimeType,
		Size:        size,
		Width:       width,
		Height:      height,
	}
}

// Type returns the type of the entity.
func (i *Image) Type() string {
	return hm.DefaultType(imageType)
}

// GetID returns the unique identifier of the entity.
func (i *Image) GetID() uuid.UUID {
	return i.ID
}

// GenID delegates to the functional helper.
func (i *Image) GenID() {
	hm.GenID(i)
}

// SetID sets the unique identifier of the entity.
func (i *Image) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if i.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		i.ID = id
	}
}

// GetShortID returns the short ID portion of the slug.
func (i *Image) GetShortID() string {
	return i.ShortID
}

// GenShortID delegates to the functional helper.
func (i *Image) GenShortID() {
	hm.GenShortID(i)
}

// SetShortID sets the short ID of the entity.
func (i *Image) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if i.ShortID == "" || shouldForce {
		i.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (i *Image) TypeID() string {
	return hm.Normalize(i.Type()) + "-" + i.GetShortID()
}

// IsZero returns true if the Image is uninitialized.
func (i *Image) IsZero() bool {
	return i.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (i *Image) Slug() string {
	return hm.Normalize(i.Name) + "-" + i.GetShortID()
}

func (i *Image) OptValue() string {
	return i.GetID().String()
}

func (i *Image) OptLabel() string {
	return i.Name
}

// StringID returns the unique identifier of the entity as a string.
func (i *Image) StringID() string {
	return i.GetID().String()
}

// ToWebImage converts a feat.Image model to a web.Image model.
func ToWebImage(featImage feat.Image) Image {
	return Image{
		ID:      featImage.ID,
		ShortID: featImage.ShortID,
		Name:    featImage.Title, // Map Title to Name
		// Path and URL are not directly in feat.Image, they come from variants
		AltText: featImage.AltText,
		Width:   featImage.Width,
		Height:  featImage.Height,
	}
}

// ToWebImages converts a slice of feat.Image models to a slice of web.Image models.
func ToWebImages(featImages []feat.Image) []Image {
	webImages := make([]Image, len(featImages))
	for i, featImage := range featImages {
		webImages[i] = ToWebImage(featImage)
	}
	return webImages
}

// ImageForm represents the form data for an Image.
type ImageForm struct {
	*hm.BaseForm                       // Embed BaseForm
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	AltText      string                `json:"altText"`
	File         *multipart.FileHeader `json:"file"` // For file upload
}

// NewImageForm creates a new ImageForm.
func NewImageForm(r *http.Request) ImageForm {
	return ImageForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// ToFeatImage converts an ImageForm to a feat.Image.
func ToFeatImage(form ImageForm) feat.Image {
	id, _ := uuid.Parse(form.ID)
	return feat.Image{
		ID:      id,
		Title:   form.Name, // Map Name to Title
		AltText: form.AltText,
		// Path, URL, MimeType, Size, Width, Height are set in webhandlerimage.go from file upload
	}
}

// ToImageForm converts a web.Image to an ImageForm.
func ToImageForm(r *http.Request, image Image) ImageForm {
	form := NewImageForm(r)
	form.ID = image.ID.String()
	form.Name = image.Name
	form.Description = image.Description
	form.AltText = image.AltText
	return form
}

// Validate validates the ImageForm.
func (f *ImageForm) Validate() {
	// Reset validation errors
	f.BaseForm.SetValidation(&hm.Validation{})

	validation := f.BaseForm.Validation() // Get the pointer to validation

	if f.Name == "" {
		validation.AddFieldError("name", f.Name, "Name cannot be empty")
	}

	if f.File == nil || f.File.Size == 0 {
		validation.AddFieldError("file", "", "Image file is required")
	}
}
