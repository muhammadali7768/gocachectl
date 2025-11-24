package cache

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetGoEnv retrieves a Go environment variable
func GetGoEnv(key string) (string, error) {
	cmd := exec.Command("go", "env", key)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run 'go env %s': %w", key, err)
	}

	value := strings.TrimSpace(string(output))
	if value == "" {
		return "", fmt.Errorf("%s is empty", key)
	}

	return value, nil
}

// GetGoVersion retrieves the Go version
func GetGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run 'go version': %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatCount formats a count with thousands separator
func FormatCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}
	return addCommas(count)
}

func addCommas(n int) string {
	str := fmt.Sprintf("%d", n)
	var result strings.Builder
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}
	return result.String()
}
