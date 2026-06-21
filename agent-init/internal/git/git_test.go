package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsRepository(t *testing.T) {
	tmp := t.TempDir()

	if IsRepository(tmp) {
		t.Error("non-git directory should not be a repository")
	}

	if err := exec.Command("git", "init", "-q", tmp).Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	if !IsRepository(tmp) {
		t.Error("git directory should be a repository")
	}
}

func TestHasChanges(t *testing.T) {
	tmp := t.TempDir()
	if err := exec.Command("git", "init", "-q", tmp).Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	changed, err := HasChanges(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("empty repo should have no changes")
	}

	if err := os.WriteFile(filepath.Join(tmp, "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	changed, err = HasChanges(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("repo with new file should have changes")
	}
}

func TestCommit(t *testing.T) {
	tmp := t.TempDir()
	if err := exec.Command("git", "init", "-q", tmp).Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}
	if err := exec.Command("git", "-C", tmp, "config", "user.email", "test@example.com").Run(); err != nil {
		t.Fatalf("failed to config user.email: %v", err)
	}
	if err := exec.Command("git", "-C", tmp, "config", "user.name", "Test").Run(); err != nil {
		t.Fatalf("failed to config user.name: %v", err)
	}

	file := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	if err := Commit(tmp, []string{"file.txt"}, "chore(test): initial commit"); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	out, err := exec.Command("git", "-C", tmp, "log", "--oneline", "-1").Output()
	if err != nil {
		t.Fatalf("failed to read log: %v", err)
	}
	if !strings.Contains(string(out), "chore(test): initial commit") {
		t.Errorf("commit message not found in log: %s", out)
	}
}
