package main

import (
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Wather17/Agents/agent-init/internal/templates"
)

func TestRunUpgradeContinuesWhenCLISelfUpdateFails(t *testing.T) {
	target := t.TempDir()
	prompt, err := templates.Read("files/GEMINI.md")
	if err != nil {
		t.Fatalf("failed to read prompt template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "GEMINI.md"), prompt, 0o644); err != nil {
		t.Fatalf("failed to create existing prompt: %v", err)
	}

	originalUpdateCLI := updateCLI
	originalCommandLine := flag.CommandLine
	originalArgs := os.Args
	t.Cleanup(func() {
		updateCLI = originalUpdateCLI
		flag.CommandLine = originalCommandLine
		os.Args = originalArgs
	})

	updateCLI = func() error {
		return errors.New("network unavailable")
	}
	flag.CommandLine = flag.NewFlagSet("agent-init", flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = []string{"agent-init", "upgrade", "--path", target, "--no-commit"}

	if err := runUpgradeCommand(); err != nil {
		t.Fatalf("upgrade should continue when self-update fails: %v", err)
	}

	for _, path := range []string{
		filepath.Join(target, ".agents", "skills", "audit-issues", "SKILL.md"),
		filepath.Join(target, ".agents", "skills", "refine-issues", "SKILL.md"),
		filepath.Join(target, ".agents", "skills", "autonomous-batch", "SKILL.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("upgrade should install %s using current embedded templates: %v", path, err)
		}
	}
}

func TestRunUpgradeRelaunchesAfterSuccessfulSelfUpdate(t *testing.T) {
	target := t.TempDir()

	originalUpdateCLI := updateCLI
	originalRelaunchUpgrade := relaunchUpgrade
	originalCommandLine := flag.CommandLine
	originalArgs := os.Args
	t.Cleanup(func() {
		updateCLI = originalUpdateCLI
		relaunchUpgrade = originalRelaunchUpgrade
		flag.CommandLine = originalCommandLine
		os.Args = originalArgs
	})

	updated := false
	relaunched := false
	updateCLI = func() error {
		updated = true
		return nil
	}
	relaunchUpgrade = func() error {
		relaunched = true
		return nil
	}
	flag.CommandLine = flag.NewFlagSet("agent-init", flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = []string{"agent-init", "upgrade", "--path", target, "--no-commit"}
	t.Setenv(skipSelfUpdateEnv, "")

	if err := runUpgradeCommand(); err != nil {
		t.Fatalf("upgrade should relaunch after self-update: %v", err)
	}
	if !updated {
		t.Error("upgrade should attempt to update the CLI")
	}
	if !relaunched {
		t.Error("upgrade should relaunch the updated CLI")
	}
}
