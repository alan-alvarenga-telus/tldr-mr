package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// execCommand is the function used to create git commands. It is a variable so
// that tests can replace it with a fake implementation without spawning a real
// git process.
var execCommand = exec.Command

func GetCurrentBranch() (string, error) {
	cmd := execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
