package templates

import "testing"

func TestFilesForGemini(t *testing.T) {
	files, err := FilesFor(Gemini)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	if files[0].TargetPath != "GEMINI.md" || !files[0].Ignored {
		t.Errorf("GEMINI.md should be ignored, got %+v", files[0])
	}

	if files[1].TargetPath != "scripts/sync-issues.sh" || files[1].Ignored {
		t.Errorf("sync-issues.sh should not be ignored, got %+v", files[1])
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

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	if files[0].TargetPath != "AGENTS.md" || !files[0].Ignored {
		t.Errorf("AGENTS.md should be ignored, got %+v", files[0])
	}

	if files[1].TargetPath != "scripts/sync-issues.sh" || files[1].Ignored {
		t.Errorf("sync-issues.sh should not be ignored, got %+v", files[1])
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
