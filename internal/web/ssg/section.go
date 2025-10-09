package ssg

import (
	"github.com/google/uuid"

	"github.com/hermesgen/hm"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

const (
	sectionType = "section"
)

// Section model.
type Section struct {
	ID          uuid.UUID `json:"id"`
	ShortID     string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	LayoutID    uuid.UUID `json:"layout_id"`
	Header      string    `json:"header"`
	BlogHeader  string    `json:"blog_header"`
	LayoutName  string    `json:"layout_name"`
}

// NewSection creates a new Section.
func NewSection(name, description, path string, layoutID uuid.UUID) Section {
	s := Section{
		Name:        name,
		Description: description,
		Path:        path,
		LayoutID:    layoutID,
	}

	return s
}

// Type returns the type of the entity.
func (s *Section) Type() string {
	return hm.DefaultType(sectionType)
}

// GetID returns the unique identifier of the entity.
func (s *Section) GetID() uuid.UUID {
	return s.ID
}

// GenID delegates to the functional helper.
func (s *Section) GenID() {
	hm.GenID(s)
}

// SetID sets the unique identifier of the entity.
func (s *Section) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if s.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		s.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (s *Section) GetShortID() string {
	return s.ShortID
}

// GenShortID delegates to the functional helper.
func (s *Section) GenShortID() {
	hm.GenShortID(s)
}

// SetShortID sets the short ID of the entity.
func (s *Section) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if s.ShortID == "" || shouldForce {
		s.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (s *Section) TypeID() string {
	return hm.Normalize(s.Type()) + "-" + s.GetShortID()
}

// IsZero returns true if the Section is uninitialized.
func (s *Section) IsZero() bool {
	return s.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (s *Section) Slug() string {
	return hm.Normalize(s.Name) + "-" + s.GetShortID()
}

func (s *Section) OptValue() string {
	return s.GetID().String()
}

func (s *Section) OptLabel() string {
	return s.Name
}

// ToWebSection converts a feat.Section model to a web.Section model.
// NOTE: Probably we want to avoid this coupling in the future.
func ToWebSection(featSection feat.Section) Section {
	return Section{
		ID:          featSection.ID,
		ShortID:     featSection.ShortID,
		Name:        featSection.Name,
		Description: featSection.Description,
		Path:        featSection.Path,
		LayoutID:    featSection.LayoutID,
		Header:      "", // TODO: Get header via relationship
		BlogHeader:  "", // TODO: Get blog header via relationship
		LayoutName:  featSection.LayoutName,
	}
}

// ToWebSections converts a slice of feat.Section models to a slice of web.Section models.
// NOTE: Prpbably we want to avoid this coupling in the future.
func ToWebSections(featSections []feat.Section) []Section {
	webSections := make([]Section, len(featSections))
	for i, featSection := range featSections {
		webSections[i] = ToWebSection(featSection)
	}
	return webSections
}
