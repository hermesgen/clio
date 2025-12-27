package ssg

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SiteManager handles site lifecycle operations.
type SiteManager struct {
	hm.Core
	siteRepo SiteRepo
	repo     Repo
	assetsFS embed.FS
	engine   string
}

// NewSiteManager creates a new site manager.
func NewSiteManager(repo Repo, assetsFS embed.FS, engine string, params hm.XParams) *SiteManager {
	return &SiteManager{
		Core:     hm.NewCore("site-manager", params),
		repo:     repo,
		assetsFS: assetsFS,
		engine:   engine,
	}
}

func (sm *SiteManager) Setup(ctx context.Context) error {
	// Extract DB from repo for SiteRepo
	type DBGetter interface {
		GetDB() *sqlx.DB
	}
	if dbGetter, ok := sm.repo.(DBGetter); ok {
		sm.siteRepo = NewSiteRepo(dbGetter.GetDB())
	} else {
		return fmt.Errorf("repo does not provide GetDB() method")
	}
	return nil
}

// CreateSite creates a new site with full initialization.
// Steps:
// - Create site record in sites.db
// - Create directory structure
// - Create and initialize site database
// - Run migrations
// - Run seeding
func (sm *SiteManager) CreateSite(ctx context.Context, name, slug, mode string, userID uuid.UUID) (Site, error) {
	sm.Log().Info("Creating new site", "slug", slug, "mode", mode)

	// Validate mode
	if mode != "structured" && mode != "blog" {
		return Site{}, fmt.Errorf("invalid mode: must be 'structured' or 'blog'")
	}

	// Check if slug is available
	existing, err := sm.siteRepo.GetSiteBySlug(ctx, slug)
	if err == nil && !existing.IsZero() {
		return Site{}, fmt.Errorf("site with slug '%s' already exists", slug)
	}

	// Create site record in sites.db
	site := NewSite(name, slug, mode)
	site.GenID()
	site.GenShortID()
	site.GenCreateValues(userID)

	if err := sm.siteRepo.CreateSite(ctx, &site); err != nil {
		return Site{}, fmt.Errorf("failed to create site record: %w", err)
	}

	sm.Log().Info("Site record created", "id", site.ID, "slug", slug)

	// Create directory structure
	if err := sm.createSiteDirectories(slug); err != nil {
		// Rollback: delete site record
		sm.siteRepo.DeleteSite(ctx, site.ID)
		return Site{}, fmt.Errorf("failed to create directories: %w", err)
	}

	sm.Log().Info("Site directories created", "slug", slug)

	// Initialize site database
	if err := sm.initializeSiteDatabase(ctx, slug, userID); err != nil {
		// Rollback: delete directories and site record
		sm.deleteSiteDirectories(slug)
		sm.siteRepo.DeleteSite(ctx, site.ID)
		return Site{}, fmt.Errorf("failed to initialize database: %w", err)
	}

	sm.Log().Info("Site database initialized", "slug", slug)
	sm.Log().Info("Site created successfully", "id", site.ID, "slug", slug)

	return site, nil
}

