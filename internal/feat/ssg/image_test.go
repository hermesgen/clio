package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewImage(t *testing.T) {
	img := NewImage()

	if img.ID == uuid.Nil {
		t.Error("NewImage() did not generate UUID")
	}
}

func TestImageType(t *testing.T) {
	img := Image{}
	got := img.Type()
	want := "image"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestImageGetID(t *testing.T) {
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
			img := Image{ID: tt.id}
			got := img.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestImageGenID(t *testing.T) {
	img := Image{}
	img.GenID()

	if img.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestImageSetID(t *testing.T) {
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
			expected: uuid.Nil, // Will be set to newID
		},
		{
			name:     "does not override existing ID without force",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    nil,
			expected: uuid.Nil, // Will keep initial
		},
		{
			name:     "overrides existing ID with force true",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{true},
			expected: uuid.Nil, // Will be set to newID
		},
		{
			name:     "does not override with force false",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{false},
			expected: uuid.Nil, // Will keep initial
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := Image{ID: tt.initial}
			img.SetID(tt.newID, tt.force...)

			// Determine expected based on logic
			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if img.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", img.ID, want)
			}
		})
	}
}

func TestImageGetShortID(t *testing.T) {
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
			img := Image{ShortID: tt.shortID}
			got := img.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestImageGenShortID(t *testing.T) {
	img := Image{}
	img.GenShortID()

	if img.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(img.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(img.ShortID))
	}
}

func TestImageSetShortID(t *testing.T) {
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
			img := Image{ShortID: tt.initial}
			img.SetShortID(tt.newID, tt.force...)

			if img.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", img.ShortID, tt.expected)
			}
		})
	}
}

func TestImageGenCreateValues(t *testing.T) {
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
			img := Image{}
			beforeTime := time.Now()
			img.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if img.CreatedAt.Before(beforeTime) || img.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", img.CreatedAt)
			}

			if img.UpdatedAt.Before(beforeTime) || img.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", img.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if img.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", img.CreatedBy, tt.userID[0])
				}
				if img.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", img.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestImageGenUpdateValues(t *testing.T) {
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
			img := Image{}
			beforeTime := time.Now()
			img.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if img.UpdatedAt.Before(beforeTime) || img.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", img.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if img.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", img.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestImageGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	img := Image{CreatedBy: userID}
	got := img.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestImageGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	img := Image{UpdatedBy: userID}
	got := img.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestImageGetCreatedAt(t *testing.T) {
	now := time.Now()
	img := Image{CreatedAt: now}
	got := img.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestImageGetUpdatedAt(t *testing.T) {
	now := time.Now()
	img := Image{UpdatedAt: now}
	got := img.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestImageSetCreatedAt(t *testing.T) {
	now := time.Now()
	img := Image{}
	img.SetCreatedAt(now)

	if img.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", img.CreatedAt, now)
	}
}

func TestImageSetUpdatedAt(t *testing.T) {
	now := time.Now()
	img := Image{}
	img.SetUpdatedAt(now)

	if img.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", img.UpdatedAt, now)
	}
}

func TestImageSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	img := Image{}
	img.SetCreatedBy(userID)

	if img.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", img.CreatedBy, userID)
	}
}

func TestImageSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	img := Image{}
	img.SetUpdatedBy(userID)

	if img.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", img.UpdatedBy, userID)
	}
}

func TestImageIsZero(t *testing.T) {
	tests := []struct {
		name string
		img  Image
		want bool
	}{
		{
			name: "returns true for uninitialized image",
			img:  Image{},
			want: true,
		},
		{
			name: "returns false for initialized image",
			img:  Image{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.img.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		fileName string
		shortID  string
		want     string
	}{
		{
			name:     "uses title when available",
			title:    "Beautiful Sunset",
			fileName: "img001.jpg",
			shortID:  "abc123",
			want:     "beautiful-sunset-abc123",
		},
		{
			name:     "uses filename when title is empty",
			title:    "",
			fileName: "vacation-photo.jpg",
			shortID:  "xyz789",
			want:     "vacation-photo.jpg-xyz789",
		},
		{
			name:     "handles special characters in title",
			title:    "Hello, World!",
			fileName: "test.png",
			shortID:  "def456",
			want:     "hello,-world!-def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := Image{
				Title:    tt.title,
				FileName: tt.fileName,
				ShortID:  tt.shortID,
			}
			got := img.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageRef(t *testing.T) {
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
			img := Image{ref: tt.ref}
			got := img.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestImageSetRef(t *testing.T) {
	img := Image{}
	ref := "new-ref"
	img.SetRef(ref)

	if img.ref != ref {
		t.Errorf("SetRef() set %v, want %v", img.ref, ref)
	}
}
