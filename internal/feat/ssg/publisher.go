package ssg

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	hm "github.com/hermesgen/hm"
)

// PublisherConfig holds all configuration needed for a publishing operation.
type PublisherConfig struct {
	RepoURL      string // Full URL to the GitHub repository
	Branch       string // Target branch for publishing (e.g., "gh-pages")
	PagesSubdir  string // Subdirectory within the repo (e.g., "" for root, "docs")
	Auth         hm.GitAuth
	CommitAuthor hm.GitCommit
}

// Publisher defines the interface for orchestrating the publishing process.
type Publisher interface {
	// Validate checks if the provided configuration is valid for publishing.
	Validate(cfg PublisherConfig) error

	// Publish takes the source directory (containing generated HTML) and publishes it
	// to the configured GitHub repo.
	Publish(ctx context.Context, cfg PublisherConfig, sourceDir string) (commitURL string, err error)

	// Plan performs a dry-run, showing what changes would be made without
	// actually pushing to the remote.
	Plan(ctx context.Context, cfg PublisherConfig, sourceDir string) (PlanReport, error)
}

type PlanReport struct {
	Added    []string
	Modified []string
	Removed  []string
	Summary  string
}

type publisher struct {
	hm.Core
	gitClient hm.GitClient
}

// NewPublisherWithParams creates a Publisher with XParams.
func NewPublisher(gitClient hm.GitClient, params hm.XParams) *publisher {
	return &publisher{
		Core:      hm.NewCore("ssg-pub", params),
		gitClient: gitClient,
	}
}

// Validate implementation
func (p *publisher) Validate(cfg PublisherConfig) error {
	// NOTE: See if we can use hm.Validator approach here
	if cfg.RepoURL == "" {
		return fmt.Errorf("repo URL cannot be empty")
	}

	if cfg.Branch == "" {
		return fmt.Errorf("publish branch cannot be empty")
	}

	return nil
}

func (p *publisher) Publish(ctx context.Context, cfg PublisherConfig, sourceDir string) (commitURL string, err error) {
	p.Log().Info("Starting publish process")

	// Temp dir for the publisher's work
	parentTempDir, err := os.MkdirTemp("", "clio-publish-parent-*")
	if err != nil {
		return "", fmt.Errorf("cannot create parent temp dir: %w", err)
	}
	defer os.RemoveAll(parentTempDir)

	// The actual repository will be cloned into a subdirectory
	tempDir := filepath.Join(parentTempDir, "repo")

	// NOTE: This is a temporary hack and we need to get rid of it, but for now
	// it does the trick.
	// Create a temporary script to provide the GitHub token for the publisher's git operations
	askpassScriptPath := filepath.Join(parentTempDir, "git-askpass.sh")
	err = os.WriteFile(askpassScriptPath, []byte(fmt.Sprintf("#!/bin/sh\necho %s", cfg.Auth.Token)), 0700)
	if err != nil {
		return "", fmt.Errorf("cannot create askpass script: %w", err)
	}

	// Set GIT_ASKPASS environment variable for all git commands in tempDir
	env := os.Environ()
	env = append(env, "GIT_ASKPASS="+askpassScriptPath)

	if err := p.gitClient.Clone(ctx, cfg.RepoURL, tempDir, cfg.Auth, env); err != nil {
		return "", fmt.Errorf("cannot clone repo: %w", err)
	}
	p.Log().Info("Repo cloned")

	// Checkout target branch
	if err := p.gitClient.Checkout(ctx, tempDir, cfg.Branch, false, env); err != nil {
		return "", fmt.Errorf("cannot checkout branch: %w", err)
	}
	p.Log().Info("Checked out branch", "branch", cfg.Branch)

	// Clean and copy source dir content
	targetDir := filepath.Join(tempDir, cfg.PagesSubdir)
	p.Log().Info("Cleaning target directory", "path", targetDir)

	if cfg.PagesSubdir == "" {
		// Remove all contents except .git
		dirs, err := os.ReadDir(tempDir)
		if err != nil {
			return "", fmt.Errorf("cannot read temp dir: %w", err)
		}

		for _, d := range dirs {
			if d.Name() == ".git" {
				continue
			}
			if err := os.RemoveAll(filepath.Join(tempDir, d.Name())); err != nil {
				return "", fmt.Errorf("cannot remove %s from temp dir: %w", d.Name(), err)
			}
		}

	} else {
		// When publishing to a subdirectory we remove the subdirectory
		if err := os.RemoveAll(targetDir); err != nil {
			return "", fmt.Errorf("cannot clean target dir: %w", err)
		}

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return "", fmt.Errorf("cannot create target dir: %w", err)
		}
	}

	p.Log().Info("Copying generated site to target directory")
	if err := copyDir(sourceDir, targetDir); err != nil {
		return "", fmt.Errorf("cannot copy site content: %w", err)
	}

	// Stage
	p.Log().Info("Staging changes")
	if err := p.gitClient.Add(ctx, tempDir, ".", env); err != nil {
		return "", fmt.Errorf("cannot stage changes: %w", err)
	}

	// Commit
	p.Log().Info("Committing changes")
	commitHash, err := p.gitClient.Commit(ctx, tempDir, cfg.CommitAuthor, env)
	if err != nil {
		return "", fmt.Errorf("cannot commit changes: %w", err)
	}
	p.Log().Info("Changes committed", "hash", commitHash)

	statusOutput, err := p.gitClient.Status(ctx, tempDir, env)
	if err != nil {
		p.Log().Error("cannot get git status after commit", "error", err)
	}

	p.Log().Info("DEBUG: git status after commit", "output", statusOutput)
	logOutput, err := p.gitClient.GitLog(ctx, tempDir, []string{"-1", "--pretty=format:%s"}, env)

	if err != nil {
		p.Log().Error("cannot get git log after commit", "error", err)
	}

	p.Log().Info("DEBUG: git log after commit", "output", logOutput)

	// Push
	p.Log().Info("Pushing changes to remote")
	if err := p.gitClient.Push(ctx, tempDir, cfg.Auth, "origin", cfg.Branch, env); err != nil {
		return "", fmt.Errorf("cannot push changes: %w", err)
	}

	// // NOTE: We need to find a neater way to do this
	commitURL = fmt.Sprintf("%s/commit/%s", cfg.RepoURL, commitHash)
	p.Log().Info("Publish process completed successfully", "commit_url", commitURL)

	return commitURL, nil
}

