// Package installer handles the installation of agent configuration files.
package installer

import (
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

// Result describes an installed file.
type Result struct {
	Path    string
	Ignored bool
}

// Install writes the embedded template files to the target directory.
func Install(opts Options) ([]Result, error) {
	files, err := templates.FilesFor(opts.Agent)
	if err != nil {
		return nil, err
	}

	installed := make([]Result, 0, len(files))

	for _, file := range files {
		content, err := templates.Read(file.SourcePath)
		if err != nil {
			return installed, fmt.Errorf("reading embedded file %q: %w", file.SourcePath, err)
		}

		target := filepath.Join(opts.TargetPath, file.TargetPath)

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return installed, fmt.Errorf("creating directory for %q: %w", target, err)
		}

		if _, err := os.Stat(target); err == nil && !opts.Force {
			return installed, fmt.Errorf("file already exists: %s (use --force to overwrite)", target)
		} else if err != nil && !os.IsNotExist(err) {
			return installed, fmt.Errorf("checking file %q: %w", target, err)
		}

		mode := os.FileMode(0o644)
		if file.Executable {
			mode = 0o755
		}

		if err := os.WriteFile(target, content, mode); err != nil {
			return installed, fmt.Errorf("writing file %q: %w", target, err)
		}

		installed = append(installed, Result{Path: target, Ignored: file.Ignored})
	}

	if err := ensureGitignore(opts.TargetPath, opts.Agent); err != nil {
		return installed, fmt.Errorf("updating .gitignore: %w", err)
	}

	return installed, nil
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
