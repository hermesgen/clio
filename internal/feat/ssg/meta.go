package ssg

import (
	"time"

	"github.com/google/uuid"

	"github.com/hermesgen/hm"
)

type Meta struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ShortID         string    `json:"-" db:"short_id"`
	ref             string    `json:"-"`
	SiteID          uuid.UUID `json:"site_id" db:"site_id"`
	ContentID       uuid.UUID `json:"content_id" db:"content_id"`
	Summary         string    `json:"summary" db:"summary"`
	Excerpt         string    `json:"excerpt" db:"excerpt"`
	Description     string    `json:"description" db:"description"`
	Keywords        string    `json:"keywords" db:"keywords"`
	Robots          string    `json:"robots" db:"robots"`
	CanonicalURL    string    `json:"canonical_url" db:"canonical_url"`
	Sitemap         string    `json:"sitemap" db:"sitemap"`
	TableOfContents bool      `json:"table_of_contents" db:"table_of_contents"`
	Share           bool      `json:"share" db:"share"`
	Comments        bool      `json:"comments" db:"comments"`
	CreatedBy       uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy       uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt       time.Time `json:"-" db:"created_at"`
	UpdatedAt       time.Time `json:"-" db:"updated_at"`
}

func NewMeta(contentID uuid.UUID) Meta {
	return Meta{
		ContentID: contentID,
	}
}

func (m *Meta) Type() string {
	return "meta"
}

func (m *Meta) GetID() uuid.UUID {
	return m.ID
}

func (m *Meta) GenID() {
	hm.GenID(m)
}

func (m *Meta) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if m.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		m.ID = id
	}
}

func (m *Meta) GetShortID() string {
	return m.ShortID
}

func (m *Meta) GenShortID() {
	hm.GenShortID(m)
}

func (m *Meta) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if m.ShortID == "" || shouldForce {
		m.ShortID = shortID
	}
}

func (m *Meta) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(m, userID...)
}

func (m *Meta) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(m, userID...)
}

func (m *Meta) GetCreatedBy() uuid.UUID {
	return m.CreatedBy
}

func (m *Meta) GetUpdatedBy() uuid.UUID {
	return m.UpdatedBy
}

func (m *Meta) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Meta) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *Meta) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

func (m *Meta) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t
}

func (m *Meta) SetCreatedBy(u uuid.UUID) {
	m.CreatedBy = u
}

func (m *Meta) SetUpdatedBy(u uuid.UUID) {
	m.UpdatedBy = u
}

func (m *Meta) IsZero() bool {
	return m.ID == uuid.Nil
}

func (m *Meta) Slug() string {
	return m.GetShortID()
}

func (m *Meta) Ref() string {
	return m.ref
}

func (m *Meta) SetRef(ref string) {
	m.ref = ref
}
