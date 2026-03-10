package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileWithDefault(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("explicit path returns content", func(t *testing.T) {
		path := filepath.Join(tmpDir, "explicit.txt")
		os.WriteFile(path, []byte("test content"), 0644)

		content, err := loadFileWithDefault(path, "default.txt", "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if content != "test content" {
			t.Errorf("expected 'test content', got %q", content)
		}
	})

	t.Run("explicit path not found returns error", func(t *testing.T) {
		_, err := loadFileWithDefault(filepath.Join(tmpDir, "nonexistent.txt"), "default.txt", "test")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("default file used when explicit is empty", func(t *testing.T) {
		path := filepath.Join(tmpDir, "default.txt")
		os.WriteFile(path, []byte("default content"), 0644)

		content, err := loadFileWithDefault("", path, "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if content != "default content" {
			t.Errorf("expected 'default content', got %q", content)
		}
	})

	t.Run("returns empty when no file exists", func(t *testing.T) {
		content, err := loadFileWithDefault("", filepath.Join(tmpDir, "nonexistent.txt"), "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if content != "" {
			t.Errorf("expected empty string, got %q", content)
		}
	})
}

func TestDetermineBranchComparison(t *testing.T) {
	mockResolver := &mockBranchResolver{}

	t.Run("legacy positional arg returns as-is", func(t *testing.T) {
		result, err := determineBranchComparison([]string{"main..feature"}, mockResolver)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "main..feature" {
			t.Errorf("expected 'main..feature', got %q", result)
		}
	})

	t.Run("uses headFlag when provided", func(t *testing.T) {
		headFlag = "my-feature"
		defer func() { headFlag = "" }()

		result, err := determineBranchComparison([]string{}, mockResolver)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "dev...my-feature" {
			t.Errorf("expected 'dev...my-feature', got %q", result)
		}
	})

	t.Run("falls back to resolver when headFlag empty", func(t *testing.T) {
		headFlag = ""
		mockResolver.branch = "current-branch"

		result, err := determineBranchComparison([]string{}, mockResolver)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "dev...current-branch" {
			t.Errorf("expected 'dev...current-branch', got %q", result)
		}
	})

	t.Run("uses custom base branch", func(t *testing.T) {
		headFlag = "feature"
		baseFlag = "main"
		defer func() { baseFlag = "dev" }()

		result, err := determineBranchComparison([]string{}, mockResolver)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "main...feature" {
			t.Errorf("expected 'main...feature', got %q", result)
		}
	})
}

type mockBranchResolver struct {
	branch string
	err    error
}

func (m *mockBranchResolver) GetCurrentBranch() (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.branch, nil
}
