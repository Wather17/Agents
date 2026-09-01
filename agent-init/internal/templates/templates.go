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
	// Antigravity is the default agent template.
	Antigravity Agent = "antigravity"
	// Gemini is a backwards-compatible name for Antigravity.
	Gemini Agent = Antigravity
	// OpenCode is the OpenCode agent template.
	OpenCode Agent = "opencode"
	// Codex is the OpenAI Codex agent template.
	Codex Agent = "codex"
	// Claude is the Claude Code agent template.
	Claude Agent = "claude"
)

// ParseAgent normalizes a CLI agent name to its canonical identifier.
func ParseAgent(name string) (Agent, error) {
	switch Agent(name) {
	case "gemini", Antigravity:
		return Antigravity, nil
	case OpenCode:
		return OpenCode, nil
	case Codex:
		return Codex, nil
	case Claude:
		return Claude, nil
	default:
		return "", fmt.Errorf("unsupported agent: %q", name)
	}
}

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
	case Antigravity:
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
	case Codex:
		return []File{
			{SourcePath: "files/AGENTS.md", TargetPath: "AGENTS.md", Executable: false, Ignored: true, Prompt: true},
			{SourcePath: "files/skills/audit-issues/SKILL.md", TargetPath: ".agents/skills/audit-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/refine-issues/SKILL.md", TargetPath: ".agents/skills/refine-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/autonomous-batch/SKILL.md", TargetPath: ".agents/skills/autonomous-batch/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/sync-issues.sh", TargetPath: "scripts/sync-issues.sh", Executable: true, Ignored: false},
		}, nil
	case Claude:
		return []File{
			{SourcePath: "files/AGENTS.md", TargetPath: "CLAUDE.md", Executable: false, Ignored: true, Prompt: true},
			{SourcePath: "files/skills/audit-issues/SKILL.md", TargetPath: ".claude/skills/audit-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/refine-issues/SKILL.md", TargetPath: ".claude/skills/refine-issues/SKILL.md", Executable: false, Ignored: true},
			{SourcePath: "files/skills/autonomous-batch/SKILL.md", TargetPath: ".claude/skills/autonomous-batch/SKILL.md", Executable: false, Ignored: true},
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
	case Antigravity:
		return []string{
			"# AI agent configuration files",
			".agent-init.json",
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
			".agent-init.json",
			"AGENTS.md",
			".opencode/skill/audit-issues/SKILL.md",
			".opencode/skill/refine-issues/SKILL.md",
			".opencode/skill/autonomous-batch/SKILL.md",
			"",
			"# Synced GitHub issues",
			"issues/",
		}, nil
	case Codex:
		return []string{
			"# AI agent configuration files",
			".agent-init.json",
			"AGENTS.md",
			".agents/skills/audit-issues/SKILL.md",
			".agents/skills/refine-issues/SKILL.md",
			".agents/skills/autonomous-batch/SKILL.md",
			"",
			"# Synced GitHub issues",
			"issues/",
		}, nil
	case Claude:
		return []string{
			"# AI agent configuration files",
			".agent-init.json",
			"CLAUDE.md",
			".claude/skills/audit-issues/SKILL.md",
			".claude/skills/refine-issues/SKILL.md",
			".claude/skills/autonomous-batch/SKILL.md",
			"",
			"# Synced GitHub issues",
			"issues/",
		}, nil
	default:
		return nil, fmt.Errorf("unsupported agent: %q", agent)
	}
}
