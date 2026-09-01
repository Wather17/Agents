package installer

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Wather17/Agents/agent-init/internal/templates"
)

func TestManifestTracksCanonicalAgentsInOrder(t *testing.T) {
	target := t.TempDir()
	for _, agent := range []templates.Agent{templates.OpenCode, templates.Antigravity, templates.Claude, templates.Codex} {
		if _, _, err := Install(Options{TargetPath: target, Agent: agent}); err != nil {
			t.Fatalf("installing %s: %v", agent, err)
		}
	}

	state, exists, err := readManifest(target)
	if err != nil || !exists {
		t.Fatalf("reading manifest: exists=%v err=%v", exists, err)
	}
	want := []templates.Agent{templates.Antigravity, templates.Claude, templates.Codex, templates.OpenCode}
	if !reflect.DeepEqual(state.Agents, want) {
		t.Fatalf("manifest agents = %v, want %v", state.Agents, want)
	}

	data, err := os.ReadFile(filepath.Join(target, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), manifestFilename) != 1 {
		t.Fatalf("expected %s once in .gitignore: %s", manifestFilename, data)
	}
}

func TestInstallRejectsInvalidManifestBeforeWritingTemplates(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, manifestFilename), []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Install(Options{TargetPath: target, Agent: templates.Codex}); err == nil {
		t.Fatal("expected invalid manifest to fail")
	}
	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("AGENTS.md should not be written, got err=%v", err)
	}
}

func TestUpgradeMigratesLegacyPrompts(t *testing.T) {
	target := t.TempDir()
	for _, name := range []string{"GEMINI.md", "AGENTS.md"} {
		if err := os.WriteFile(filepath.Join(target, name), []byte("legacy"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := Upgrade(target); err != nil {
		t.Fatal(err)
	}
	state, exists, err := readManifest(target)
	if err != nil || !exists {
		t.Fatalf("reading migrated manifest: exists=%v err=%v", exists, err)
	}
	want := []templates.Agent{templates.Antigravity, templates.OpenCode}
	if !reflect.DeepEqual(state.Agents, want) {
		t.Fatalf("migrated agents = %v, want %v", state.Agents, want)
	}
	if _, err := os.Stat(filepath.Join(target, ".opencode/skill/audit-issues/SKILL.md")); err != nil {
		t.Fatalf("legacy AGENTS.md should migrate as opencode: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, ".claude/skills/audit-issues/SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("upgrade should not infer claude, got err=%v", err)
	}
}