// createSiteDirectories creates the directory structure for a site.
func (sm *SiteManager) createSiteDirectories(slug string) error {
	sitesBasePath := sm.Cfg().StrValOrDef(SSGKey.SitesBasePath, "_workspace/sites")

	dirs := []string{
		GetSiteDBPath(sitesBasePath, slug), // e.g., _workspace/sites/slug/db/clio.db
		GetSiteMarkdownPath(sitesBasePath, slug),
		GetSiteHTMLPath(sitesBasePath, slug),
		GetSiteImagesPath(sitesBasePath, slug),
	}

	for _, dir := range dirs {
		// For the DB path, we only want the directory (parent of clio.db), not the file
		if dir == GetSiteDBPath(sitesBasePath, slug) {
			dir = filepath.Join(sitesBasePath, slug, "db") // Just the directory, not the .db file
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// deleteSiteDirectories removes the directory structure for a site.
func (sm *SiteManager) deleteSiteDirectories(slug string) {
	sitesBasePath := sm.Cfg().StrValOrDef(SSGKey.SitesBasePath, "_workspace/sites")

	// Remove entire site directory (includes db, documents, etc.)
	siteBase := GetSiteBasePath(sitesBasePath, slug)
	if err := os.RemoveAll(siteBase); err != nil {
		sm.Log().Error("Failed to remove site directories", "slug", slug, "error", err)
	}
}

// initializeSiteDatabase initializes data for a new site.
// With single-DB architecture, migrations are already run.
// This only performs seeding for the new site.
func (sm *SiteManager) initializeSiteDatabase(ctx context.Context, slug string, userID uuid.UUID) error {
	sm.Log().Info("Initializing site data", "slug", slug)

	params := hm.XParams{
		Cfg: sm.Cfg(),
		Log: sm.Log(),
	}

	// Type assert repo to auth.Repo for seeding
	authRepo, ok := sm.repo.(auth.Repo)
	if !ok {
		return fmt.Errorf("repository does not implement auth.Repo interface")
	}

	// Run auth seeding for this site
	authSeeder := auth.NewSeeder(sm.assetsFS, sm.engine, authRepo, params)
	if err := authSeeder.Setup(ctx); err != nil {
		return fmt.Errorf("failed to seed auth data: %w", err)
	}

	sm.Log().Info("Auth seeding completed", "slug", slug)

	// Run SSG seeding for this site
	ssgSeeder := NewSeeder(sm.assetsFS, sm.engine, sm.repo, params)
	if err := ssgSeeder.Setup(ctx); err != nil {
		return fmt.Errorf("failed to seed SSG data: %w", err)
	}

	sm.Log().Info("SSG seeding completed", "slug", slug)

	return nil
}

// ListSites returns all sites (optionally only active ones).
// Auto-deletes sites whose database files don't exist.
func (sm *SiteManager) ListSites(ctx context.Context, activeOnly bool) ([]Site, error) {
	sites, err := sm.siteRepo.ListSites(ctx, activeOnly)
	if err != nil {
		return nil, err
	}

	sitesBasePath := sm.Cfg().StrValOrDef(SSGKey.SitesBasePath, "_workspace/sites")
	validSites := make([]Site, 0, len(sites))

	for _, site := range sites {
		// Check if site directory exists (not just DB file, since DB is created on first access)
		siteDir := GetSiteBasePath(sitesBasePath, site.Slug())
		if _, err := os.Stat(siteDir); err == nil {
			validSites = append(validSites, site)
		} else {
			sm.Log().Info("Site directory not found, auto-deleting orphaned record", "slug", site.Slug(), "path", siteDir)
			if delErr := sm.siteRepo.DeleteSite(ctx, site.ID); delErr != nil {
				sm.Log().Error("Failed to auto-delete orphaned site", "slug", site.Slug(), "error", delErr)
			}
		}
	}

	return validSites, nil
}

// GetSiteBySlug retrieves a site by its slug.
func (sm *SiteManager) GetSiteBySlug(ctx context.Context, slug string) (Site, error) {
	return sm.siteRepo.GetSiteBySlug(ctx, slug)
}

// GetSite retrieves a site by ID.
func (sm *SiteManager) GetSite(ctx context.Context, id uuid.UUID) (Site, error) {
	return sm.siteRepo.GetSite(ctx, id)
}

// DeleteSite removes a site from the database (files remain for backup).
func (sm *SiteManager) DeleteSite(ctx context.Context, id uuid.UUID) (string, error) {
	site, err := sm.siteRepo.GetSite(ctx, id)
	if err != nil {
		return "", fmt.Errorf("site not found: %w", err)
	}

	if err := sm.siteRepo.DeleteSite(ctx, id); err != nil {
		return "", fmt.Errorf("failed to delete site: %w", err)
	}

	sm.Log().Info("Site deleted from database, files preserved", "slug", site.Slug(), "id", id)

	return site.Slug(), nil
}
