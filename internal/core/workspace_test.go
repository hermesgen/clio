package core_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hermesgen/clio/internal/core"
	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func TestWorkspaceSetup(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get user home directory: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "clio-test-*")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("could not change to temp dir: %v", err)
	}
	defer os.Chdir(originalWd)

	testCases := []struct {
		name          string
		env           string
		expectedPaths map[string]string
	}{
		{
			name: "dev mode",
			env:  "dev",
			expectedPaths: map[string]string{
				hm.Key.DBSQLiteDSN:       "file:" + filepath.Join(tempDir, "_workspace", "db", "clio.db") + "?cache=shared&mode=rwc",
				ssg.SSGKey.WorkspacePath: filepath.Join(tempDir, "_workspace"),
				ssg.SSGKey.DocsPath:      filepath.Join(tempDir, "_workspace", "documents"),
				ssg.SSGKey.MarkdownPath:  filepath.Join(tempDir, "_workspace", "documents", "markdown"),
				ssg.SSGKey.HTMLPath:      filepath.Join(tempDir, "_workspace", "documents", "html"),
				ssg.SSGKey.AssetsPath:    filepath.Join(tempDir, "_workspace", "documents", "assets"),
				ssg.SSGKey.ImagesPath:    filepath.Join(tempDir, "_workspace", "documents", "assets", "images"),
			},
		},
		{
			name: "prod mode",
			env:  "prod",
			expectedPaths: map[string]string{
				ssg.SSGKey.WorkspacePath: filepath.Join(homeDir, ".clio"),
				ssg.SSGKey.DocsPath:      filepath.Join(homeDir, "Documents", "Clio"),
				ssg.SSGKey.MarkdownPath:  filepath.Join(homeDir, "Documents", "Clio", "markdown"),
				ssg.SSGKey.HTMLPath:      filepath.Join(homeDir, "Documents", "Clio", "html"),
				ssg.SSGKey.AssetsPath:    filepath.Join(homeDir, "Documents", "Clio", "assets"),
				ssg.SSGKey.ImagesPath:    filepath.Join(homeDir, "Documents", "Clio", "assets", "images"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := hm.NewConfig()
			cfg.Set(hm.Key.AppEnv, tc.env)

			logger := hm.NewLogger("")
			ws := core.NewWorkspace(hm.WithCfg(cfg), hm.WithLog(logger))

			if err := ws.Setup(context.Background()); err != nil {
				t.Fatalf("ws.Setup() failed: %v", err)
			}

			for key, expectedPath := range tc.expectedPaths {
				actualPath := cfg.StrValOrDef(key, "")
				if actualPath != expectedPath {
					t.Errorf("config value for key %q: got %q, want %q", key, actualPath, expectedPath)
				}
			}

			if tc.name == "dev mode" {
				for key, path := range tc.expectedPaths {
					if key == hm.Key.DBSQLiteDSN {
						continue
					}
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("directory %q should have been created in dev mode", path)
					}
				}
			}
		})
	}
}
