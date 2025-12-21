package ssg

import (
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

// Site represents a site instance in the multi-site system.
// Site records are stored in sites.db (metadata only).
// Each site has its own complete database in sites/{slug}/db/clio.db
type Site struct {
	ID      uuid.UUID `json:"id" db:"id"`
	ShortID string    `json:"-" db:"short_id"`
	ref     string    `json:"-"`

	Name      string `json:"name" db:"name"`
	SlugValue string `json:"slug" db:"slug"`
	Mode      string `json:"mode" db:"mode"`
	Active    int    `json:"active" db:"active"`

	// Audit fields
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewSite creates a new Site instance.
func NewSite(name, slug, mode string) Site {
	return Site{
		Name:      name,
		SlugValue: slug,
		Mode:      mode,
		Active:    1,
	}
}

// Type returns the entity type.
func (s *Site) Type() string {
	return "site"
}

// GetID returns the unique identifier.
func (s *Site) GetID() uuid.UUID {
	return s.ID
}

// GenID generates a new UUID for the site.
func (s *Site) GenID() {
	hm.GenID(s)
}

// SetID sets the unique identifier.
func (s *Site) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if s.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		s.ID = id
	}
}

// GetShortID returns the short ID.
func (s *Site) GetShortID() string {
	return s.ShortID
}

// GenShortID generates a short ID.
func (s *Site) GenShortID() {
	hm.GenShortID(s)
}

// SetShortID sets the short ID.
func (s *Site) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if s.ShortID == "" || shouldForce {
		s.ShortID = shortID
	}
}

func (s *Site) Slug() string {
	return s.SlugValue
}

// GenCreateValues sets creation audit values.
func (s *Site) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(s, userID...)
}

// GenUpdateValues sets update audit values.
func (s *Site) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(s, userID...)
}

// Audit method implementations
func (s *Site) GetCreatedBy() uuid.UUID { return s.CreatedBy }
func (s *Site) GetUpdatedBy() uuid.UUID { return s.UpdatedBy }
func (s *Site) GetCreatedAt() time.Time { return s.CreatedAt }
func (s *Site) GetUpdatedAt() time.Time { return s.UpdatedAt }
func (s *Site) SetCreatedAt(t time.Time) { s.CreatedAt = t }
func (s *Site) SetUpdatedAt(t time.Time) { s.UpdatedAt = t }
func (s *Site) SetCreatedBy(u uuid.UUID) { s.CreatedBy = u }
func (s *Site) SetUpdatedBy(u uuid.UUID) { s.UpdatedBy = u }

// IsZero checks if the site is uninitialized.
func (s *Site) IsZero() bool {
	return s.ID == uuid.Nil
}

// Ref returns the reference.
func (s *Site) Ref() string {
	return s.ref
}

// SetRef sets the reference.
func (s *Site) SetRef(ref string) {
	s.ref = ref
}

// OptValue for form select options.
func (s *Site) OptValue() string {
	return s.GetID().String()
}

// OptLabel for form select options.
func (s *Site) OptLabel() string {
	return s.Name
}
