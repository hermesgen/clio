package ssg

import (
	"time"

	"github.com/google/uuid"
)

// ContentImage represents the relationship between content and images
type ContentImage struct {
	// Common
	ID         uuid.UUID `json:"id" db:"id"`
	ContentID  uuid.UUID `json:"content_id" db:"content_id"`
	ImageID    uuid.UUID `json:"image_id" db:"image_id"`
	IsHeader   bool      `json:"is_header" db:"is_header"`
	IsFeatured bool      `json:"is_featured" db:"is_featured"`
	OrderNum   int       `json:"order_num" db:"order_num"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	// Related objects (populated by joins)
	Image   *Image   `json:"image,omitempty" db:"-"`
	Content *Content `json:"content,omitempty" db:"-"`
}

// NewContentImage creates a new ContentImage
func NewContentImage(contentID, imageID uuid.UUID, isHeader bool) *ContentImage {
	now := time.Now()
	return &ContentImage{
		ID:         uuid.New(),
		ContentID:  contentID,
		ImageID:    imageID,
		IsHeader:   isHeader,
		IsFeatured: false,
		OrderNum:   0,
		CreatedAt:  now,
	}
}

// SectionImage represents the relationship between sections and images
type SectionImage struct {
	// Common
	ID         uuid.UUID `json:"id" db:"id"`
	SectionID  uuid.UUID `json:"section_id" db:"section_id"`
	ImageID    uuid.UUID `json:"image_id" db:"image_id"`
	IsHeader   bool      `json:"is_header" db:"is_header"`
	IsFeatured bool      `json:"is_featured" db:"is_featured"`
	OrderNum   int       `json:"order_num" db:"order_num"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	// Related objects (populated by joins)
	Image *Image `json:"image,omitempty" db:"-"`
}

// NewSectionImage creates a new SectionImage
func NewSectionImage(sectionID uuid.UUID, imageID uuid.UUID, isHeader bool) *SectionImage {
	now := time.Now()
	return &SectionImage{
		ID:         uuid.New(),
		SectionID:  sectionID,
		ImageID:    imageID,
		IsHeader:   isHeader,
		IsFeatured: false,
		OrderNum:   0,
		CreatedAt:  now,
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
