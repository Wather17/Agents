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

var (
	updateCLI       = selfUpdate
	relaunchUpgrade = relaunchUpgradeCommand
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "update":
			if err := selfUpdate(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			return
		case "upgrade":
			if err := runUpgradeCommand(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			return
		case "version":
			printVersion()
			return
		}
	}

	var (
		agentFlag = flag.String("agent", "antigravity", "Agent template to install (antigravity, codex, opencode, or claude)")
		forceFlag = flag.Bool("force", false, "Overwrite existing files")
		pathFlag  = flag.String("path", ".", "Target repository path")
		noCommit  = flag.Bool("no-commit", false, "Skip creating a git commit")
	)
	flag.Parse()

	if err := runInstall(*agentFlag, *pathFlag, *forceFlag, *noCommit); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runUpgradeCommand() error {
	var pathFlag = flag.String("path", ".", "Target repository path")
	var noCommit = flag.Bool("no-commit", false, "Skip creating a git commit")
	flag.CommandLine.Parse(os.Args[2:])

	absPath, err := filepath.Abs(*pathFlag)
	if err != nil {
		return fmt.Errorf("resolving target path: %w", err)
	}

	if os.Getenv(skipSelfUpdateEnv) != "1" {
		if err := updateCLI(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v; continuing with current embedded templates\n", err)
		} else {
			return relaunchUpgrade()
		}
	}

	installed, skipped, err := installer.Upgrade(absPath)
	if err != nil {
		return err
	}

	printResults(installed, skipped)

	if *noCommit {
		fmt.Println("Skipping git commit (--no-commit).")
		return nil
	}

	return commitIfNeeded(absPath, installed)
}

func runInstall(agentName, targetPath string, force, skipCommit bool) error {
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("resolving target path: %w", err)
	}

	agent, err := templates.ParseAgent(agentName)
	if err != nil {
		return err
	}

	installed, skipped, err := installer.Install(installer.Options{
		TargetPath: absPath,
		Agent:      agent,
		Force:      force,
	})
	if err != nil {
		return err
	}

	printResults(installed, skipped)

	if skipCommit {
		fmt.Println("Skipping git commit (--no-commit).")
		return nil
	}

	return commitIfNeeded(absPath, installed)
}

func printResults(installed, skipped []installer.Result) {
	if len(installed) > 0 {
		fmt.Println("Installed files:")
	}
	for _, r := range installed {
		status := ""
		if r.Ignored {
			status = " (ignored by git)"
		}
		fmt.Printf("  - %s%s\n", r.Path, status)
	}

	if len(skipped) > 0 {
		fmt.Println("Skipped (already exists and identical):")
		for _, r := range skipped {
			fmt.Printf("  - %s\n", r.Path)
		}
	}
}

func commitIfNeeded(absPath string, installed []installer.Result) error {
	var commitable []string
	for _, r := range installed {
		if !r.Ignored {
			commitable = append(commitable, r.Path)
		}
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

	msg := "chore(agent-init): update agent configuration files"
	if len(commitable) > 0 {
		msg = fmt.Sprintf("chore(agent-init): update %s", filepath.Base(commitable[0]))
	}
	if err := git.Commit(absPath, relFiles, msg); err != nil {
		return err
	}

	fmt.Println("Committed changes.")
	return nil
}
