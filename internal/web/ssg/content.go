package ssg

import (
	"time"

	"github.com/google/uuid"

	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

const (
	contentType = "content"
)

// Content model.
type Content struct {
	ID          uuid.UUID  `json:"id"`
	ShortID     string     `json:"-"`
	UserID      uuid.UUID  `json:"user_id"`
	SectionID   uuid.UUID  `json:"section_id"`
	Heading     string     `json:"heading"`
	Body        string     `json:"body"`
	Image       string     `json:"image"`
	Draft       bool       `json:"draft"`
	Featured    bool       `json:"featured"`
	PublishedAt *time.Time `json:"published_at"`
	Tags        []feat.Tag `json:"tags"`
	Meta        feat.Meta  `json:"meta"`
	SectionPath string     `json:"section_path,omitempty"`
	SectionName string     `json:"section_name,omitempty"`
}

// NewContent creates a new Content.
func NewContent(heading, body string) Content {
	c := Content{
		Heading: heading,
		Body:    body,
	}

	return c
}

// Type returns the type of the entity.
func (c *Content) Type() string {
	return hm.DefaultType(contentType)
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

// TypeID returns a universal identifier for a specific model instance.
func (c *Content) TypeID() string {
	return hm.Normalize(c.Type()) + "-" + c.GetShortID()
}

// IsZero returns true if the Content is uninitialized.
func (c *Content) IsZero() bool {
	return c.ID == uuid.Nil
}

// Slug returns the slug for the content.
func (c *Content) Slug() string {
	return "content"
}

func (c *Content) OptValue() string {
	return c.GetID().String()
}

func (c *Content) OptLabel() string {
	return c.Heading
}

// ToWebContent converts a feat.Content model to a web.Content model.
func ToWebContent(featContent feat.Content) Content {
	return Content{
		ID:          featContent.ID,
		ShortID:     featContent.ShortID,
		UserID:      featContent.UserID,
		SectionID:   featContent.SectionID,
		Heading:     featContent.Heading,
		Body:        featContent.Body,
		Image:       "", // TODO: Get image via relationship
		Draft:       featContent.Draft,
		Featured:    featContent.Featured,
		PublishedAt: featContent.PublishedAt,
		Tags:        featContent.Tags,
		Meta:        featContent.Meta,
		SectionPath: featContent.SectionPath,
		SectionName: featContent.SectionName,
	}
}

// ToWebContents converts a slice of feat.Content models to a slice of web.Content models.
func ToWebContents(featContents []feat.Content) []Content {
	webContents := make([]Content, len(featContents))
	for i, featContent := range featContents {
		webContents[i] = ToWebContent(featContent)
	}
	return webContents
}
