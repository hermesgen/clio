package ssg

import (
	"time"

	"github.com/google/uuid"

	"github.com/hermesgen/hm"
)

// Section model.
type Section struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	ShortID string    `json:"-" db:"short_id"`
	ref     string    `json:"-"`

	// Section specific fields
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Path        string    `json:"path" db:"path"`
	LayoutID    uuid.UUID `json:"layout_id" db:"layout_id"`
	LayoutName  string    `json:"layout_name" db:"layout_name"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
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
	return "section"
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

// GenCreateValues delegates to the functional helper.
func (s *Section) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(s, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (s *Section) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(s, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (s *Section) GetCreatedBy() uuid.UUID {
	return s.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (s *Section) GetUpdatedBy() uuid.UUID {
	return s.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (s *Section) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (s *Section) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (s *Section) SetCreatedAt(t time.Time) {
	s.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (s *Section) SetUpdatedAt(t time.Time) {
	s.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (s *Section) SetCreatedBy(u uuid.UUID) {
	s.CreatedBy = u
}

// SetUpdatedBy implements the Auditable interface.
func (s *Section) SetUpdatedBy(u uuid.UUID) {
	s.UpdatedBy = u
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

func (s *Section) Ref() string {
	return s.ref
}

func (s *Section) SetRef(ref string) {
	s.ref = ref
}
