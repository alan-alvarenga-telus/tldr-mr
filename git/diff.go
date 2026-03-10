package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetDiff(comparison string) (string, error) {
	cmd := exec.Command("git", "diff", comparison, "--histogram", "--unified=5", "-M", "--stat", "-w")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff failed: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
