package ssg

import (
	"github.com/google/uuid"

	"github.com/hermesgen/hm"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

const (
	tagType = "tag"
)

// Tag model for the web layer.
type Tag struct {
	ID        uuid.UUID `json:"id"`
	ShortID   string    `json:"-"`
	Name      string    `json:"name"`
	SlugField string    `json:"slug"`
}

// NewTag creates a new Tag for the web layer.
func NewTag(name string) Tag {
	return Tag{
		Name: name,
	}
}

// Type returns the type of the entity.
func (t *Tag) Type() string {
	return hm.DefaultType(tagType)
}

// GetID returns the unique identifier of the entity.
func (t *Tag) GetID() uuid.UUID {
	return t.ID
}

// GenID delegates to the functional helper.
func (t *Tag) GenID() {
	hm.GenID(t)
}

// SetID sets the unique identifier of the entity.
func (t *Tag) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if t.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		t.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (t *Tag) GetShortID() string {
	return t.ShortID
}

// GenShortID delegates to the functional helper.
func (t *Tag) GenShortID() {
	hm.GenShortID(t)
}

// SetShortID sets the short ID of the entity.
func (t *Tag) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if t.ShortID == "" || shouldForce {
		t.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (t *Tag) TypeID() string {
	return hm.Normalize(t.Type()) + "-" + t.GetShortID()
}

// IsZero returns true if the Tag is uninitialized.
func (t *Tag) IsZero() bool {
	return t.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (t *Tag) Slug() string {
	if t.SlugField != "" {
		return t.SlugField
	}
	return hm.Normalize(t.Name) + "-" + t.GetShortID()
}

func (t *Tag) OptValue() string {
	return t.GetID().String()
}

func (t *Tag) OptLabel() string {
	return t.Name
}

// ToWebTag converts a feat.Tag model to a web.Tag model.
func ToWebTag(featTag feat.Tag) Tag {
	return Tag{
		ID:        featTag.ID,
		ShortID:   featTag.ShortID,
		Name:      featTag.Name,
		SlugField: featTag.Slug(), // Assuming feat.Tag has a Slug() method
	}
}

// ToWebTags converts a slice of feat.Tag models to a slice of web.Tag models.
func ToWebTags(featTags []feat.Tag) []Tag {
	webTags := make([]Tag, len(featTags))
	for i, featTag := range featTags {
		webTags[i] = ToWebTag(featTag)
	}
	return webTags
}
