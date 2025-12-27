package ssg_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func TestPublisherPublishFakeClient(t *testing.T) {
	tests := []struct {
		name          string
		config        ssg.PublisherConfig
		sourceContent map[string]string
		gitClient     *fake.GithubClient
		expectedError error

		expectedCloneCalls    int
		expectedCheckoutCalls int
		expectedAddCalls      int
		expectedCommitCalls   int
		expectedPushCalls     int
	}{
		{
			name: "successful publish",
			config: ssg.PublisherConfig{
				RepoURL:     "https://github.com/test/repo.git",
				Branch:      "gh-pages",
				PagesSubdir: "",
				Auth:        hm.GitAuth{Method: hm.AuthToken, Token: "test-token"},
				CommitAuthor: hm.GitCommit{
					UserName:  "Test User",
					UserEmail: "test@example.com",
					Message:   "Test commit",
				},
			},
			sourceContent: map[string]string{
				"index.html":    "<html><body>Hello</body></html>",
				"css/style.css": "body { color: red; }",
			},
			gitClient: &fake.GithubClient{
				CloneFn: func(ctx context.Context, repoURL, localPath string, auth hm.GitAuth, env []string) error {
					if err := os.MkdirAll(localPath, 0755); err != nil {
						return err
					}
					return nil
				},
				CheckoutFn: func(ctx context.Context, localRepoPath, branch string, create bool, env []string) error { return nil },
				AddFn:      func(ctx context.Context, localRepoPath, pathspec string, env []string) error { return nil },
				CommitFn: func(ctx context.Context, localRepoPath string, commit hm.GitCommit, env []string) (string, error) {
					return "fake-hash", nil
				},
				PushFn: func(ctx context.Context, localRepoPath string, auth hm.GitAuth, remote, branch string, env []string) error {
					return nil
				},
			},
			expectedError:         nil,
			expectedCloneCalls:    1,
			expectedCheckoutCalls: 1,
			expectedAddCalls:      1,
			expectedCommitCalls:   1,
			expectedPushCalls:     1,
		},
		{
			name: "publish fails on clone",
			config: ssg.PublisherConfig{
				RepoURL:     "https://github.com/test/repo.git",
				Branch:      "gh-pages",
				PagesSubdir: "",
				Auth:        hm.GitAuth{Method: hm.AuthToken, Token: "test-token"},
				CommitAuthor: hm.GitCommit{
					UserName:  "Test User",
					UserEmail: "test@example.com",
					Message:   "Test commit",
				},
			},
			sourceContent: map[string]string{
				"index.html": "<html><body>Hello</body></html>",
			},
			gitClient: &fake.GithubClient{
				CloneFn: func(ctx context.Context, repoURL, localPath string, auth hm.GitAuth, env []string) error {
					return errors.New("clone error")
				},
			},
			expectedError:         errors.New("cannot clone repo: clone error"),
			expectedCloneCalls:    1,
			expectedCheckoutCalls: 0,
			expectedAddCalls:      0,
			expectedCommitCalls:   0,
			expectedPushCalls:     0,
		},
		// TODO: Add more cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceDir, err := os.MkdirTemp("", "test-source-*")
			if err != nil {
				t.Fatalf("cannot create temp source dir: %v", err)
			}
			defer os.RemoveAll(sourceDir)

			for filename, content := range tt.sourceContent {
				filePath := filepath.Join(sourceDir, filename)
				if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
					t.Fatalf("cannot create dir for source file: %v", err)
				}
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("cannot write source file: %v", err)
				}
			}

			params := hm.XParams{Log: hm.NewLogger("debug")}
			publisher := ssg.NewPublisher(tt.gitClient, params)

			_, err = publisher.Publish(context.Background(), tt.config, sourceDir)

			if tt.expectedError != nil {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError.Error()) {
					t.Errorf("Expected error containing \"%v\", got \"%v\"", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got \"%v\"", err)
			}

			if len(tt.gitClient.CloneCalls) != tt.expectedCloneCalls {
				t.Errorf("Expected %d Clone calls, got %d", tt.expectedCloneCalls, len(tt.gitClient.CloneCalls))
			}
			if len(tt.gitClient.CheckoutCalls) != tt.expectedCheckoutCalls {
				t.Errorf("Expected %d Checkout calls, got %d", tt.expectedCheckoutCalls, len(tt.gitClient.CheckoutCalls))
			}
			if len(tt.gitClient.AddCalls) != tt.expectedAddCalls {
				t.Errorf("Expected %d Add calls, got %d", tt.expectedAddCalls, len(tt.gitClient.AddCalls))
			}
			if len(tt.gitClient.CommitCalls) != tt.expectedCommitCalls {
				t.Errorf("Expected %d Commit calls, got %d", tt.expectedCommitCalls, len(tt.gitClient.CommitCalls))
			}
			if len(tt.gitClient.PushCalls) != tt.expectedPushCalls {
				t.Errorf("Expected %d Push calls, got %d", tt.expectedPushCalls, len(tt.gitClient.PushCalls))
			}

			if tt.expectedCloneCalls > 0 && len(tt.gitClient.CloneCalls) > 0 {
				call := tt.gitClient.CloneCalls[0]
				if call.RepoURL != tt.config.RepoURL {
					t.Errorf("CloneCalls[0].RepoURL expected %s, got %s", tt.config.RepoURL, call.RepoURL)
				}
			}
		})
	}
}
