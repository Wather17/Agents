// agent-init installs AI agent configuration files into a project repository.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wather17/Agents/agent-init/internal/git"
	"github.com/Wather17/Agents/agent-init/internal/installer"
	"github.com/Wather17/Agents/agent-init/internal/templates"
)

func main() {
	var (
		agentFlag = flag.String("agent", "gemini", "Agent template to install (e.g., gemini)")
		forceFlag = flag.Bool("force", false, "Overwrite existing files")
		pathFlag  = flag.String("path", ".", "Target repository path")
		noCommit  = flag.Bool("no-commit", false, "Skip creating a git commit")
	)
	flag.Parse()

	if err := run(*agentFlag, *pathFlag, *forceFlag, *noCommit); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(agentName, targetPath string, force, skipCommit bool) error {
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("resolving target path: %w", err)
	}

	agent := templates.Agent(agentName)

	installed, err := installer.Install(installer.Options{
		TargetPath: absPath,
		Agent:      agent,
		Force:      force,
	})
	if err != nil {
		return err
	}

	fmt.Println("Installed files:")
	var commitable []string
	for _, r := range installed {
		status := ""
		if r.Ignored {
			status = " (ignored by git)"
		} else {
			commitable = append(commitable, r.Path)
		}
		fmt.Printf("  - %s%s\n", r.Path, status)
	}

	if skipCommit {
		fmt.Println("Skipping git commit (--no-commit).")
		return nil
	}

	if !git.IsRepository(absPath) {
		fmt.Println("Target path is not a git repository; skipping commit.")
		return nil
	}

	gitignorePath := filepath.Join(absPath, ".gitignore")
	filesToCommit := append([]string{gitignorePath}, commitable...)

	relFiles := make([]string, len(filesToCommit))
	for i, f := range filesToCommit {
		rel, err := filepath.Rel(absPath, f)
		if err != nil {
			return err
		}
		relFiles[i] = rel
	}

	hasChanges, err := git.HasChanges(absPath)
	if err != nil {
		return err
	}
	if !hasChanges {
		fmt.Println("No changes to commit.")
		return nil
	}

	msg := fmt.Sprintf("chore(agent-init): add %s agent configuration and issue sync script", agent)
	if err := git.Commit(absPath, relFiles, msg); err != nil {
		return err
	}

	fmt.Println("Committed changes.")
	return nil
}
