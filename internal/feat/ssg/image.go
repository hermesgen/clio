package ssg

import (
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

// Image represents an image asset with its metadata.
type Image struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	ShortID string    `json:"-" db:"short_id"`
	ref     string    `json:"-"`

	// Site relationship
	SiteID uuid.UUID `json:"site_id" db:"site_id"`

	// File information
	FileName string `json:"file_name" db:"file_name"`
	FilePath string `json:"file_path" db:"file_path"`
	Width    int    `json:"width" db:"width"`
	Height   int    `json:"height" db:"height"`

	// Accessibility fields
	Title   string `json:"title" db:"title"`
	AltText string `json:"alt_text" db:"alt_text"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewImage creates a new Image instance with default values.
func NewImage() Image {
	return Image{
		ID: uuid.New(),
	}
}

// Type returns the type of the entity.
func (i *Image) Type() string {
	return "image"
}

// GetID returns the unique identifier of the entity.
func (i Image) GetID() uuid.UUID {
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

// GenCreateValues delegates to the functional helper.
func (i *Image) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(i, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (i *Image) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(i, userID...)
}

// GetCreatedBy returns the UUID of the user who created the entity.
func (i *Image) GetCreatedBy() uuid.UUID {
	return i.CreatedBy
}

// GetUpdatedBy returns the UUID of the user who last updated the entity.
func (i *Image) GetUpdatedBy() uuid.UUID {
	return i.UpdatedBy
}

// GetCreatedAt returns the creation time of the entity.
func (i *Image) GetCreatedAt() time.Time {
	return i.CreatedAt
}

// GetUpdatedAt returns the last update time of the entity.
func (i *Image) GetUpdatedAt() time.Time {
	return i.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (i *Image) SetCreatedAt(t time.Time) {
	i.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (i *Image) SetUpdatedAt(t time.Time) {
	i.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (i *Image) SetCreatedBy(id uuid.UUID) {
	i.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (i *Image) SetUpdatedBy(id uuid.UUID) {
	i.UpdatedBy = id
}

// IsZero returns true if the Image is uninitialized.
func (i *Image) IsZero() bool {
	return i.ID == uuid.Nil
}

// Slug returns a slug for the image.
func (i *Image) Slug() string {
	if i.Title != "" {
		return hm.Normalize(i.Title) + "-" + i.GetShortID()
	}
	return hm.Normalize(i.FileName) + "-" + i.GetShortID()
}

func (i *Image) Ref() string {
	return i.ref
}

func (i *Image) SetRef(ref string) {
	i.ref = ref
}
