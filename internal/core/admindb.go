package core

import (
	"context"
	"embed"

	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type AdminDBManager struct {
	hm.Core
	assetsFS embed.FS
	engine   string
	migrator *hm.Migrator
	db       *sqlx.DB
}

func NewAdminDBManager(assetsFS embed.FS, engine string, migrator *hm.Migrator, params hm.XParams) *AdminDBManager {
	return &AdminDBManager{
		Core:     hm.NewCore("admin-db-manager", params),
		assetsFS: assetsFS,
		engine:   engine,
		migrator: migrator,
	}
}

func (m *AdminDBManager) Setup(ctx context.Context) error {
	dsn := m.Cfg().StrValOrDef(hm.Key.DBSQLiteDSN, "file:_workspace/db/clio.db?cache=shared&mode=rwc")

	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return err
	}
	m.db = db

	m.Log().Infof("Connected to database: %s", dsn)

	m.migrator.SetDB(m.db.DB)
	if err := m.migrator.Setup(ctx); err != nil {
		return err
	}

	return nil
}

func (m *AdminDBManager) Stop(ctx context.Context) error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *AdminDBManager) GetDB() *sqlx.DB {
	return m.db
}
