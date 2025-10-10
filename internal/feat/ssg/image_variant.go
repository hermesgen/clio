package ssg

import (
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

// ImageVariant represents a specific rendition of an image.
type ImageVariant struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	ShortID string    `json:"-" db:"short_id"`
	ref     string    `json:"-"`

	ImageID      uuid.UUID `json:"image_id" db:"image_id"`
	Kind         string    `json:"kind" db:"kind"` // e.g., 'original', 'web', 'thumb'
	Width        int       `json:"width" db:"width"`
	Height       int       `json:"height" db:"height"`
	FilesizeByte int64     `json:"filesize_bytes" db:"filesize_bytes"`
	Mime         string    `json:"mime" db:"mime"`
	BlobRef      string    `json:"blob_ref" db:"blob_ref"` // Abstract reference to the stored file

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewImageVariant creates a new ImageVariant instance with default values.
func NewImageVariant() ImageVariant {
	return ImageVariant{
		ID: uuid.New(),
	}
}

// Type returns the type of the entity.
func (iv *ImageVariant) Type() string {
	return "image-variant"
}

// GetID returns the unique identifier of the entity.
func (iv ImageVariant) GetID() uuid.UUID {
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

// GenCreateValues delegates to the functional helper.
func (iv *ImageVariant) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(iv, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (iv *ImageVariant) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(iv, userID...)
}

// GetCreatedBy returns the UUID of the user who created the entity.
func (iv *ImageVariant) GetCreatedBy() uuid.UUID {
	return iv.CreatedBy
}

// GetUpdatedBy returns the UUID of the user who last updated the entity.
func (iv *ImageVariant) GetUpdatedBy() uuid.UUID {
	return iv.UpdatedBy
}

// GetCreatedAt returns the creation time of the entity.
func (iv *ImageVariant) GetCreatedAt() time.Time {
	return iv.CreatedAt
}

// GetUpdatedAt returns the last update time of the entity.
func (iv *ImageVariant) GetUpdatedAt() time.Time {
	return iv.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (iv *ImageVariant) SetCreatedAt(t time.Time) {
	iv.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (iv *ImageVariant) SetUpdatedAt(t time.Time) {
	iv.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (iv *ImageVariant) SetCreatedBy(id uuid.UUID) {
	iv.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (iv *ImageVariant) SetUpdatedBy(id uuid.UUID) {
	iv.UpdatedBy = id
}

// IsZero returns true if the ImageVariant is uninitialized.
func (iv *ImageVariant) IsZero() bool {
	return iv.ID == uuid.Nil
}

// Slug returns a slug for the image variant.
func (iv *ImageVariant) Slug() string {
	return hm.Normalize(iv.Kind) + "-" + iv.GetShortID()
}

func (iv *ImageVariant) Ref() string {
	return iv.ref
}

func (iv *ImageVariant) SetRef(ref string) {
	iv.ref = ref
}