// Plan implementation
func (p *publisher) Plan(ctx context.Context, cfg PublisherConfig, sourceDir string) (PlanReport, error) {
	p.Log().Info("Starting plan dry-run process")

	var report PlanReport

	parentTempDir, err := os.MkdirTemp("", "clio-plan-parent-*")
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot create parent temp dir for plan: %w", err)
	}
	defer os.RemoveAll(parentTempDir)

	tempDir := filepath.Join(parentTempDir, "repo") // Git will create this

	askpassScriptPath := filepath.Join(parentTempDir, "git-askpass.sh")
	err = os.WriteFile(askpassScriptPath, []byte(fmt.Sprintf("#!/bin/sh\necho %s", cfg.Auth.Token)), 0700)
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot create askpass script for plan: %w", err)
	}

	env := os.Environ()
	env = append(env, "GIT_ASKPASS="+askpassScriptPath)

	if err := p.gitClient.Clone(ctx, cfg.RepoURL, tempDir, cfg.Auth, env); err != nil {
		return PlanReport{}, fmt.Errorf("cannot clone repo for plan: %w", err)
	}
	p.Log().Info("Repo cloned for plan")

	if err := p.gitClient.Checkout(ctx, tempDir, cfg.Branch, false, env); err != nil {
		return PlanReport{}, fmt.Errorf("cannot checkout branch for plan: %w", err)
	}
	p.Log().Info("Checked out branch for plan", "branch", cfg.Branch)

	targetDir := filepath.Join(tempDir, cfg.PagesSubdir)
	p.Log().Info("Cleaning target directory for plan", "path", targetDir)
	if err := os.RemoveAll(targetDir); err != nil {
		return PlanReport{}, fmt.Errorf("cannot clean target dir for plan: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return PlanReport{}, fmt.Errorf("cannot create target dir for plan: %w", err)
	}

	p.Log().Info("Copying generated site to target directory for plan")
	if err := copyDir(sourceDir, targetDir); err != nil {
		return PlanReport{}, fmt.Errorf("cannot copy site content for plan: %w", err)
	}

	p.Log().Info("Staging changes for plan")
	if err := p.gitClient.Add(ctx, tempDir, ".", env); err != nil {
		return PlanReport{}, fmt.Errorf("cannot stage changes for plan: %w", err)
	}

	p.Log().Info("Getting git status for plan")
	statusOutput, err := p.gitClient.Status(ctx, tempDir, env)
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot get git status for plan: %w", err)
	}

	lines := strings.Split(statusOutput, "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status := line[0:2]
		filename := strings.TrimSpace(line[3:])

		switch status {
		case "A ":
			report.Added = append(report.Added, filename)
		case "M ":
			report.Modified = append(report.Modified, filename)
		case "D ":
			report.Removed = append(report.Removed, filename)
		case "??":
			report.Added = append(report.Added, filename)
		}
	}

	report.Summary = fmt.Sprintf("Added: %d, Modified: %d, Removed: %d", len(report.Added), len(report.Modified), len(report.Removed))
	p.Log().Info("Plan dry-run process completed successfully", "summary", report.Summary)

	return report, nil
}

// copyDir copies the contents of src to dst. It is not recursive!
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// NOTE: calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
