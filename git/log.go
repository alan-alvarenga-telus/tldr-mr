package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
