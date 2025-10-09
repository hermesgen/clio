package ssg

import (
	"github.com/google/uuid"

	"github.com/hermesgen/hm"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

const (
	layoutType = "layout"
)

// Layout model.
type Layout struct {
	// Common
	ID          uuid.UUID `json:"id"`
	ShortID     string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
}

// Newlayout creates a new Layout.
func Newlayout(name, description, code string) Layout {
	l := Layout{
		Name:        name,
		Description: description,
		Code:        code,
	}

	return l
}

// Type returns the type of the entity.
func (l *Layout) Type() string {
	return hm.DefaultType(layoutType)
}

// GetID returns the unique identifier of the entity.
func (l *Layout) GetID() uuid.UUID {
	return l.ID
}

// GenID delegates to the functional helper.
func (l *Layout) GenID() {
	hm.GenID(l)
}

// SetID sets the unique identifier of the entity.
func (l *Layout) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if l.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		l.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (l *Layout) GetShortID() string {
	return l.ShortID
}

// GenShortID delegates to the functional helper.
func (l *Layout) GenShortID() {
	hm.GenShortID(l)
}

// SetShortID sets the short ID of the entity.
func (l *Layout) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if l.ShortID == "" || shouldForce {
		l.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (l *Layout) TypeID() string {
	return hm.Normalize(l.Type()) + "-" + l.GetShortID()
}

// IsZero returns true if the Layout is uninitialized.
func (l *Layout) IsZero() bool {
	return l.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (l *Layout) Slug() string {
	return hm.Normalize(l.Name) + "-" + l.GetShortID()
}

func (l *Layout) OptValue() string {
	return l.GetID().String()
}

func (l *Layout) OptLabel() string {
	return l.Name
}

// StringID returns the unique identifier of the entity as a string.
func (l *Layout) StringID() string {
	return l.GetID().String()
}

// ToWebLayout converts a feat.Layout model to a web.Layout model.
func ToWebLayout(featLayout feat.Layout) Layout {
	return Layout{
		ID:          featLayout.ID,
		ShortID:     featLayout.ShortID,
		Name:        featLayout.Name,
		Description: featLayout.Description,
		Code:        featLayout.Code,
	}
}

// ToWebLayouts converts a slice of feat.Layout models to a slice of web.Layout models.
func ToWebLayouts(featLayouts []feat.Layout) []Layout {
	webLayouts := make([]Layout, len(featLayouts))
	for i, featLayout := range featLayouts {
		webLayouts[i] = ToWebLayout(featLayout)
	}
	return webLayouts
}
