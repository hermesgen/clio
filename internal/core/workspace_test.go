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
		t.Fatalf("os.UserHomeDir: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "clio-test-*")
	if err != nil {
		t.Fatalf("os.MkdirTemp: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("os.Chdir: %v", err)
	}
	defer os.Chdir(originalWd)

	tests := []struct {
		name          string
		env           string
		expectedPaths map[string]string
		expectedDirs  []string
	}{
		{
			name: "dev mode",
			env:  "dev",
			expectedPaths: map[string]string{
				hm.Key.DBSQLiteDSN:       "file:" + filepath.Join(tempDir, "_workspace", "db", "clio.db") + "?cache=shared&mode=rwc",
				ssg.SSGKey.WorkspacePath: filepath.Join(tempDir, "_workspace"),
				ssg.SSGKey.SitesBasePath: filepath.Join(tempDir, "_workspace", "sites"),
			},
			expectedDirs: []string{
				filepath.Join(tempDir, "_workspace", "config"),
				filepath.Join(tempDir, "_workspace", "db"),
				filepath.Join(tempDir, "_workspace", "sites"),
			},
		},
		{
			name: "prod mode",
			env:  "prod",
			expectedPaths: map[string]string{
				hm.Key.DBSQLiteDSN:       "file:" + filepath.Join(homeDir, ".clio", "clio.db") + "?cache=shared&mode=rwc",
				ssg.SSGKey.WorkspacePath: filepath.Join(homeDir, "Documents", "Clio"),
				ssg.SSGKey.SitesBasePath: filepath.Join(homeDir, "Documents", "Clio", "sites"),
			},
			expectedDirs: []string{
				filepath.Join(homeDir, ".clio"),
				filepath.Join(homeDir, ".config", "clio"),
				filepath.Join(homeDir, "Documents", "Clio", "sites"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := hm.NewConfig()
			cfg.Set(hm.Key.AppEnv, tt.env)

			logger := hm.NewLogger("")
			params := hm.XParams{Cfg: cfg, Log: logger}
			ws := core.NewWorkspace(params)

			if err := ws.Setup(context.Background()); err != nil {
				t.Fatalf("Setup: %v", err)
			}

			for key, want := range tt.expectedPaths {
				got := cfg.StrValOrDef(key, "")
				if got != want {
					t.Errorf("cfg.Get(%q) = %q, want %q", key, got, want)
				}
			}

			if tt.env == "dev" {
				for _, dir := range tt.expectedDirs {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Errorf("directory not created: %q", dir)
					}
				}
			}
		})
	}
}
