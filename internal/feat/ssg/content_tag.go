package ssg

import (
	"time"

	"github.com/google/uuid"

	"github.com/hermesgen/hm"
)

// ContentTag model represents the many-to-many relationship between Content and Tag.
type ContentTag struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	ShortID string    `json:"-" db:"short_id"`
	ref     string    `json:"-"`

	ContentID uuid.UUID `json:"content_id" db:"content_id"`
	TagID     uuid.UUID `json:"tag_id" db:"tag_id"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewContentTag creates a new ContentTag.
func NewContentTag(contentID, tagID uuid.UUID) ContentTag {
	ct := ContentTag{
		ContentID: contentID,
		TagID:     tagID,
	}

	return ct
}

// Type returns the type of the entity.
func (ct *ContentTag) Type() string {
	return "content-tag"
}


// GetID returns the unique identifier of the entity.
func (ct *ContentTag) GetID() uuid.UUID {
	return ct.ID
}

// GenID delegates to the functional helper.
func (ct *ContentTag) GenID() {
	hm.GenID(ct)
}

// SetID sets the unique identifier of the entity.
func (ct *ContentTag) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if ct.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		ct.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (ct *ContentTag) GetShortID() string {
	return ct.ShortID
}

// GenShortID delegates to the functional helper.
func (ct *ContentTag) GenShortID() {
	hm.GenShortID(ct)
}

// SetShortID sets the short ID of the entity.
func (ct *ContentTag) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if ct.ShortID == "" || shouldForce {
		ct.ShortID = shortID
	}
}

// GenCreateValues delegates to the functional helper.
func (ct *ContentTag) GenCreateValues(userID ...uuid.UUID) {
	hm.SetCreateValues(ct, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (ct *ContentTag) GenUpdateValues(userID ...uuid.UUID) {
	hm.SetUpdateValues(ct, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (ct *ContentTag) GetCreatedBy() uuid.UUID {
	return ct.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (ct *ContentTag) GetUpdatedBy() uuid.UUID {
	return ct.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (ct *ContentTag) GetCreatedAt() time.Time {
	return ct.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (ct *ContentTag) GetUpdatedAt() time.Time {
	return ct.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (ct *ContentTag) SetCreatedAt(createdAt time.Time) {
	ct.CreatedAt = createdAt
}

// SetUpdatedAt implements the Auditable interface.
func (ct *ContentTag) SetUpdatedAt(updatedAt time.Time) {
	ct.UpdatedAt = updatedAt
}

// SetCreatedBy implements the Auditable interface.
func (ct *ContentTag) SetCreatedBy(createdBy uuid.UUID) {
	ct.CreatedBy = createdBy
}

// SetUpdatedBy implements the Auditable interface.
func (ct *ContentTag) SetUpdatedBy(updatedBy uuid.UUID) {
	ct.UpdatedBy = updatedBy
}

// IsZero returns true if the ContentTag is uninitialized.
func (ct *ContentTag) IsZero() bool {
	return ct.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (ct *ContentTag) Slug() string {
	return hm.Normalize(ct.Type()) + "-" + ct.GetShortID()
}

func (ct *ContentTag) OptValue() string {
	return ct.GetID().String()
}

func (ct *ContentTag) OptLabel() string {
	return ct.GetShortID()
}

func (ct *ContentTag) Ref() string {
	return ct.ref
}

func (ct *ContentTag) SetRef(ref string) {
	ct.ref = ref
}
