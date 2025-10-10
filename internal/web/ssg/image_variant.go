package ssg

import (
	"net/http" // Import http

	"github.com/google/uuid"

	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

const (
	imageVariantType = "imageVariant"
)

// ImageVariant model for the web layer.
type ImageVariant struct {
	ID       uuid.UUID `json:"id"`
	ShortID  string    `json:"-"`
	ImageID  uuid.UUID `json:"imageID"`
	Name     string    `json:"name"` // Maps to feat.ImageVariant.Kind
	Path     string    `json:"path"` // Maps to feat.ImageVariant.BlobRef
	URL      string    `json:"url"`  // Derived from BlobRef
	MimeType string    `json:"mimeType"`
	Size     int64     `json:"size"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
}

// NewImageVariant creates a new ImageVariant for the web layer.
func NewImageVariant(imageID uuid.UUID, name, path, url, mimeType string, size int64, width, height int) ImageVariant {
	return ImageVariant{
		ImageID:  imageID,
		Name:     name,
		Path:     path,
		URL:      url,
		MimeType: mimeType,
		Size:     size,
		Width:    width,
		Height:   height,
	}
}

// Type returns the type of the entity.
func (iv *ImageVariant) Type() string {
	return hm.DefaultType(imageVariantType)
}

// GetID returns the unique identifier of the entity.
func (iv *ImageVariant) GetID() uuid.UUID {
	return iv.ID
}

// GenID delegates to the functional helper.
func (iv *ImageVariant) GenID() {
	hm.GenID(iv)
}

// SetID sets the unique identifier of the entity.
func (iv *ImageVariant) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if iv.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		iv.ID = id
	}
}

// GetShortID returns the short ID portion of the slug.
func (iv *ImageVariant) GetShortID() string {
	return iv.ShortID
}

// GenShortID delegates to the functional helper.
func (iv *ImageVariant) GenShortID() {
	hm.GenShortID(iv)
}

// SetShortID sets the short ID of the entity.
func (iv *ImageVariant) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if iv.ShortID == "" || shouldForce {
		iv.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (iv *ImageVariant) TypeID() string {
	return hm.Normalize(iv.Type()) + "-" + iv.GetShortID()
}

// IsZero returns true if the ImageVariant is uninitialized.
func (iv *ImageVariant) IsZero() bool {
	return iv.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (iv *ImageVariant) Slug() string {
	return hm.Normalize(iv.Name) + "-" + iv.GetShortID()
}

func (iv *ImageVariant) OptValue() string {
	return iv.GetID().String()
}

func (iv *ImageVariant) OptLabel() string {
	return iv.Name
}

// StringID returns the unique identifier of the entity as a string.
func (iv *ImageVariant) StringID() string {
	return iv.GetID().String()
}

// ToWebImageVariant converts a feat.ImageVariant model to a web.ImageVariant model.
func ToWebImageVariant(featImageVariant feat.ImageVariant) ImageVariant {
	return ImageVariant{
		ID:       featImageVariant.ID,
		ShortID:  featImageVariant.ShortID,
		ImageID:  featImageVariant.ImageID,
		Name:     featImageVariant.Kind,    // Map Kind to Name
		Path:     featImageVariant.BlobRef, // Map BlobRef to Path
		URL:      featImageVariant.BlobRef, // Assuming BlobRef can be used as URL for now
		MimeType: featImageVariant.Mime,
		Size:     featImageVariant.FilesizeByte,
		Width:    featImageVariant.Width,
		Height:   featImageVariant.Height,
	}
}

// ToWebImageVariants converts a slice of feat.ImageVariant models to a slice of web.ImageVariant models.
func ToWebImageVariants(featImageVariants []feat.ImageVariant) []ImageVariant {
	webImageVariants := make([]ImageVariant, len(featImageVariants))
	for i, featImageVariant := range featImageVariants {
		webImageVariants[i] = ToWebImageVariant(featImageVariant)
	}
	return webImageVariants
}

// ImageVariantForm represents the form data for an ImageVariant.
type ImageVariantForm struct {
	*hm.BaseForm        // Embed BaseForm
	ID           string `json:"id"`
	ImageID      string `json:"imageID"`
	Name         string `json:"name"` // Maps to feat.ImageVariant.Kind
}

// NewImageVariantForm creates a new ImageVariantForm.
func NewImageVariantForm(r *http.Request) ImageVariantForm {
	return ImageVariantForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// ToFeatImageVariant converts an ImageVariantForm to a feat.ImageVariant.
func ToFeatImageVariant(form ImageVariantForm) feat.ImageVariant {
	id, _ := uuid.Parse(form.ID)
	imageID, _ := uuid.Parse(form.ImageID)
	return feat.ImageVariant{
		ID:      id,
		ImageID: imageID,
		Kind:    form.Name, // Map Name to Kind
	}
}

// ToImageVariantForm converts a web.ImageVariant to an ImageVariantForm.
func ToImageVariantForm(r *http.Request, imageVariant ImageVariant) ImageVariantForm {
	form := NewImageVariantForm(r)
	form.ID = imageVariant.ID.String()
	form.ImageID = imageVariant.ImageID.String()
	form.Name = imageVariant.Name
	return form
}

// Validate validates the ImageVariantForm.
func (f *ImageVariantForm) Validate() {
	// Reset validation errors
	f.BaseForm.SetValidation(&hm.Validation{})

	validation := f.BaseForm.Validation() // Get the pointer to validation

	if f.ImageID == "" {
		validation.AddFieldError("imageID", f.ImageID, "Image ID cannot be empty")
	}

	if f.Name == "" {
		validation.AddFieldError("name", f.Name, "Name cannot be empty")
	}
}
