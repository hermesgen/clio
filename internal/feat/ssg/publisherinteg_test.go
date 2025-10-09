package ssg_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hermesgen/hm"
	"github.com/hermesgen/hm/github"
	"github.com/hermesgen/clio/internal/feat/ssg"
)

var (
	testRepoOwner   string
	testRepoName    string
	testGithubToken string
	gitClient       hm.GitClient
)

func TestMain(m *testing.M) {
	testRepoOwner = os.Getenv("GITHUB_TEST_REPO_OWNER")
	testRepoName = os.Getenv("GITHUB_TEST_REPO_NAME")
	testGithubToken = os.Getenv("GITHUB_TEST_TOKEN")

	if testRepoOwner == "" || testRepoName == "" || testGithubToken == "" {
		fmt.Println("Skipping integration tests: GITHUB_TEST_REPO_OWNER, GITHUB_TEST_REPO_NAME, or GITHUB_TEST_TOKEN not set.")
		os.Exit(0)
	}

	// Initialize real GitHub client
	logger := hm.NewLogger("info")
	cfg := hm.NewConfig() // Basic config for tests
	xparams := hm.XParams{Cfg: cfg, Log: logger}
	gitClient = github.NewClient(xparams)

	// Prepare the remote repository
	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", testRepoOwner, testRepoName)
	testBranch := fmt.Sprintf("test-run-%d", time.Now().UnixNano())
	os.Setenv("TEST_BRANCH", testBranch) // Pass branch to tests

	setupDir, err := os.MkdirTemp("", "clio-test-setup-*")
	if err != nil {
		fmt.Printf("Failed to create temp dir for setup: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(setupDir)

	fmt.Println("Cloning repository for test setup...")
	cmd := exec.Command("git", "clone", repoURL, setupDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to clone repo for setup: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	fmt.Printf("Creating and checking out test branch: %s\n", testBranch)

	// We need to check if the branch already exists locally and delete it if it does
	cmd = exec.Command("git", "branch", "-D", testBranch)
	cmd.Dir = setupDir
	cmd.Run() // Branch might not exist, we can ignore this error

	cmd = exec.Command("git", "checkout", "-b", testBranch)
	cmd.Dir = setupDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to create test branch: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	fmt.Println("Cleaning repository by removing all files...")
	// NOTE: Remove all files and dirs except .git
	err = filepath.Walk(setupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		// We only remove contents of setupDir
		if strings.HasPrefix(path, setupDir) && path != setupDir {
			return os.RemoveAll(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Failed to clean repo contents: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = setupDir
	cmd.Run()

	commitMsg := "chore: prepare for integration test"
	cmd = exec.Command("git", "commit", "--allow-empty", "-m", commitMsg)
	cmd.Dir = setupDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to commit cleanup changes: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	fmt.Println("Pushing clean state to remote test branch...")
	cmd = exec.Command("git", "push", "-u", "origin", testBranch)
	cmd.Dir = setupDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to push clean test branch: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	fmt.Printf("Test setup complete. Running tests on branch: %s\n", testBranch)

	code := m.Run()

	fmt.Println("Running teardown: deleting remote test branch...")

	// Switch to main branch before deleting the test branch
	cmd = exec.Command("git", "checkout", "main")
	cmd.Dir = setupDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("WARNING: Failed to checkout main branch for teardown: %v\nOutput: %s\n", err, string(output))
		// NOTE We continue attempting to delete the branch even if checkout fails
	}

	cmd = exec.Command("git", "push", "origin", "--delete", testBranch)
	cmd.Dir = setupDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("WARNING: Failed to delete remote test branch %s: %v\nOutput: %s\n", testBranch, err, string(output))
	} else {
		fmt.Printf("Successfully deleted remote branch: %s\n", testBranch)
	}

	os.Exit(code)
}

func TestPublisherIntegration(t *testing.T) {
	if testRepoOwner == "" || testRepoName == "" || testGithubToken == "" {
		t.Skip("Skipping integration tests: GITHUB_TEST_REPO_OWNER, GITHUB_TEST_REPO_NAME, or GITHUB_TEST_TOKEN not set.")
	}

	testBranch := os.Getenv("TEST_BRANCH")

	tests := []struct {
		name                string
		initialContent      string
		publishContent      string
		publishPath         string
		commitMessage       string
		expectedFileContent string
		expectedError       error
	}{
		{
			name:                "Publish new file",
			initialContent:      "",
			publishContent:      "Hello, integration test!",
			publishPath:         "test-file.md",
			commitMessage:       "feat: add test-file.md",
			expectedFileContent: "Hello, integration test!",
			expectedError:       nil,
		},
		{
			name:                "Update existing file",
			initialContent:      "Initial content.",
			publishContent:      "Updated content.",
			publishPath:         "test-file-to-update.md",
			commitMessage:       "feat: update test-file-to-update.md",
			expectedFileContent: "Updated content.",
			expectedError:       nil,
		},
		// TODO: Add more scenarios: delete file, publish to different path, etc.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We clone the repo into a fresh working directory for each sub-test.
			parentWorkDir, err := os.MkdirTemp("", "clio-workdir-parent-*")
			require.NoError(t, err)
			defer os.RemoveAll(parentWorkDir)

			workDir := filepath.Join(parentWorkDir, "repo")

			askpassScriptPath := filepath.Join(parentWorkDir, "git-askpass.sh")
			err = os.WriteFile(askpassScriptPath, []byte(fmt.Sprintf("#!/bin/sh\necho %s", testGithubToken)), 0700)
			require.NoError(t, err)

			env := os.Environ()
			env = append(env, "GIT_ASKPASS="+askpassScriptPath)

			repoURL := fmt.Sprintf("https://github.com/%s/%s.git", testRepoOwner, testRepoName)
			cmd := exec.Command("git", "clone", "--branch", testBranch, "--single-branch", repoURL, workDir)
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Failed to clone repo into workDir: %s", string(output))

			authRepoURL := fmt.Sprintf("https://oauth2:%s@github.com/%s/%s.git", testGithubToken, testRepoOwner, testRepoName)

			if tt.initialContent != "" {
				err = os.WriteFile(filepath.Join(workDir, tt.publishPath), []byte(tt.initialContent), 0644)
				require.NoError(t, err)

				cmd = exec.Command("git", "add", ".")
				cmd.Dir = workDir
				cmd.Env = env
				_, err = cmd.CombinedOutput()
				require.NoError(t, err)

				cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("feat: add initial content for %s", t.Name()))
				cmd.Dir = workDir
				cmd.Env = env
				_, err = cmd.CombinedOutput()
				require.NoError(t, err)

				cmd = exec.Command("git", "push")
				cmd.Dir = workDir
				cmd.Env = env
				output, err = cmd.CombinedOutput()
				require.NoError(t, err, "Failed to push initial content: %s", string(output))
			}

			sourceDir, err := os.MkdirTemp("", "clio-source-*")
			require.NoError(t, err)
			defer os.RemoveAll(sourceDir)

			err = os.WriteFile(filepath.Join(sourceDir, tt.publishPath), []byte(tt.publishContent), 0644)
			require.NoError(t, err)

			logger := hm.NewLogger("info")
			publisher := ssg.NewPublisher(gitClient, hm.WithLog(logger))

			pubCfg := ssg.PublisherConfig{
				RepoURL:     authRepoURL, // NOTE: We wse the authenticated
				Branch:      testBranch,
				PagesSubdir: "",
				Auth: hm.GitAuth{
					Method: hm.AuthToken,
					Token:  testGithubToken,
				},
				CommitAuthor: hm.GitCommit{
					UserName:  "Clio Test Bot",
					UserEmail: "ci@clio.dev",
					Message:   tt.commitMessage,
				},
			}

			_, err = publisher.Publish(context.Background(), pubCfg, sourceDir)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)

				t.Logf("DEBUG: Verifying publisher's push...")
				checkDir, err := os.MkdirTemp("", "clio-check-*")
				require.NoError(t, err)
				defer os.RemoveAll(checkDir)

				// Clone the remote into check dir to see if the commit is there
				cmd = exec.Command("git", "clone", "--branch", testBranch, "--single-branch", authRepoURL, checkDir)
				output, err = cmd.CombinedOutput()
				require.NoError(t, err, "Failed to clone repo into checkDir: %s", string(output))

				cmd = exec.Command("git", "log", "-1", "--pretty=format:%s")
				cmd.Dir = checkDir
				output, err = cmd.CombinedOutput()
				require.NoError(t, err, "Cannot get git log from checkDir: %s", string(output))
				lastCommitMessage := strings.TrimSpace(string(output))
				assert.Equal(t, tt.commitMessage, lastCommitMessage, "Last commit message in remote does not match expected")
				t.Logf("DEBUG: Last commit message in remote: %s", lastCommitMessage)

				// We verify file content in the real repository by pulling changes into our work dir
				t.Logf("DEBUG: workDir before pull: %s", workDir)
				debugCmd := exec.Command("git", "status")
				debugCmd.Dir = workDir
				debugCmd.Env = env
				debugOutput, _ := debugCmd.CombinedOutput()
				t.Logf("DEBUG: git status in workDir:\n%s", string(debugOutput))
				debugCmd = exec.Command("ls", "-R")
				debugCmd.Dir = workDir
				debugCmd.Env = env
				debugOutput, _ = debugCmd.CombinedOutput()
				t.Logf("DEBUG: ls -R in workDir:\n%s", string(debugOutput))

				t.Logf("DEBUG: Waiting 2 seconds before pulling changes...")
				time.Sleep(2 * time.Second)

				cmd = exec.Command("git", "pull", "origin", testBranch)
				cmd.Dir = workDir
				cmd.Env = env
				output, err = cmd.CombinedOutput()
				require.NoError(t, err, "Failed to pull changes for verification: %s", string(output))

				t.Logf("DEBUG: workDir after pull: %s", workDir)
				debugCmd = exec.Command("git", "status")
				debugCmd.Dir = workDir
				debugCmd.Env = env
				debugOutput, _ = debugCmd.CombinedOutput()
				t.Logf("DEBUG: git status in workDir after pull:\n%s", string(debugOutput))
				debugCmd = exec.Command("ls", "-R")
				debugCmd.Dir = workDir
				debugCmd.Env = env
				debugOutput, _ = debugCmd.CombinedOutput()
				t.Logf("DEBUG: ls -R in workDir after pull:\n%s", string(debugOutput))

				content, err := os.ReadFile(filepath.Join(workDir, tt.publishPath))
				require.NoError(t, err, "Failed to read file from verified repo")

				assert.Equal(t, tt.expectedFileContent, string(content), "Verified file content does not match expected content")
			}
		})
	}
}
