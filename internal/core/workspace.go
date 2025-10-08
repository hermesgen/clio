package core

import (
	"context"
	"os"
	"path/filepath"

	hm "github.com/hermesgen/hm"
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

		w.Cfg().Set(key.SSGWorkspacePath, base)
		w.Cfg().Set(key.SSGDocsPath, filepath.Join(base, "documents"))
		w.Cfg().Set(key.SSGMarkdownPath, filepath.Join(base, "documents", "markdown"))
		w.Cfg().Set(key.SSGHTMLPath, filepath.Join(base, "documents", "html"))
		w.Cfg().Set(key.SSGAssetsPath, filepath.Join(base, "documents", "assets"))
		w.Cfg().Set(key.SSGImagesPath, filepath.Join(base, "documents", "assets", "images"))

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

		w.Cfg().Set(key.SSGWorkspacePath, basePath)
		w.Cfg().Set(key.SSGDocsPath, docsPath)
		w.Cfg().Set(key.SSGMarkdownPath, markdownPath)
		w.Cfg().Set(key.SSGHTMLPath, htmlPath)
		w.Cfg().Set(key.SSGAssetsPath, assetsPath)
		w.Cfg().Set(key.SSGImagesPath, imagesPath)
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
