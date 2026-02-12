package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetCurrentBranch returns the name of the current git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// GetLog executes git log for the given branch comparison
func GetLog(comparison string) (string, error) {
	cmd := exec.Command("git", "log", comparison, "--oneline")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git log failed: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// GetDiff executes git diff for the given branch comparison
func GetDiff(comparison string) (string, error) {
	cmd := exec.Command("git", "diff", comparison, "--histogram", "--unified=4", "-M", "--stat", "-w")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff failed: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
