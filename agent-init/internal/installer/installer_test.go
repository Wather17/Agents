package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Wather17/Agents/agent-init/internal/templates"
)

func TestInstallCreatesFiles(t *testing.T) {
	target := t.TempDir()

	results, err := Install(Options{
		TargetPath: target,
		Agent:      templates.Gemini,
		Force:      false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 installed files, got %d", len(results))
	}

	if _, err := os.Stat(filepath.Join(target, "GEMINI.md")); err != nil {
		t.Errorf("GEMINI.md was not created: %v", err)
	}

	info, err := os.Stat(filepath.Join(target, "scripts", "sync-issues.sh"))
	if err != nil {
		t.Errorf("sync-issues.sh was not created: %v", err)
	} else if info.Mode().Perm()&0o111 == 0 {
		t.Errorf("sync-issues.sh should be executable")
	}

	gitignore, err := os.ReadFile(filepath.Join(target, ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore was not created: %v", err)
	}
	if !strings.Contains(string(gitignore), "GEMINI.md") || !strings.Contains(string(gitignore), "issues/") {
		t.Errorf(".gitignore missing expected entries: %s", gitignore)
	}
}

func TestInstallFailsWhenFilesExistWithoutForce(t *testing.T) {
	target := t.TempDir()

	if _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("first install should succeed: %v", err)
	}

	if _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err == nil {
		t.Fatal("second install without force should fail")
	}
}

func TestInstallSucceedsWithForce(t *testing.T) {
	target := t.TempDir()

	if _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("first install should succeed: %v", err)
	}

	if _, err := Install(Options{TargetPath: target, Agent: templates.Gemini, Force: true}); err != nil {
		t.Fatalf("install with force should succeed: %v", err)
	}
}

func TestInstallDoesNotDuplicateGitignoreEntries(t *testing.T) {
	target := t.TempDir()

	if err := os.WriteFile(filepath.Join(target, ".gitignore"), []byte("GEMINI.md\nissues/\n"), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	if _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("install should succeed: %v", err)
	}

	gitignore, err := os.ReadFile(filepath.Join(target, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	count := strings.Count(string(gitignore), "GEMINI.md")
	if count != 1 {
		t.Errorf("expected GEMINI.md once in .gitignore, found %d", count)
	}
}
