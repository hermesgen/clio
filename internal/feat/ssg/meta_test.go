package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewMeta(t *testing.T) {
	tests := []struct {
		name      string
		contentID uuid.UUID
		want      Meta
	}{
		{
			name:      "creates meta with content ID",
			contentID: uuid.New(),
			want:      Meta{},
		},
		{
			name:      "creates meta with nil content ID",
			contentID: uuid.Nil,
			want:      Meta{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMeta(tt.contentID)

			if got.ContentID != tt.contentID {
				t.Errorf("ContentID = %v, want %v", got.ContentID, tt.contentID)
			}
		})
	}
}

func TestMetaType(t *testing.T) {
	m := Meta{}
	got := m.Type()
	want := "meta"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestMetaGetID(t *testing.T) {
	tests := []struct {
		name string
		id   uuid.UUID
	}{
		{
			name: "returns set ID",
			id:   uuid.New(),
		},
		{
			name: "returns nil UUID when not set",
			id:   uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{ID: tt.id}
			got := m.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestMetaGenID(t *testing.T) {
	m := Meta{}
	m.GenID()

	if m.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestMetaSetID(t *testing.T) {
	tests := []struct {
		name     string
		initial  uuid.UUID
		newID    uuid.UUID
		force    []bool
		expected uuid.UUID
	}{
		{
			name:     "sets ID when nil",
			initial:  uuid.Nil,
			newID:    uuid.New(),
			force:    nil,
			expected: uuid.Nil,
		},
		{
			name:     "does not override existing ID without force",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    nil,
			expected: uuid.Nil,
		},
		{
			name:     "overrides existing ID with force true",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{true},
			expected: uuid.Nil,
		},
		{
			name:     "does not override with force false",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{false},
			expected: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{ID: tt.initial}
			m.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if m.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", m.ID, want)
			}
		})
	}
}

func TestMetaGetShortID(t *testing.T) {
	tests := []struct {
		name    string
		shortID string
	}{
		{
			name:    "returns set short ID",
			shortID: "abc123",
		},
		{
			name:    "returns empty string when not set",
			shortID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{ShortID: tt.shortID}
			got := m.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestMetaGenShortID(t *testing.T) {
	m := Meta{}
	m.GenShortID()

	if m.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(m.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(m.ShortID))
	}
}

func TestMetaSetShortID(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		newID    string
		force    []bool
		expected string
	}{
		{
			name:     "sets short ID when empty",
			initial:  "",
			newID:    "xyz789",
			force:    nil,
			expected: "xyz789",
		},
		{
			name:     "does not override existing short ID without force",
			initial:  "abc123",
			newID:    "xyz789",
			force:    nil,
			expected: "abc123",
		},
		{
			name:     "overrides existing short ID with force true",
			initial:  "abc123",
			newID:    "xyz789",
			force:    []bool{true},
			expected: "xyz789",
		},
		{
			name:     "does not override with force false",
			initial:  "abc123",
			newID:    "xyz789",
			force:    []bool{false},
			expected: "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{ShortID: tt.initial}
			m.SetShortID(tt.newID, tt.force...)

			if m.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", m.ShortID, tt.expected)
			}
		})
	}
}

func TestMetaGenCreateValues(t *testing.T) {
	tests := []struct {
		name   string
		userID []uuid.UUID
	}{
		{
			name:   "sets create values with user ID",
			userID: []uuid.UUID{uuid.New()},
		},
		{
			name:   "sets create values without user ID",
			userID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{}
			beforeTime := time.Now()
			m.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if m.CreatedAt.Before(beforeTime) || m.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", m.CreatedAt)
			}

			if m.UpdatedAt.Before(beforeTime) || m.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", m.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if m.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", m.CreatedBy, tt.userID[0])
				}
				if m.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", m.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestMetaGenUpdateValues(t *testing.T) {
	tests := []struct {
		name   string
		userID []uuid.UUID
	}{
		{
			name:   "sets update values with user ID",
			userID: []uuid.UUID{uuid.New()},
		},
		{
			name:   "sets update values without user ID",
			userID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{}
			beforeTime := time.Now()
			m.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if m.UpdatedAt.Before(beforeTime) || m.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", m.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if m.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", m.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestMetaGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	m := Meta{CreatedBy: userID}
	got := m.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestMetaGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	m := Meta{UpdatedBy: userID}
	got := m.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestMetaGetCreatedAt(t *testing.T) {
	now := time.Now()
	m := Meta{CreatedAt: now}
	got := m.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestMetaGetUpdatedAt(t *testing.T) {
	now := time.Now()
	m := Meta{UpdatedAt: now}
	got := m.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestMetaSetCreatedAt(t *testing.T) {
	now := time.Now()
	m := Meta{}
	m.SetCreatedAt(now)

	if m.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", m.CreatedAt, now)
	}
}

func TestMetaSetUpdatedAt(t *testing.T) {
	now := time.Now()
	m := Meta{}
	m.SetUpdatedAt(now)

	if m.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", m.UpdatedAt, now)
	}
}

func TestMetaSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	m := Meta{}
	m.SetCreatedBy(userID)

	if m.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", m.CreatedBy, userID)
	}
}

func TestMetaSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	m := Meta{}
	m.SetUpdatedBy(userID)

	if m.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", m.UpdatedBy, userID)
	}
}

func TestMetaIsZero(t *testing.T) {
	tests := []struct {
		name string
		m    Meta
		want bool
	}{
		{
			name: "returns true for uninitialized meta",
			m:    Meta{},
			want: true,
		},
		{
			name: "returns false for initialized meta",
			m:    Meta{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetaSlug(t *testing.T) {
	tests := []struct {
		name    string
		shortID string
		want    string
	}{
		{
			name:    "returns short ID as slug",
			shortID: "abc123def456",
			want:    "abc123def456",
		},
		{
			name:    "returns empty string when short ID not set",
			shortID: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{ShortID: tt.shortID}
			got := m.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetaRef(t *testing.T) {
	tests := []struct {
		name string
		ref  string
	}{
		{
			name: "returns set ref",
			ref:  "test-ref",
		},
		{
			name: "returns empty string when not set",
			ref:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Meta{ref: tt.ref}
			got := m.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestMetaSetRef(t *testing.T) {
	m := Meta{}
	ref := "new-ref"
	m.SetRef(ref)

	if m.ref != ref {
		t.Errorf("SetRef() set %v, want %v", m.ref, ref)
	}
}
