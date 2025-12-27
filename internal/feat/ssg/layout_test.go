package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewlayout(t *testing.T) {
	tests := []struct {
		name        string
		layoutName  string
		description string
		code        string
		want        Layout
	}{
		{
			name:        "creates layout with all fields",
			layoutName:  "Main Layout",
			description: "Primary site layout",
			code:        "<html>...</html>",
			want: Layout{
				Name:        "Main Layout",
				Description: "Primary site layout",
				Code:        "<html>...</html>",
			},
		},
		{
			name:        "creates layout with empty fields",
			layoutName:  "",
			description: "",
			code:        "",
			want: Layout{
				Name:        "",
				Description: "",
				Code:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Newlayout(tt.layoutName, tt.description, tt.code)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Code != tt.want.Code {
				t.Errorf("Code = %v, want %v", got.Code, tt.want.Code)
			}
		})
	}
}

func TestLayoutType(t *testing.T) {
	l := Layout{}
	got := l.Type()
	want := "layout"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestLayoutGetID(t *testing.T) {
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
			l := Layout{ID: tt.id}
			got := l.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestLayoutGenID(t *testing.T) {
	l := Layout{}
	l.GenID()

	if l.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestLayoutSetID(t *testing.T) {
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
			l := Layout{ID: tt.initial}
			l.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if l.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", l.ID, want)
			}
		})
	}
}

func TestLayoutSetSiteID(t *testing.T) {
	siteID := uuid.New()
	l := Layout{}
	l.SetSiteID(siteID)

	if l.SiteID != siteID {
		t.Errorf("SetSiteID() set %v, want %v", l.SiteID, siteID)
	}
}

func TestLayoutGetShortID(t *testing.T) {
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
			l := Layout{ShortID: tt.shortID}
			got := l.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestLayoutGenShortID(t *testing.T) {
	l := Layout{}
	l.GenShortID()

	if l.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(l.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(l.ShortID))
	}
}

func TestLayoutSetShortID(t *testing.T) {
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
			l := Layout{ShortID: tt.initial}
			l.SetShortID(tt.newID, tt.force...)

			if l.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", l.ShortID, tt.expected)
			}
		})
	}
}

func TestLayoutGenCreateValues(t *testing.T) {
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
			l := Layout{}
			beforeTime := time.Now()
			l.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if l.CreatedAt.Before(beforeTime) || l.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", l.CreatedAt)
			}

			if l.UpdatedAt.Before(beforeTime) || l.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", l.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if l.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", l.CreatedBy, tt.userID[0])
				}
				if l.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", l.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestLayoutGenUpdateValues(t *testing.T) {
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
			l := Layout{}
			beforeTime := time.Now()
			l.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if l.UpdatedAt.Before(beforeTime) || l.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", l.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if l.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", l.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestLayoutGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	l := Layout{CreatedBy: userID}
	got := l.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestLayoutGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	l := Layout{UpdatedBy: userID}
	got := l.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestLayoutGetCreatedAt(t *testing.T) {
	now := time.Now()
	l := Layout{CreatedAt: now}
	got := l.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestLayoutGetUpdatedAt(t *testing.T) {
	now := time.Now()
	l := Layout{UpdatedAt: now}
	got := l.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestLayoutSetCreatedAt(t *testing.T) {
	now := time.Now()
	l := Layout{}
	l.SetCreatedAt(now)

	if l.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", l.CreatedAt, now)
	}
}

func TestLayoutSetUpdatedAt(t *testing.T) {
	now := time.Now()
	l := Layout{}
	l.SetUpdatedAt(now)

	if l.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", l.UpdatedAt, now)
	}
}

func TestLayoutSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	l := Layout{}
	l.SetCreatedBy(userID)

	if l.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", l.CreatedBy, userID)
	}
}

func TestLayoutSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	l := Layout{}
	l.SetUpdatedBy(userID)

	if l.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", l.UpdatedBy, userID)
	}
}

func TestLayoutSetHeaderImageID(t *testing.T) {
	tests := []struct {
		name    string
		imageID *uuid.UUID
	}{
		{
			name: "sets header image ID",
			imageID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
		},
		{
			name:    "sets nil header image ID",
			imageID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Layout{}
			l.SetHeaderImageID(tt.imageID)

			if tt.imageID == nil {
				if l.HeaderImageID != nil {
					t.Errorf("SetHeaderImageID() set %v, want nil", l.HeaderImageID)
				}
			} else {
				if l.HeaderImageID == nil || *l.HeaderImageID != *tt.imageID {
					t.Errorf("SetHeaderImageID() set %v, want %v", l.HeaderImageID, *tt.imageID)
				}
			}
		})
	}
}

func TestLayoutGetHeaderImageID(t *testing.T) {
	tests := []struct {
		name    string
		imageID *uuid.UUID
	}{
		{
			name: "returns set header image ID",
			imageID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
		},
		{
			name:    "returns nil when not set",
			imageID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Layout{HeaderImageID: tt.imageID}
			got := l.GetHeaderImageID()

			if tt.imageID == nil {
				if got != nil {
					t.Errorf("GetHeaderImageID() = %v, want nil", got)
				}
			} else {
				if got == nil || *got != *tt.imageID {
					t.Errorf("GetHeaderImageID() = %v, want %v", got, *tt.imageID)
				}
			}
		})
	}
}

func TestLayoutIsZero(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
		want bool
	}{
		{
			name: "returns true for uninitialized layout",
			l:    Layout{},
			want: true,
		},
		{
			name: "returns false for initialized layout",
			l:    Layout{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.l.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayoutSlug(t *testing.T) {
	tests := []struct {
		name       string
		layoutName string
		shortID    string
		want       string
	}{
		{
			name:       "generates slug from name and short ID",
			layoutName: "Main Layout",
			shortID:    "abc123",
			want:       "main-layout-abc123",
		},
		{
			name:       "handles special characters",
			layoutName: "Blog, Home!",
			shortID:    "xyz789",
			want:       "blog,-home!-xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Layout{
				Name:    tt.layoutName,
				ShortID: tt.shortID,
			}
			got := l.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayoutOptValue(t *testing.T) {
	id := uuid.New()
	l := Layout{ID: id}
	got := l.OptValue()
	want := id.String()

	if got != want {
		t.Errorf("OptValue() = %v, want %v", got, want)
	}
}

func TestLayoutOptLabel(t *testing.T) {
	l := Layout{Name: "Main Layout"}
	got := l.OptLabel()
	want := "Main Layout"

	if got != want {
		t.Errorf("OptLabel() = %v, want %v", got, want)
	}
}

func TestLayoutRef(t *testing.T) {
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
			l := Layout{RefValue: tt.ref}
			got := l.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestLayoutSetRef(t *testing.T) {
	l := Layout{}
	ref := "new-ref"
	l.SetRef(ref)

	if l.RefValue != ref {
		t.Errorf("SetRef() set %v, want %v", l.RefValue, ref)
	}
}

func TestLayoutStringID(t *testing.T) {
	id := uuid.New()
	l := Layout{ID: id}
	got := l.StringID()
	want := id.String()

	if got != want {
		t.Errorf("StringID() = %v, want %v", got, want)
	}
}
