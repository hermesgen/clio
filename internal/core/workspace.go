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
		dbDir := filepath.Join(base, "db")
		dirs = []string{
			filepath.Join(base, "config"),
			dbDir,
			filepath.Join(base, "documents", "markdown"),
			filepath.Join(base, "documents", "html"),
			filepath.Join(base, "documents", "assets", "images"),
		}

		// Override config values for dev mode
		devDSN := "file:" + filepath.Join(dbDir, "clio.db") + "?cache=shared&mode=rwc"
		w.Cfg().Set(key.DBSQLiteDSN, devDSN)

		w.Cfg().Set(ssg.SSGKey.WorkspacePath, base)
		w.Cfg().Set(ssg.SSGKey.DocsPath, filepath.Join(base, "documents"))
		w.Cfg().Set(ssg.SSGKey.MarkdownPath, filepath.Join(base, "documents", "markdown"))
		w.Cfg().Set(ssg.SSGKey.HTMLPath, filepath.Join(base, "documents", "html"))
		w.Cfg().Set(ssg.SSGKey.AssetsPath, filepath.Join(base, "documents", "assets"))
		w.Cfg().Set(ssg.SSGKey.ImagesPath, filepath.Join(base, "documents", "assets", "images"))

		w.Log().Info("Overriding config for DEV mode", "key", key.DBSQLiteDSN, "value", devDSN)

	} else {
		w.Log().Info("Running in PROD mode, using system paths.")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			w.Log().Error("Cannot get user home directory", "error", err)
			return err
		}

		basePath := filepath.Join(homeDir, ".clio")
		docsPath := filepath.Join(homeDir, "Documents", "Clio")
		markdownPath := filepath.Join(docsPath, "markdown")
		htmlPath := filepath.Join(docsPath, "html")
		assetsPath := filepath.Join(docsPath, "assets")
		imagesPath := filepath.Join(assetsPath, "images")

		dirs = []string{
			filepath.Join(homeDir, ".config", "clio"),
			basePath,
			markdownPath,
			htmlPath,
			imagesPath,
		}

		w.Cfg().Set(ssg.SSGKey.WorkspacePath, basePath)
		w.Cfg().Set(ssg.SSGKey.DocsPath, docsPath)
		w.Cfg().Set(ssg.SSGKey.MarkdownPath, markdownPath)
		w.Cfg().Set(ssg.SSGKey.HTMLPath, htmlPath)
		w.Cfg().Set(ssg.SSGKey.AssetsPath, assetsPath)
		w.Cfg().Set(ssg.SSGKey.ImagesPath, imagesPath)
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
