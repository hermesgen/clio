package core

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

var key = hm.Key

type Workspace struct {
	hm.Core
}

func NewWorkspace(params hm.XParams) *Workspace {
	core := hm.NewCore("ssg-workspace", params)
	w := &Workspace{
		Core: core,
	}
	return w
}

func (w *Workspace) Setup(ctx context.Context) error {
	return w.setupDirs()
}

func (w *Workspace) setupDirs() error {
	var dirs []string
	env := w.Cfg().StrValOrDef(key.AppEnv, "prod")
	w.Log().Info("Read environment mode", "key", key.AppEnv, "value", env)

	if env == "dev" {
		w.Log().Info("Running in DEV mode, using local paths.")
		wd, err := os.Getwd()
		if err != nil {
			w.Log().Error("Cannot get working directory", "error", err)
			return err
		}
		base := filepath.Join(wd, "_workspace")
		configDir := filepath.Join(base, "config")
		dbBase := filepath.Join(base, "db")
		sitesBase := filepath.Join(base, "sites")

		dirs = []string{
			configDir,
			dbBase,
			sitesBase,
		}

		// Set single DB path (all data)
		dbDSN := "file:" + filepath.Join(dbBase, "clio.db") + "?cache=shared&mode=rwc"
		w.Cfg().Set(hm.Key.DBSQLiteDSN, dbDSN)

		// Set base paths for multi-site structure
		w.Cfg().Set(ssg.SSGKey.WorkspacePath, base)
		w.Cfg().Set(ssg.SSGKey.SitesBasePath, sitesBase)

		w.Log().Info("Overriding config for DEV mode", "db_dsn", dbDSN)

	} else {
		w.Log().Info("Running in PROD mode, using system paths.")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			w.Log().Error("Cannot get user home directory", "error", err)
			return err
		}

		dataDir := filepath.Join(homeDir, ".clio")
		configDir := filepath.Join(homeDir, ".config", "clio")
		sitesBase := filepath.Join(homeDir, "Documents", "Clio", "sites")

		dirs = []string{
			dataDir,
			configDir,
			sitesBase,
		}

		// Set single DB path (all data)
		dbDSN := "file:" + filepath.Join(dataDir, "clio.db") + "?cache=shared&mode=rwc"
		w.Cfg().Set(hm.Key.DBSQLiteDSN, dbDSN)

		// Set base paths for multi-site structure
		w.Cfg().Set(ssg.SSGKey.WorkspacePath, filepath.Join(homeDir, "Documents", "Clio"))
		w.Cfg().Set(ssg.SSGKey.SitesBasePath, sitesBase)

		w.Log().Info("Running in PROD mode", "db_dsn", dbDSN)
	}

	w.Log().Info("Ensuring base directory structure exists...")
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			w.Log().Error("Error creating directory", "path", dir, "error", err)
			return err
		}
	}
	w.Log().Info("Base directory structure verified.")

	return nil
}
