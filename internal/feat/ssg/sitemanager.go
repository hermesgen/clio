package ssg

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SiteManager handles site lifecycle operations.
type SiteManager struct {
	hm.Core
	siteRepo    SiteRepo
	assetsFS    embed.FS
	engine      string
	repoFactory RepoFactory
}

// NewSiteManager creates a new site manager.
func NewSiteManager(siteRepo SiteRepo, assetsFS embed.FS, engine string, repoFactory RepoFactory, params hm.XParams) *SiteManager {
	return &SiteManager{
		Core:        hm.NewCore("site-manager", params),
		siteRepo:    siteRepo,
		assetsFS:    assetsFS,
		engine:      engine,
		repoFactory: repoFactory,
	}
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
	if mode != "normal" && mode != "blog" {
		return Site{}, fmt.Errorf("invalid mode: must be 'normal' or 'blog'")
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
		GetSiteDBPath(sitesBasePath, slug), // Will create parent db/ dir
		GetSiteMarkdownPath(sitesBasePath, slug),
		GetSiteHTMLPath(sitesBasePath, slug),
		GetSiteImagesPath(sitesBasePath, slug),
	}

	for _, dir := range dirs {
		// For the DB path, we only want the directory, not the file
		if dir == GetSiteDBPath(sitesBasePath, slug) {
			dir = GetSiteBasePath(sitesBasePath, slug) + "/db"
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
	siteBase := GetSiteBasePath(sitesBasePath, slug)

	if err := os.RemoveAll(siteBase); err != nil {
		sm.Log().Error("Failed to remove site directories", "slug", slug, "error", err)
	}
}

// initializeSiteDatabase creates and initializes the database for a site.
func (sm *SiteManager) initializeSiteDatabase(ctx context.Context, slug string, userID uuid.UUID) error {
	sitesBasePath := sm.Cfg().StrValOrDef(SSGKey.SitesBasePath, "_workspace/sites")
	dsn := GetSiteDBDSN(sitesBasePath, slug)

	// Open database connection
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to site database: %w", err)
	}
	defer db.Close()

	sm.Log().Info("Connected to site database", "dsn", dsn)

	// Create XParams for this database
	params := hm.XParams{
		Cfg: sm.Cfg(),
		Log: sm.Log(),
	}

	// Run migrations
	migrator := hm.NewMigrator(sm.assetsFS, sm.engine, params)
	migrator.SetDB(db.DB) // Set the database connection
	if err := migrator.Setup(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	sm.Log().Info("Migrations completed", "slug", slug)

	// Create repository for seeding
	queryManager := hm.NewQueryManager(sm.assetsFS, sm.engine, params)
	if err := queryManager.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup query manager: %w", err)
	}

	repo := sm.repoFactory(queryManager, params)

	// Set the database connection on the repo
	if setter, ok := repo.(interface{ SetDB(*sqlx.DB) }); ok {
		setter.SetDB(db)
	}

	if err := repo.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup repository: %w", err)
	}

	// Type assert repo to auth.Repo for seeding
	authRepo, ok := repo.(auth.Repo)
	if !ok {
		return fmt.Errorf("repository does not implement auth.Repo interface")
	}

	// Run auth seeding
	authSeeder := auth.NewSeeder(sm.assetsFS, sm.engine, authRepo, params)
	if err := authSeeder.Setup(ctx); err != nil {
		return fmt.Errorf("failed to seed auth data: %w", err)
	}

	sm.Log().Info("Auth seeding completed", "slug", slug)

	// Run SSG seeding
	ssgSeeder := NewSeeder(sm.assetsFS, sm.engine, repo, params)
	if err := ssgSeeder.Setup(ctx); err != nil {
		return fmt.Errorf("failed to seed SSG data: %w", err)
	}

	sm.Log().Info("SSG seeding completed", "slug", slug)

	return nil
}

// ListSites returns all sites (optionally only active ones).
func (sm *SiteManager) ListSites(ctx context.Context, activeOnly bool) ([]Site, error) {
	return sm.siteRepo.ListSites(ctx, activeOnly)
}

// GetSiteBySlug retrieves a site by its slug.
func (sm *SiteManager) GetSiteBySlug(ctx context.Context, slug string) (Site, error) {
	return sm.siteRepo.GetSiteBySlug(ctx, slug)
}

// GetSite retrieves a site by ID.
func (sm *SiteManager) GetSite(ctx context.Context, id uuid.UUID) (Site, error) {
	return sm.siteRepo.GetSite(ctx, id)
}
