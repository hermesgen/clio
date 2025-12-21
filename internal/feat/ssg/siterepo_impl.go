package ssg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SiteRepoImpl implements SiteRepo interface.
type SiteRepoImpl struct {
	db *sqlx.DB
}

// NewSiteRepo creates a new site repository.
func NewSiteRepo(db *sqlx.DB) SiteRepo {
	return &SiteRepoImpl{db: db}
}

// CreateSite inserts a new site into sites.db.
func (r *SiteRepoImpl) CreateSite(ctx context.Context, site *Site) error {
	query := `INSERT INTO site (id, short_id, name, slug, mode, active, created_by, updated_by, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		site.ID, site.ShortID, site.Name, site.Slug(), site.Mode, site.Active,
		site.CreatedBy, site.UpdatedBy, site.CreatedAt, site.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create site: %w", err)
	}

	return nil
}

// GetSite retrieves a site by ID from sites.db.
func (r *SiteRepoImpl) GetSite(ctx context.Context, id uuid.UUID) (Site, error) {
	var site Site
	query := `SELECT * FROM site WHERE id = ? LIMIT 1`

	err := r.db.GetContext(ctx, &site, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return Site{}, fmt.Errorf("site not found: %w", err)
		}
		return Site{}, fmt.Errorf("failed to get site: %w", err)
	}

	return site, nil
}

// GetSiteBySlug retrieves a site by slug from sites.db.
func (r *SiteRepoImpl) GetSiteBySlug(ctx context.Context, slug string) (Site, error) {
	var site Site
	query := `SELECT * FROM site WHERE slug = ? LIMIT 1`

	err := r.db.GetContext(ctx, &site, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return Site{}, fmt.Errorf("site not found: %w", err)
		}
		return Site{}, fmt.Errorf("failed to get site by slug: %w", err)
	}

	return site, nil
}

// ListSites retrieves all sites from sites.db.
func (r *SiteRepoImpl) ListSites(ctx context.Context, activeOnly bool) ([]Site, error) {
	var sites []Site
	query := `SELECT * FROM site`

	if activeOnly {
		query += ` WHERE active = 1`
	}

	query += ` ORDER BY name ASC`

	err := r.db.SelectContext(ctx, &sites, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	return sites, nil
}

// UpdateSite updates an existing site in sites.db.
func (r *SiteRepoImpl) UpdateSite(ctx context.Context, site *Site) error {
	query := `UPDATE site
              SET name = ?, slug = ?, mode = ?, active = ?, updated_by = ?, updated_at = ?
              WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		site.Name, site.Slug(), site.Mode, site.Active, site.UpdatedBy, site.UpdatedAt, site.ID)

	if err != nil {
		return fmt.Errorf("failed to update site: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("site not found")
	}

	return nil
}

// DeleteSite removes a site from sites.db.
func (r *SiteRepoImpl) DeleteSite(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM site WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("site not found")
	}

	return nil
}
