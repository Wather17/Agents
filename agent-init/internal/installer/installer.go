// Package installer handles the installation of agent configuration files.
package installer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wather17/Agents/agent-init/internal/templates"
)

// Options controls how files are installed.
type Options struct {
	TargetPath string
	Agent      templates.Agent
	Force      bool
}

// Result describes a file handled during installation.
type Result struct {
	Path    string
	Ignored bool
}

// Install writes the embedded template files to the target directory.
// It returns the files that were installed, the files that were skipped
// because they already exist with identical content, and any error.
func Install(opts Options) (installed []Result, skipped []Result, err error) {
	files, err := templates.FilesFor(opts.Agent)
	if err != nil {
		return nil, nil, err
	}

	installed = make([]Result, 0, len(files))
	skipped = make([]Result, 0, len(files))

	for _, file := range files {
		content, err := templates.Read(file.SourcePath)
		if err != nil {
			return installed, skipped, fmt.Errorf("reading embedded file %q: %w", file.SourcePath, err)
		}

		target := filepath.Join(opts.TargetPath, file.TargetPath)

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return installed, skipped, fmt.Errorf("creating directory for %q: %w", target, err)
		}

		exists, identical, err := compareFile(target, content)
		if err != nil {
			return installed, skipped, fmt.Errorf("checking file %q: %w", target, err)
		}

		if exists {
			if identical {
				skipped = append(skipped, Result{Path: target, Ignored: file.Ignored})
				continue
			}
			if !opts.Force {
				return installed, skipped, fmt.Errorf("file already exists with different content: %s (use --force to overwrite)", target)
			}
		}

		mode := os.FileMode(0o644)
		if file.Executable {
			mode = 0o755
		}

		if err := os.WriteFile(target, content, mode); err != nil {
			return installed, skipped, fmt.Errorf("writing file %q: %w", target, err)
		}

		installed = append(installed, Result{Path: target, Ignored: file.Ignored})
	}

	if err := ensureGitignore(opts.TargetPath, opts.Agent); err != nil {
		return installed, skipped, fmt.Errorf("updating .gitignore: %w", err)
	}

	return installed, skipped, nil
}

// Upgrade updates already-installed agent prompt files and completes missing
// auxiliary files in the target directory. It does not install a new agent
// prompt or overwrite shared files such as scripts/sync-issues.sh. An agent is
// considered installed when its primary prompt file exists. It returns the
// files that were installed or updated and the files that were skipped because
// they already match the latest template.
func Upgrade(targetPath string) (installed []Result, skipped []Result, err error) {
	installed = make([]Result, 0)
	skipped = make([]Result, 0)

	agents := []templates.Agent{templates.Gemini, templates.OpenCode}

	seen := make(map[string]bool)
	files := make([]templates.File, 0)

	for _, agent := range agents {
		agentFiles, err := templates.FilesFor(agent)
		if err != nil {
			return installed, skipped, err
		}

		installedAgent := false
		for _, file := range agentFiles {
			if !file.Prompt {
				continue
			}

			exists, err := fileExists(filepath.Join(targetPath, file.TargetPath))
			if err != nil {
				return installed, skipped, fmt.Errorf("checking agent prompt %q: %w", file.TargetPath, err)
			}
			installedAgent = exists
			break
		}
		if !installedAgent {
			continue
		}

		for _, file := range agentFiles {
			// Only manage ignored agent files, never shared scripts, and never
			// process the same shared file twice across agents.
			if !file.Ignored || seen[file.TargetPath] {
				continue
			}
			seen[file.TargetPath] = true
			files = append(files, file)
		}

		if err := ensureGitignore(targetPath, agent); err != nil {
			return installed, skipped, fmt.Errorf("updating .gitignore: %w", err)
		}
	}

	for _, file := range files {
		target := filepath.Join(targetPath, file.TargetPath)
		exists, err := fileExists(target)
		if err != nil {
			return installed, skipped, fmt.Errorf("checking file %q: %w", target, err)
		}
		if !exists {
			if file.Prompt {
				// The prompt existed during agent detection. Do not create it if
				// it was removed while the upgrade was running.
				continue
			}
		}

		content, err := templates.Read(file.SourcePath)
		if err != nil {
			return installed, skipped, fmt.Errorf("reading embedded file %q: %w", file.SourcePath, err)
		}

		exists, identical, err := compareFile(target, content)
		if err != nil {
			return installed, skipped, fmt.Errorf("comparing file %q: %w", target, err)
		}
		if exists && identical {
			skipped = append(skipped, Result{Path: target, Ignored: file.Ignored})
			continue
		}

		mode := os.FileMode(0o644)
		if file.Executable {
			mode = 0o755
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return installed, skipped, fmt.Errorf("creating directory for %q: %w", target, err)
		}
		if err := os.WriteFile(target, content, mode); err != nil {
			return installed, skipped, fmt.Errorf("writing file %q: %w", target, err)
		}
		installed = append(installed, Result{Path: target, Ignored: file.Ignored})
	}

	return installed, skipped, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// compareFile checks whether the file at path exists and whether its content
// is identical to the provided content. It returns (false, false, nil) when the
// file does not exist.
func compareFile(path string, content []byte) (exists bool, identical bool, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return true, bytes.Equal(data, content), nil
}

func ensureGitignore(targetPath string, agent templates.Agent) error {
	entries, err := templates.IgnoredEntries(agent)
	if err != nil {
		return err
	}

	gitignorePath := filepath.Join(targetPath, ".gitignore")

	existing := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		existing = string(data)
	} else if !os.IsNotExist(err) {
		return err
	}

	missing := filterMissing(entries, existing)
	if len(missing) == 0 {
		return nil
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if existing != "" && !endsWithNewline(existing) {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}

	for _, entry := range missing {
		if _, err := f.WriteString(entry + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func filterMissing(entries []string, existing string) []string {
	var missing []string
	for _, entry := range entries {
		if entry == "" {
			continue
		}
		if !containsLine(existing, entry) {
			missing = append(missing, entry)
		}
	}
	return missing
}

func containsLine(text, line string) bool {
	idx := 0
	for {
		i := indexOf(text[idx:], line)
		if i == -1 {
			return false
		}
		abs := idx + i
		before := abs == 0 || text[abs-1] == '\n'
		after := abs+len(line) == len(text) || text[abs+len(line)] == '\n'
		if before && after {
			return true
		}
		idx = abs + 1
	}
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func endsWithNewline(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\n'
}
