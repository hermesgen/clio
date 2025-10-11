package ssg

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hermesgen/hm"
)

type Content struct {
	ID      uuid.UUID `json:"id" db:"id"`
	ShortID string    `json:"-" db:"short_id"`
	ref     string    `json:"-"`

	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	SectionID   uuid.UUID  `json:"section_id" db:"section_id"`
	Kind        string     `json:"kind" db:"kind"`
	Heading     string     `json:"heading" db:"heading"`
	Summary     string     `json:"summary" db:"summary"`
	Body        string     `json:"body" db:"body"`
	Draft       bool       `json:"draft" db:"draft"`
	Featured    bool       `json:"featured" db:"featured"`
	Series      string     `json:"series,omitempty" db:"series"`
	SeriesOrder int        `json:"series_order,omitempty" db:"series_order"`
	PublishedAt *time.Time `json:"published_at" db:"published_at"`
	Tags        []Tag      `json:"tags"`
	Meta        Meta       `json:"meta"`

	ThumbnailURL       string `json:"thumbnail_url,omitempty" db:"-"`
	HeaderImageURL     string `json:"header_image_url,omitempty" db:"-"`
	HeaderImageAlt     string `json:"header_image_alt,omitempty" db:"-"`
	HeaderImageCaption string `json:"header_image_caption,omitempty" db:"-"`

	SectionPath string `json:"section_path,omitempty" db:"section_path"`
	SectionName string `json:"section_name,omitempty" db:"section_name"`

	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewContent creates a new Content.
func NewContent(heading, body string) Content {
	c := Content{
		Heading: heading,
		Body:    body,
		Draft:   true,
	}

	return c
}

// Type returns the type of the entity.
func (c *Content) Type() string {
	return "content"
}

// GetID returns the unique identifier of the entity.
func (c *Content) GetID() uuid.UUID {
	return c.ID
}

// GenID delegates to the functional helper.
func (c *Content) GenID() {
	hm.GenID(c)
}

// SetID sets the unique identifier of the entity.
func (c *Content) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if c.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		c.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (c *Content) GetShortID() string {
	return c.ShortID
}

// GenShortID delegates to the functional helper.
func (c *Content) GenShortID() {
	hm.GenShortID(c)
}

// SetShortID sets the short ID of the entity.
func (c *Content) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if c.ShortID == "" || shouldForce {
		c.ShortID = shortID
	}
}

// GenCreateValues delegates to the functional helper.
func (c *Content) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(c, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (c *Content) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(c, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (c *Content) GetCreatedBy() uuid.UUID {
	return c.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (c *Content) GetUpdatedBy() uuid.UUID {
	return c.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (c *Content) GetCreatedAt() time.Time {
	return c.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (c *Content) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (c *Content) SetCreatedAt(t time.Time) {
	c.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (c *Content) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (c *Content) SetCreatedBy(u uuid.UUID) {
	c.CreatedBy = u
}

// SetUpdatedBy implements the Auditable interface.
func (c *Content) SetUpdatedBy(u uuid.UUID) {
	c.UpdatedBy = u
}

// IsZero returns true if the Content is uninitialized.
func (c *Content) IsZero() bool {
	return c.ID == uuid.Nil
}

// Slug returns the slug for the content.
func (c *Content) Slug() string {
	return hm.Normalize(c.Heading) + "-" + c.GetShortID()
}

func (c *Content) OptValue() string {
	return c.GetID().String()
}

func (c *Content) OptLabel() string {
	return c.Heading
}

func (c *Content) Ref() string {
	return c.ref
}

func (c *Content) SetRef(ref string) {
	c.ref = ref
}

// String implements the fmt.Stringer interface to provide a clean log representation
// that truncates the Body field to prevent cluttering logs with long markdown content
func (c Content) String() string {
	bodyPreview := c.Body
	if len(bodyPreview) > 50 {
		bodyPreview = bodyPreview[:47] + "..."
	}
	return fmt.Sprintf("Content{ID: %s, Heading: %q, Body: %q, Draft: %t}", 
		c.ID.String()[:8]+"...", c.Heading, bodyPreview, c.Draft)
}
