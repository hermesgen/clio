package ssg

import (
	"context"
	"embed"
	"fmt"

	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// RepoFactory is a function that creates a Repo from a QueryManager.
type RepoFactory func(*hm.QueryManager, hm.XParams) Repo

// RepoManager creates repository instances for specific sites.
type RepoManager struct {
	hm.Core
	assetsFS    embed.FS
	engine      string
	repoFactory RepoFactory
}

// NewRepoManager creates a new repository manager.
func NewRepoManager(assetsFS embed.FS, engine string, repoFactory RepoFactory, params hm.XParams) *RepoManager {
	return &RepoManager{
		Core:        hm.NewCore("repo-manager", params),
		assetsFS:    assetsFS,
		engine:      engine,
		repoFactory: repoFactory,
	}
}

// GetRepoForSite creates a repository instance connected to a specific site's database.
func (rm *RepoManager) GetRepoForSite(ctx context.Context, siteSlug string) (Repo, error) {
	dbBasePath := rm.Cfg().StrValOrDef(SSGKey.DBBasePath, "_workspace/db")
	dsn := GetSiteDBDSN(dbBasePath, siteSlug)

	// Open database connection
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to site database: %w", err)
	}

	rm.Log().Debug("Connected to site database", "slug", siteSlug, "dsn", dsn)

	// Create query manager
	params := hm.XParams{
		Cfg: rm.Cfg(),
		Log: rm.Log(),
	}

	queryManager := hm.NewQueryManager(rm.assetsFS, rm.engine, params)
	if err := queryManager.Setup(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to setup query manager: %w", err)
	}

	// Create repository using factory
	repo := rm.repoFactory(queryManager, params)

	// Set the database connection on the repo
	if setter, ok := repo.(interface{ SetDB(*sqlx.DB) }); ok {
		setter.SetDB(db)
	}

	return repo, nil
}

// CloseRepo closes a repository's database connection.
func (rm *RepoManager) CloseRepo(repo Repo) error {
	// The repo doesn't expose Close directly, we'd need to modify it
	// For now, connections will be closed by GC
	// TODO: Add Close() to Repo interface if needed
	return nil
}
