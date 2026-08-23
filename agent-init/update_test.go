package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseAssetNameFor(t *testing.T) {
	tests := []struct {
		name   string
		goos   string
		goarch string
		want   string
	}{
		{name: "linux amd64", goos: "linux", goarch: "amd64", want: "agent-init-linux-amd64"},
		{name: "linux arm64", goos: "linux", goarch: "arm64", want: "agent-init-linux-arm64"},
		{name: "windows amd64", goos: "windows", goarch: "amd64", want: "agent-init-windows-amd64.exe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := releaseAssetNameFor(tt.goos, tt.goarch)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("releaseAssetNameFor(%q, %q) = %q, want %q", tt.goos, tt.goarch, got, tt.want)
			}
		})
	}
}

func TestReleaseAssetNameForRejectsUnsupportedPlatform(t *testing.T) {
	if _, err := releaseAssetNameFor("darwin", "arm64"); err == nil {
		t.Fatal("expected unsupported platform error")
	}
}

func TestVerifyChecksum(t *testing.T) {
	target := t.TempDir()
	assetName := "agent-init-linux-amd64"
	assetPath := filepath.Join(target, assetName)
	checksumPath := filepath.Join(target, "checksums.txt")
	content := []byte("agent-init release binary")
	if err := os.WriteFile(assetPath, content, 0o644); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}
	digest := sha256.Sum256(content)
	checksums := fmt.Sprintf("%x  %s\n", digest, assetName)
	if err := os.WriteFile(checksumPath, []byte(checksums), 0o644); err != nil {
		t.Fatalf("failed to create checksums: %v", err)
	}

	if err := verifyChecksum(assetPath, checksumPath, assetName); err != nil {
		t.Fatalf("valid checksum should pass: %v", err)
	}

	if err := os.WriteFile(assetPath, []byte("tampered binary"), 0o644); err != nil {
		t.Fatalf("failed to tamper asset: %v", err)
	}
	if err := verifyChecksum(assetPath, checksumPath, assetName); err == nil {
		t.Fatal("tampered asset should fail checksum verification")
	}
}

func TestVerifyChecksumRequiresAssetEntry(t *testing.T) {
	target := t.TempDir()
	assetPath := filepath.Join(target, "agent-init-linux-amd64")
	checksumPath := filepath.Join(target, "checksums.txt")
	if err := os.WriteFile(assetPath, []byte("binary"), 0o644); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}
	if err := os.WriteFile(checksumPath, []byte("deadbeef  another-file\n"), 0o644); err != nil {
		t.Fatalf("failed to create checksums: %v", err)
	}

	if err := verifyChecksum(assetPath, checksumPath, filepath.Base(assetPath)); err == nil {
		t.Fatal("missing asset checksum entry should fail")
	}
}

func TestReplaceExecutable(t *testing.T) {
	target := t.TempDir()
	sourcePath := filepath.Join(target, "downloaded")
	targetPath := filepath.Join(target, "agent-init")
	if err := os.WriteFile(sourcePath, []byte("new binary"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("old binary"), 0o755); err != nil {
		t.Fatalf("failed to create target: %v", err)
	}

	if err := replaceExecutable(sourcePath, targetPath); err != nil {
		t.Fatalf("replaceExecutable should succeed: %v", err)
	}
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read replaced executable: %v", err)
	}
	if string(content) != "new binary" {
		t.Errorf("replaced executable content = %q, want %q", content, "new binary")
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("failed to stat replaced executable: %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Error("replaced executable should be executable")
	}
}

func TestEnvironmentWithValueReplacesExistingValue(t *testing.T) {
	environment := []string{"PATH=/bin", "AGENT_INIT_SKIP_SELF_UPDATE=0", "HOME=/tmp"}
	updated := environmentWithValue(environment, "AGENT_INIT_SKIP_SELF_UPDATE", "1")
	joined := strings.Join(updated, "\n")
	if strings.Contains(joined, "AGENT_INIT_SKIP_SELF_UPDATE=0") {
		t.Error("old environment value should be removed")
	}
	if !strings.Contains(joined, "AGENT_INIT_SKIP_SELF_UPDATE=1") {
		t.Error("new environment value should be present")
	}
}
