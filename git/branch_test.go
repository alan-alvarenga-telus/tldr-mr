package git

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// fakeExecCommand returns a *exec.Cmd that, when Run(), re-invokes the test
// binary itself with the special flag -test.run=TestHelperProcess. The helper
// process reads FAKE_STDOUT / FAKE_STDERR / FAKE_EXIT from the environment and
// writes them to the appropriate file descriptors before exiting.
func fakeExecCommand(stdout, stderr string, exitCode int) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			"FAKE_STDOUT="+stdout,
			"FAKE_STDERR="+stderr,
			"FAKE_EXIT="+exitCodeStr(exitCode),
		)
		return cmd
	}
}

func exitCodeStr(code int) string {
	if code == 0 {
		return "0"
	}
	return "1"
}

// TestHelperProcess is NOT a real test. It is invoked as a subprocess by the
// fake exec.Command created above. It must be present in the test binary.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	stdout := os.Getenv("FAKE_STDOUT")
	stderr := os.Getenv("FAKE_STDERR")
	exitStr := os.Getenv("FAKE_EXIT")

	if stdout != "" {
		os.Stdout.WriteString(stdout) //nolint:errcheck
	}
	if stderr != "" {
		os.Stderr.WriteString(stderr) //nolint:errcheck
	}
	if exitStr == "1" {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestGetCurrentBranch(t *testing.T) {
	// Save and restore the real exec.Command after each test.
	original := execCommand
	t.Cleanup(func() { execCommand = original })

	t.Run("returns trimmed branch name on success", func(t *testing.T) {
		execCommand = fakeExecCommand("feature/my-branch\n", "", 0)

		branch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "feature/my-branch" {
			t.Errorf("expected 'feature/my-branch', got %q", branch)
		}
	})

	t.Run("trims leading and trailing whitespace", func(t *testing.T) {
		execCommand = fakeExecCommand("  main  \n", "", 0)

		branch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "main" {
			t.Errorf("expected 'main', got %q", branch)
		}
	})

	t.Run("returns empty string when output is only whitespace", func(t *testing.T) {
		execCommand = fakeExecCommand("   \n", "", 0)

		branch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "" {
			t.Errorf("expected empty string, got %q", branch)
		}
	})

	t.Run("returns error when git command fails", func(t *testing.T) {
		execCommand = fakeExecCommand("", "not a git repository", 1)

		branch, err := GetCurrentBranch()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if branch != "" {
			t.Errorf("expected empty branch on error, got %q", branch)
		}
		if !strings.Contains(err.Error(), "failed to get current branch") {
			t.Errorf("expected error to contain 'failed to get current branch', got %q", err.Error())
		}
	})

	t.Run("error message includes stderr output", func(t *testing.T) {
		execCommand = fakeExecCommand("", "fatal: not a git repo", 1)

		_, err := GetCurrentBranch()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "fatal: not a git repo") {
			t.Errorf("expected error to include stderr text, got %q", err.Error())
		}
	})

	t.Run("handles branch name with slashes", func(t *testing.T) {
		execCommand = fakeExecCommand("release/v2.0.1\n", "", 0)

		branch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "release/v2.0.1" {
			t.Errorf("expected 'release/v2.0.1', got %q", branch)
		}
	})

	t.Run("HEAD branch name is returned as-is", func(t *testing.T) {
		execCommand = fakeExecCommand("HEAD\n", "", 0)

		branch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "HEAD" {
			t.Errorf("expected 'HEAD', got %q", branch)
		}
	})
}
