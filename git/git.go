package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
	cmd := exec.Command("git", "diff", comparison, "--histogram --unified=4 -M --stats -w")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff failed: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
