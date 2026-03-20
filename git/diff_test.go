package git

import (
	"strings"
	"testing"
)

func TestGetDiff(t *testing.T) {
	original := execCommand
	t.Cleanup(func() { execCommand = original })

	t.Run("returns trimmed diff on success", func(t *testing.T) {
		fakeDiff := " main.go | 5 +++++\n 1 file changed, 5 insertions(+)\n\ndiff --git a/main.go b/main.go\n"
		execCommand = fakeExecCommand(fakeDiff, "", 0)

		diff, err := GetDiff("dev...feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := strings.TrimSpace(fakeDiff)
		if diff != expected {
			t.Errorf("unexpected diff output: %q", diff)
		}
	})

	t.Run("returns empty string when no changes", func(t *testing.T) {
		execCommand = fakeExecCommand("", "", 0)

		diff, err := GetDiff("dev...feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff != "" {
			t.Errorf("expected empty string, got %q", diff)
		}
	})

	t.Run("returns error when git command fails", func(t *testing.T) {
		execCommand = fakeExecCommand("", "fatal: not a git repository", 1)

		_, err := GetDiff("dev...feature")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "git diff failed") {
			t.Errorf("expected error to contain 'git diff failed', got %q", err.Error())
		}
	})

	t.Run("error message includes stderr output", func(t *testing.T) {
		execCommand = fakeExecCommand("", "bad revision 'dev...feature'", 1)

		_, err := GetDiff("dev...feature")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "bad revision 'dev...feature'") {
			t.Errorf("expected error to include stderr text, got %q", err.Error())
		}
	})

	t.Run("trims trailing whitespace and newlines", func(t *testing.T) {
		execCommand = fakeExecCommand("diff --git a/foo.go b/foo.go\n\n\n", "", 0)

		diff, err := GetDiff("main...feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff != "diff --git a/foo.go b/foo.go" {
			t.Errorf("expected trimmed output, got %q", diff)
		}
	})
}
