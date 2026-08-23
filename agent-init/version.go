package main

import "fmt"

var (
	// Version, Commit, and Date are populated by the release workflow via ldflags.
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func printVersion() {
	fmt.Printf("agent-init %s\n", Version)
	if Commit != "none" {
		fmt.Printf("Commit: %s\n", Commit)
	}
	if Date != "unknown" {
		fmt.Printf("Build date: %s\n", Date)
	}
}
