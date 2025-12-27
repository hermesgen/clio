package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		name    string
		tagName string
		want    Tag
	}{
		{
			name:    "creates tag with name",
			tagName: "golang",
			want: Tag{
				Name: "golang",
			},
		},
		{
			name:    "creates tag with empty name",
			tagName: "",
			want: Tag{
				Name: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTag(tt.tagName)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}

func TestTagType(t *testing.T) {
	tag := Tag{}
	got := tag.Type()
	want := "tag"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestTagGetID(t *testing.T) {
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
			tag := Tag{ID: tt.id}
			got := tag.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestTagGenID(t *testing.T) {
	tag := Tag{}
	tag.GenID()

	if tag.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestTagSetID(t *testing.T) {
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
			tag := Tag{ID: tt.initial}
			tag.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if tag.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", tag.ID, want)
			}
		})
	}
}

func TestTagGetShortID(t *testing.T) {
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
			tag := Tag{ShortID: tt.shortID}
			got := tag.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestTagGenShortID(t *testing.T) {
	tag := Tag{}
	tag.GenShortID()

	if tag.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(tag.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(tag.ShortID))
	}
}

func TestTagSetShortID(t *testing.T) {
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
			tag := Tag{ShortID: tt.initial}
			tag.SetShortID(tt.newID, tt.force...)

			if tag.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", tag.ShortID, tt.expected)
			}
		})
	}
}

func TestTagGenCreateValues(t *testing.T) {
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
			tag := Tag{}
			beforeTime := time.Now()
			tag.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if tag.CreatedAt.Before(beforeTime) || tag.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", tag.CreatedAt)
			}

			if tag.UpdatedAt.Before(beforeTime) || tag.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", tag.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if tag.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", tag.CreatedBy, tt.userID[0])
				}
				if tag.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", tag.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestTagGenUpdateValues(t *testing.T) {
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
			tag := Tag{}
			beforeTime := time.Now()
			tag.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if tag.UpdatedAt.Before(beforeTime) || tag.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", tag.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if tag.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", tag.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestTagGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	tag := Tag{CreatedBy: userID}
	got := tag.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestTagGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	tag := Tag{UpdatedBy: userID}
	got := tag.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestTagGetCreatedAt(t *testing.T) {
	now := time.Now()
	tag := Tag{CreatedAt: now}
	got := tag.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestTagGetUpdatedAt(t *testing.T) {
	now := time.Now()
	tag := Tag{UpdatedAt: now}
	got := tag.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestTagSetCreatedAt(t *testing.T) {
	now := time.Now()
	tag := Tag{}
	tag.SetCreatedAt(now)

	if tag.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", tag.CreatedAt, now)
	}
}

func TestTagSetUpdatedAt(t *testing.T) {
	now := time.Now()
	tag := Tag{}
	tag.SetUpdatedAt(now)

	if tag.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", tag.UpdatedAt, now)
	}
}

func TestTagSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	tag := Tag{}
	tag.SetCreatedBy(userID)

	if tag.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", tag.CreatedBy, userID)
	}
}

func TestTagSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	tag := Tag{}
	tag.SetUpdatedBy(userID)

	if tag.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", tag.UpdatedBy, userID)
	}
}

func TestTagIsZero(t *testing.T) {
	tests := []struct {
		name string
		tag  Tag
		want bool
	}{
		{
			name: "returns true for uninitialized tag",
			tag:  Tag{},
			want: true,
		},
		{
			name: "returns false for initialized tag",
			tag:  Tag{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tag.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagSlug(t *testing.T) {
	tests := []struct {
		name      string
		tagName   string
		slugField string
		shortID   string
		want      string
	}{
		{
			name:      "uses slug field when available",
			tagName:   "Golang",
			slugField: "go-lang",
			shortID:   "abc123",
			want:      "go-lang",
		},
		{
			name:      "generates slug from name when slug field empty",
			tagName:   "Rust Language",
			slugField: "",
			shortID:   "xyz789",
			want:      "rust-language-xyz789",
		},
		{
			name:      "handles special characters in name",
			tagName:   "C++",
			slugField: "",
			shortID:   "def456",
			want:      "c++-def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := Tag{
				Name:      tt.tagName,
				SlugField: tt.slugField,
				ShortID:   tt.shortID,
			}
			got := tag.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagOptValue(t *testing.T) {
	id := uuid.New()
	tag := Tag{ID: id}
	got := tag.OptValue()
	want := id.String()

	if got != want {
		t.Errorf("OptValue() = %v, want %v", got, want)
	}
}

func TestTagOptLabel(t *testing.T) {
	tag := Tag{Name: "golang"}
	got := tag.OptLabel()
	want := "golang"

	if got != want {
		t.Errorf("OptLabel() = %v, want %v", got, want)
	}
}

func TestTagRef(t *testing.T) {
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
			tag := Tag{ref: tt.ref}
			got := tag.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestTagSetRef(t *testing.T) {
	tag := Tag{}
	ref := "new-ref"
	tag.SetRef(ref)

	if tag.ref != ref {
		t.Errorf("SetRef() set %v, want %v", tag.ref, ref)
	}
}
