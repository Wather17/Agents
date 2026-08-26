package templates

import "testing"

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

	if files[1].TargetPath != "agents/issue-architect.md" || !files[1].Ignored || files[1].Prompt {
		t.Errorf("issue-architect.md should be ignored, got %+v", files[1])
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

	if files[1].TargetPath != ".opencode/agent/issue-architect.md" || !files[1].Ignored || files[1].Prompt {
		t.Errorf("issue-architect agent should be ignored, got %+v", files[1])
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

func TestIgnoredEntriesIncludeSkillAndAgent(t *testing.T) {
	expectedByAgent := map[Agent][]string{
		Gemini: {
			"GEMINI.md",
			"agents/issue-architect.md",
			".agents/skills/refine-issues/SKILL.md",
			".agents/skills/autonomous-batch/SKILL.md",
			"issues/",
		},
		OpenCode: {
			"AGENTS.md",
			".opencode/agent/issue-architect.md",
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
	agents := []Agent{Gemini, OpenCode}
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
