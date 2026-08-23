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

	installed, skipped, err := Install(Options{
		TargetPath: target,
		Agent:      templates.Gemini,
		Force:      false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(installed) != 4 {
		t.Fatalf("expected 4 installed files, got %d", len(installed))
	}
	if len(skipped) != 0 {
		t.Fatalf("expected 0 skipped files, got %d", len(skipped))
	}

	if _, err := os.Stat(filepath.Join(target, "GEMINI.md")); err != nil {
		t.Errorf("GEMINI.md was not created: %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, "agents", "issue-architect.md")); err != nil {
		t.Errorf("agents/issue-architect.md was not created: %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, ".agents", "skills", "refine-issues.md")); err != nil {
		t.Errorf(".agents/skills/refine-issues.md was not created: %v", err)
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
	for _, entry := range []string{"GEMINI.md", "issues/", "agents/issue-architect.md", ".agents/skills/refine-issues.md"} {
		if !strings.Contains(string(gitignore), entry) {
			t.Errorf(".gitignore missing expected entry %q: %s", entry, gitignore)
		}
	}
}

func TestInstallSkipsIdenticalFiles(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("first install should succeed: %v", err)
	}

	installed, skipped, err := Install(Options{TargetPath: target, Agent: templates.Gemini})
	if err != nil {
		t.Fatalf("second install of identical files should succeed: %v", err)
	}
	if len(installed) != 0 {
		t.Fatalf("expected 0 installed files, got %d", len(installed))
	}
	if len(skipped) != 4 {
		t.Fatalf("expected 4 skipped files, got %d", len(skipped))
	}
}

func TestInstallFailsWhenFilesDifferWithoutForce(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("first install should succeed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(target, "GEMINI.md"), []byte("different content"), 0o644); err != nil {
		t.Fatalf("failed to modify GEMINI.md: %v", err)
	}

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err == nil {
		t.Fatal("install should fail when file differs without force")
	}
}

func TestInstallSucceedsWithForce(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("first install should succeed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(target, "GEMINI.md"), []byte("different content"), 0o644); err != nil {
		t.Fatalf("failed to modify GEMINI.md: %v", err)
	}

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini, Force: true}); err != nil {
		t.Fatalf("install with force should succeed: %v", err)
	}
}

func TestInstallDoesNotDuplicateGitignoreEntries(t *testing.T) {
	target := t.TempDir()

	if err := os.WriteFile(filepath.Join(target, ".gitignore"), []byte("GEMINI.md\nissues/\n"), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
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

func TestInstallCreatesOpenCodeFiles(t *testing.T) {
	target := t.TempDir()

	installed, skipped, err := Install(Options{
		TargetPath: target,
		Agent:      templates.OpenCode,
		Force:      false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(installed) != 4 {
		t.Fatalf("expected 4 installed files, got %d", len(installed))
	}
	if len(skipped) != 0 {
		t.Fatalf("expected 0 skipped files, got %d", len(skipped))
	}

	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); err != nil {
		t.Errorf("AGENTS.md was not created: %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, "agents", "issue-architect.md")); err != nil {
		t.Errorf("agents/issue-architect.md was not created: %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, ".agents", "skills", "refine-issues.md")); err != nil {
		t.Errorf(".agents/skills/refine-issues.md was not created: %v", err)
	}

	gitignore, err := os.ReadFile(filepath.Join(target, ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore was not created: %v", err)
	}
	for _, entry := range []string{"AGENTS.md", "issues/", "agents/issue-architect.md", ".agents/skills/refine-issues.md"} {
		if !strings.Contains(string(gitignore), entry) {
			t.Errorf(".gitignore missing expected entry %q: %s", entry, gitignore)
		}
	}
}

func TestInstallMultipleAgents(t *testing.T) {
	target := t.TempDir()

	installed, skipped, err := Install(Options{TargetPath: target, Agent: templates.Gemini})
	if err != nil {
		t.Fatalf("gemini install should succeed: %v", err)
	}
	if len(installed) != 4 || len(skipped) != 0 {
		t.Fatalf("unexpected gemini install result: installed=%d skipped=%d", len(installed), len(skipped))
	}

	installed, skipped, err = Install(Options{TargetPath: target, Agent: templates.OpenCode})
	if err != nil {
		t.Fatalf("opencode install after gemini should succeed: %v", err)
	}

	if len(installed) != 1 {
		t.Fatalf("expected 1 installed file for opencode, got %d", len(installed))
	}
	if installed[0].Path != filepath.Join(target, "AGENTS.md") {
		t.Errorf("expected AGENTS.md to be installed, got %s", installed[0].Path)
	}

	if len(skipped) != 3 {
		t.Fatalf("expected 3 skipped files for opencode, got %d", len(skipped))
	}
	foundSync := false
	for _, s := range skipped {
		if s.Path == filepath.Join(target, "scripts", "sync-issues.sh") {
			foundSync = true
		}
	}
	if !foundSync {
		t.Error("expected scripts/sync-issues.sh to be skipped")
	}

	gitignore, err := os.ReadFile(filepath.Join(target, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}
	if !strings.Contains(string(gitignore), "GEMINI.md") || !strings.Contains(string(gitignore), "AGENTS.md") {
		t.Errorf(".gitignore should contain both GEMINI.md and AGENTS.md: %s", gitignore)
	}
	if !strings.Contains(string(gitignore), "agents/issue-architect.md") || !strings.Contains(string(gitignore), ".agents/skills/refine-issues.md") {
		t.Errorf(".gitignore should contain skill and agent entries: %s", gitignore)
	}
}

func TestUpgradeUpdatesExistingPrompts(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("install should succeed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(target, "GEMINI.md"), []byte("modified"), 0o644); err != nil {
		t.Fatalf("failed to modify GEMINI.md: %v", err)
	}

	installed, skipped, err := Upgrade(target)
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}

	if len(installed) != 1 {
		t.Fatalf("expected 1 updated file, got %d", len(installed))
	}
	if installed[0].Path != filepath.Join(target, "GEMINI.md") {
		t.Errorf("expected GEMINI.md to be updated, got %s", installed[0].Path)
	}
	if len(skipped) != 2 {
		t.Fatalf("expected 2 skipped files, got %d", len(skipped))
	}

	content, err := os.ReadFile(filepath.Join(target, "GEMINI.md"))
	if err != nil {
		t.Fatalf("failed to read GEMINI.md: %v", err)
	}
	if strings.Contains(string(content), "modified") {
		t.Errorf("GEMINI.md was not updated to template content")
	}
}

func TestUpgradeUpdatesSkillAndAgent(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("install should succeed: %v", err)
	}

	for _, path := range []string{
		filepath.Join(target, "agents", "issue-architect.md"),
		filepath.Join(target, ".agents", "skills", "refine-issues.md"),
	} {
		if err := os.WriteFile(path, []byte("modified"), 0o644); err != nil {
			t.Fatalf("failed to modify %s: %v", path, err)
		}
	}

	installed, skipped, err := Upgrade(target)
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}

	if len(installed) != 2 {
		t.Fatalf("expected 2 updated files, got %d", len(installed))
	}
	for _, r := range installed {
		if r.Ignored != true {
			t.Errorf("updated file %s should be marked as ignored", r.Path)
		}
		content, err := os.ReadFile(r.Path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", r.Path, err)
		}
		if strings.Contains(string(content), "modified") {
			t.Errorf("%s was not updated to template content", r.Path)
		}
	}
	if len(skipped) != 1 {
		t.Fatalf("expected 1 skipped file (GEMINI.md), got %d", len(skipped))
	}
}

func TestUpgradeInstallsMissingAuxiliaryFilesForExistingAgents(t *testing.T) {
	tests := []struct {
		name          string
		promptSource  string
		promptTarget  string
		missingTarget string
	}{
		{
			name:          "gemini",
			promptSource:  "files/GEMINI.md",
			promptTarget:  "GEMINI.md",
			missingTarget: "AGENTS.md",
		},
		{
			name:          "opencode",
			promptSource:  "files/AGENTS.md",
			promptTarget:  "AGENTS.md",
			missingTarget: "GEMINI.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := t.TempDir()
			content, err := templates.Read(tt.promptSource)
			if err != nil {
				t.Fatalf("failed to read prompt template: %v", err)
			}
			if err := os.WriteFile(filepath.Join(target, tt.promptTarget), content, 0o644); err != nil {
				t.Fatalf("failed to create existing prompt: %v", err)
			}

			installed, skipped, err := Upgrade(target)
			if err != nil {
				t.Fatalf("upgrade should succeed: %v", err)
			}
			if len(installed) != 2 {
				t.Fatalf("expected 2 auxiliary files to be installed, got %d", len(installed))
			}
			if len(skipped) != 1 {
				t.Fatalf("expected the existing prompt to be skipped, got %d skipped files", len(skipped))
			}

			for _, path := range []string{
				filepath.Join(target, "agents", "issue-architect.md"),
				filepath.Join(target, ".agents", "skills", "refine-issues.md"),
			} {
				if _, err := os.Stat(path); err != nil {
					t.Errorf("expected auxiliary file %s to be installed: %v", path, err)
				}
			}
			if _, err := os.Stat(filepath.Join(target, tt.missingTarget)); err == nil {
				t.Errorf("upgrade should not install the new agent prompt %s", tt.missingTarget)
			}

			gitignore, err := os.ReadFile(filepath.Join(target, ".gitignore"))
			if err != nil {
				t.Fatalf("upgrade should create .gitignore: %v", err)
			}
			for _, entry := range []string{tt.promptTarget, "agents/issue-architect.md", ".agents/skills/refine-issues.md", "issues/"} {
				if !strings.Contains(string(gitignore), entry) {
					t.Errorf(".gitignore missing expected entry %q: %s", entry, gitignore)
				}
			}
		})
	}
}

func TestUpgradeDoesNotInstallFilesWithoutAnAgentPrompt(t *testing.T) {
	target := t.TempDir()

	installed, skipped, err := Upgrade(target)
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}
	if len(installed) != 0 || len(skipped) != 0 {
		t.Fatalf("upgrade without an installed agent should do nothing: installed=%d skipped=%d", len(installed), len(skipped))
	}

	for _, path := range []string{
		filepath.Join(target, "agents", "issue-architect.md"),
		filepath.Join(target, ".agents", "skills", "refine-issues.md"),
	} {
		if _, err := os.Stat(path); err == nil {
			t.Errorf("upgrade should not install auxiliary file %s without an agent prompt", path)
		}
	}
}

func TestUpgradeDoesNotInstallNewAgents(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("install should succeed: %v", err)
	}

	installed, skipped, err := Upgrade(target)
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}

	if len(installed) != 0 {
		t.Fatalf("expected 0 updated files, got %d", len(installed))
	}
	if len(skipped) != 3 {
		t.Fatalf("expected 3 skipped files, got %d", len(skipped))
	}

	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); err == nil {
		t.Error("AGENTS.md should not be created by upgrade")
	}
}

func TestUpgradeDoesNotTouchSharedScripts(t *testing.T) {
	target := t.TempDir()

	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Gemini}); err != nil {
		t.Fatalf("install should succeed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(target, "scripts", "sync-issues.sh"), []byte("#!/bin/bash\necho custom"), 0o755); err != nil {
		t.Fatalf("failed to modify sync-issues.sh: %v", err)
	}

	installed, skipped, err := Upgrade(target)
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}

	if len(installed) != 0 {
		t.Fatalf("expected 0 updated files, got %d", len(installed))
	}
	if len(skipped) != 3 {
		t.Fatalf("expected 3 skipped files, got %d", len(skipped))
	}

	content, err := os.ReadFile(filepath.Join(target, "scripts", "sync-issues.sh"))
	if err != nil {
		t.Fatalf("failed to read sync-issues.sh: %v", err)
	}
	if !strings.Contains(string(content), "custom") {
		t.Error("sync-issues.sh should not be modified by upgrade")
	}
}
