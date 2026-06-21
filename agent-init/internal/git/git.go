// Package git provides helpers to interact with a git repository.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsRepository returns true if the given path is inside a git repository.
func IsRepository(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	gitDir := strings.TrimSpace(string(out))
	return filepath.IsAbs(gitDir) || filepath.IsLocal(gitDir)
}

// HasChanges returns true if the repository at path has uncommitted changes.
func HasChanges(path string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("running git status: %w", err)
	}
	return len(strings.TrimSpace(string(out))) > 0, nil
}

// Commit creates a commit with the given files using a Conventional Commit message.
func Commit(path string, files []string, message string) error {
	if message == "" {
		message = "chore(agent-init): add agent configuration and issue sync script"
	}

	cmd := exec.Command("git", append([]string{"add", "--"}, files...)...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}
