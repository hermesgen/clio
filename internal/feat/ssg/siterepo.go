package ssg

import (
	"context"

	"github.com/google/uuid"
)

// SiteRepo manages site metadata (operates on sites.db).
// This is separate from the main Repo which operates on per-site databases.
type SiteRepo interface {
	CreateSite(ctx context.Context, site *Site) error
	GetSite(ctx context.Context, id uuid.UUID) (Site, error)
	GetSiteBySlug(ctx context.Context, slug string) (Site, error)
	ListSites(ctx context.Context, activeOnly bool) ([]Site, error)
	UpdateSite(ctx context.Context, site *Site) error
	DeleteSite(ctx context.Context, id uuid.UUID) error
}
