package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
