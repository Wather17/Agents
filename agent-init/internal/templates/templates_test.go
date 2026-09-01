package templates

import "testing"

func TestParseAgent(t *testing.T) {
	tests := map[string]Agent{
		"antigravity": Antigravity,
		"gemini":      Antigravity,
		"opencode":    OpenCode,
		"codex":       Codex,
		"claude":      Claude,
	}
	for input, expected := range tests {
		agent, err := ParseAgent(input)
		if err != nil {
			t.Fatalf("ParseAgent(%q) returned an error: %v", input, err)
		}
		if agent != expected {
			t.Errorf("ParseAgent(%q) = %q, want %q", input, agent, expected)
		}
	}
	if _, err := ParseAgent("unknown"); err == nil {
		t.Fatal("expected unsupported agent to fail")
	}
}

func TestFilesForGemini(t *testing.T) {
	files, err := FilesFor(Gemini)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 5 {
		t.Fatalf("expected 5 files, got %d", len(files))
	}

	if files[0].TargetPath != "GEMINI.md" || !files[0].Ignored || !files[0].Prompt {
		t.Errorf("GEMINI.md should be ignored, got %+v", files[0])
	}

	if files[1].TargetPath != ".agents/skills/audit-issues/SKILL.md" || !files[1].Ignored || files[1].Prompt {
		t.Errorf("audit-issues SKILL.md should be ignored, got %+v", files[1])
	}

	if files[2].TargetPath != ".agents/skills/refine-issues/SKILL.md" || !files[2].Ignored || files[2].Prompt {
		t.Errorf("refine-issues SKILL.md should be ignored, got %+v", files[2])
	}

	if files[3].TargetPath != ".agents/skills/autonomous-batch/SKILL.md" || !files[3].Ignored || files[3].Prompt {
		t.Errorf("autonomous-batch SKILL.md should be ignored, got %+v", files[3])
	}

	if files[4].TargetPath != "scripts/sync-issues.sh" || files[4].Ignored {
		t.Errorf("sync-issues.sh should not be ignored, got %+v", files[4])
	}
}

func TestFilesForUnsupportedAgent(t *testing.T) {
	_, err := FilesFor("unknown")
	if err == nil {
		t.Fatal("expected error for unsupported agent")
	}

	_, err = IgnoredEntries("unknown")
	if err == nil {
		t.Fatal("expected error for unsupported agent ignored entries")
	}
}

func TestFilesForOpenCode(t *testing.T) {
	files, err := FilesFor(OpenCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 5 {
		t.Fatalf("expected 5 files, got %d", len(files))
	}

	if files[0].TargetPath != "AGENTS.md" || !files[0].Ignored || !files[0].Prompt {
		t.Errorf("AGENTS.md should be ignored, got %+v", files[0])
	}

	if files[1].TargetPath != ".opencode/skill/audit-issues/SKILL.md" || !files[1].Ignored || files[1].Prompt {
		t.Errorf("audit-issues SKILL.md should be ignored, got %+v", files[1])
	}

	if files[2].TargetPath != ".opencode/skill/refine-issues/SKILL.md" || !files[2].Ignored || files[2].Prompt {
		t.Errorf("refine-issues SKILL.md should be ignored, got %+v", files[2])
	}

	if files[3].TargetPath != ".opencode/skill/autonomous-batch/SKILL.md" || !files[3].Ignored || files[3].Prompt {
		t.Errorf("autonomous-batch SKILL.md should be ignored, got %+v", files[3])
	}

	if files[4].TargetPath != "scripts/sync-issues.sh" || files[4].Ignored {
		t.Errorf("sync-issues.sh should not be ignored, got %+v", files[4])
	}

	entries, err := IgnoredEntries(OpenCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, entry := range entries {
		if entry == "AGENTS.md" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected AGENTS.md in ignored entries, got %v", entries)
	}
}

func TestFilesForCodexAndClaude(t *testing.T) {
	tests := []struct {
		agent       Agent
		prompt      string
		skillPrefix string
	}{
		{agent: Codex, prompt: "AGENTS.md", skillPrefix: ".agents/skills"},
		{agent: Claude, prompt: "CLAUDE.md", skillPrefix: ".claude/skills"},
	}
	for _, tt := range tests {
		files, err := FilesFor(tt.agent)
		if err != nil {
			t.Fatalf("FilesFor(%s) returned an error: %v", tt.agent, err)
		}
		if len(files) != 5 {
			t.Fatalf("FilesFor(%s) returned %d files, want 5", tt.agent, len(files))
		}
		if files[0].TargetPath != tt.prompt || !files[0].Prompt {
			t.Errorf("unexpected prompt for %s: %+v", tt.agent, files[0])
		}
		if files[1].TargetPath != tt.skillPrefix+"/audit-issues/SKILL.md" {
			t.Errorf("unexpected skill path for %s: %s", tt.agent, files[1].TargetPath)
		}
	}
}

func TestIgnoredEntriesIncludeSkillAndAgent(t *testing.T) {
	expectedByAgent := map[Agent][]string{
		Gemini: {
			"GEMINI.md",
			".agents/skills/audit-issues/SKILL.md",
			".agents/skills/refine-issues/SKILL.md",
			".agents/skills/autonomous-batch/SKILL.md",
			"issues/",
		},
		OpenCode: {
			"AGENTS.md",
			".opencode/skill/audit-issues/SKILL.md",
			".opencode/skill/refine-issues/SKILL.md",
			".opencode/skill/autonomous-batch/SKILL.md",
			"issues/",
		},
	}

	for agent, expected := range expectedByAgent {
		entries, err := IgnoredEntries(agent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, want := range expected {
			found := false
			for _, entry := range entries {
				if entry == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected %q in ignored entries for %s, got %v", want, agent, entries)
			}
		}
	}
}

func TestReadEmbeddedFiles(t *testing.T) {
	agents := []Agent{Antigravity, OpenCode, Codex, Claude}
	for _, agent := range agents {
		files, _ := FilesFor(agent)
		for _, file := range files {
			data, err := Read(file.SourcePath)
			if err != nil {
				t.Errorf("failed to read %q: %v", file.SourcePath, err)
				continue
			}
			if len(data) == 0 {
				t.Errorf("embedded file %q is empty", file.SourcePath)
			}
		}
	}
}
