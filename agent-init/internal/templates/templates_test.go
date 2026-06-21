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

func TestReadEmbeddedFiles(t *testing.T) {
	files, _ := FilesFor(Gemini)
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
