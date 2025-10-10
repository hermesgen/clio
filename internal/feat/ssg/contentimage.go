package ssg

import (
	"time"

	"github.com/google/uuid"
)

// ContentImage represents the relationship between content and images
type ContentImage struct {
	// Common
	ID        uuid.UUID `json:"id" db:"id"`
	ContentID uuid.UUID `json:"content_id" db:"content_id"`
	ImageID   uuid.UUID `json:"image_id" db:"image_id"`
	Purpose   string    `json:"purpose" db:"purpose"` // 'header', 'content', 'thumbnail'
	Position  int       `json:"position" db:"position"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Related objects (populated by joins)
	Image   *Image   `json:"image,omitempty" db:"-"`
	Content *Content `json:"content,omitempty" db:"-"`
}

// NewContentImage creates a new ContentImage
func NewContentImage(contentID, imageID uuid.UUID, purpose string) *ContentImage {
	now := time.Now()
	return &ContentImage{
		ID:        uuid.New(),
		ContentID: contentID,
		ImageID:   imageID,
		Purpose:   purpose,
		Position:  0,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SectionImage represents the relationship between sections and images
type SectionImage struct {
	// Common
	ID        uuid.UUID `json:"id" db:"id"`
	SectionID uuid.UUID `json:"section_id" db:"section_id"`
	ImageID   uuid.UUID `json:"image_id" db:"image_id"`
	Purpose   string    `json:"purpose" db:"purpose"` // 'header', 'blog_header'
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Related objects (populated by joins)
	Image *Image `json:"image,omitempty" db:"-"`
}

// NewSectionImage creates a new SectionImage
func NewSectionImage(sectionID uuid.UUID, imageID uuid.UUID, purpose string) *SectionImage {
	now := time.Now()
	return &SectionImage{
		ID:        uuid.New(),
		SectionID: sectionID,
		ImageID:   imageID,
		Purpose:   purpose,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ImagePurpose constants
const (
	ImagePurposeHeader     = "header"
	ImagePurposeContent    = "content"
	ImagePurposeThumbnail  = "thumbnail"
	ImagePurposeBlogHeader = "blog_header"
)

// LayoutType constants
const (
	LayoutTypeSection = "section"
	LayoutTypeLayout  = "layout"
)
