// Package templates provides the embedded files used by agent-init.
package templates

import (
	"embed"
	"fmt"
)

//go:embed files/*
var fs embed.FS

// Agent represents a supported AI agent configuration.
type Agent string

const (
	// Gemini is the default agent template.
	Gemini Agent = "gemini"
	// OpenCode is the OpenCode agent template.
	OpenCode Agent = "opencode"
)

// File represents a file to be installed for a given agent.
type File struct {
	SourcePath string
	TargetPath string
	Executable bool
	// Ignored indicates whether the file should be added to .gitignore and
	// excluded from the automatic commit.
	Ignored bool
	// Prompt identifies the primary prompt that signals an agent is installed.
	Prompt bool
}

// FilesFor returns the list of files to install for the given agent.
func FilesFor(agent Agent) ([]File, error) {
	switch agent {
	case Gemini:
		return []File{
			{SourcePath: "files/GEMINI.md", TargetPath: "GEMINI.md", Executable: false, Ignored: true, Prompt: true},
			{SourcePath: "files/skills/audit-issues/SKILL.md", TargetPath: ".agents/skills/audit-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/refine-issues/SKILL.md", TargetPath: ".agents/skills/refine-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/autonomous-batch/SKILL.md", TargetPath: ".agents/skills/autonomous-batch/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/sync-issues.sh", TargetPath: "scripts/sync-issues.sh", Executable: true, Ignored: false},
		}, nil
	case OpenCode:
		return []File{
			{SourcePath: "files/AGENTS.md", TargetPath: "AGENTS.md", Executable: false, Ignored: true, Prompt: true},
			{SourcePath: "files/skills/audit-issues/SKILL.md", TargetPath: ".opencode/skill/audit-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/refine-issues/SKILL.md", TargetPath: ".opencode/skill/refine-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/autonomous-batch/SKILL.md", TargetPath: ".opencode/skill/autonomous-batch/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/sync-issues.sh", TargetPath: "scripts/sync-issues.sh", Executable: true, Ignored: false},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported agent: %q", agent)
	}
}

// Read returns the content of an embedded file.
func Read(sourcePath string) ([]byte, error) {
	return fs.ReadFile(sourcePath)
}

// IgnoredEntries returns the gitignore entries that should be added for the
// installed files.
func IgnoredEntries(agent Agent) ([]string, error) {
	switch agent {
	case Gemini:
		return []string{
			"# AI agent configuration files",
			"GEMINI.md",
			".agents/skills/audit-issues/SKILL.md",
			".agents/skills/refine-issues/SKILL.md",
			".agents/skills/autonomous-batch/SKILL.md",
			"",
			"# Synced GitHub issues",
			"issues/",
		}, nil
	case OpenCode:
		return []string{
			"# AI agent configuration files",
			"AGENTS.md",
			".opencode/skill/audit-issues/SKILL.md",
			".opencode/skill/refine-issues/SKILL.md",
			".opencode/skill/autonomous-batch/SKILL.md",
			"",
			"# Synced GitHub issues",
			"issues/",
		}, nil
	default:
		return nil, fmt.Errorf("unsupported agent: %q", agent)
	}
}
