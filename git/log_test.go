package git

import (
	"strings"
	"testing"
)

func TestGetLog(t *testing.T) {
	original := execCommand
	t.Cleanup(func() { execCommand = original })

	t.Run("returns trimmed commit log on success", func(t *testing.T) {
		execCommand = fakeExecCommand("abc1234 fix: correct typo\ndef5678 feat: add login\n", "", 0)

		log, err := GetLog("dev...feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if log != "abc1234 fix: correct typo\ndef5678 feat: add login" {
			t.Errorf("unexpected log output: %q", log)
		}
	})

	t.Run("trims trailing whitespace and newlines", func(t *testing.T) {
		execCommand = fakeExecCommand("abc1234 single commit\n\n", "", 0)

		log, err := GetLog("main...feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if log != "abc1234 single commit" {
			t.Errorf("expected trimmed output, got %q", log)
		}
	})

	t.Run("returns empty string when no commits", func(t *testing.T) {
		execCommand = fakeExecCommand("", "", 0)

		log, err := GetLog("dev...feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if log != "" {
			t.Errorf("expected empty string, got %q", log)
		}
	})

	t.Run("returns error when git command fails", func(t *testing.T) {
		execCommand = fakeExecCommand("", "fatal: not a git repository", 1)

		_, err := GetLog("dev...feature")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "git log failed") {
			t.Errorf("expected error to contain 'git log failed', got %q", err.Error())
		}
	})

	t.Run("error message includes stderr output", func(t *testing.T) {
		execCommand = fakeExecCommand("", "ambiguous argument 'dev...feature'", 1)

		_, err := GetLog("dev...feature")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "ambiguous argument 'dev...feature'") {
			t.Errorf("expected error to include stderr text, got %q", err.Error())
		}
	})
}
