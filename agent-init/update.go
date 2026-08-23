package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	packagePath       = "github.com/Wather17/Agents/agent-init"
	releaseRepository = "Wather17/Agents"
	checksumAsset     = "checksums.txt"
	skipSelfUpdateEnv = "AGENT_INIT_SKIP_SELF_UPDATE"
)

func selfUpdate() error {
	asset, err := releaseAssetName()
	if err != nil {
		return err
	}

	updateDir, err := os.MkdirTemp("", "agent-init-update-")
	if err != nil {
		return fmt.Errorf("creating update directory: %w", err)
	}
	defer os.RemoveAll(updateDir)

	fmt.Printf("Updating %s from the latest release...\n", packagePath)
	if err := downloadReleaseAssets(updateDir, asset); err != nil {
		return err
	}

	assetPath := filepath.Join(updateDir, asset)
	checksumPath := filepath.Join(updateDir, checksumAsset)
	if err := verifyChecksum(assetPath, checksumPath, asset); err != nil {
		return err
	}

	target, err := currentExecutablePath()
	if err != nil {
		return err
	}
	if err := replaceExecutable(assetPath, target); err != nil {
		return err
	}

	fmt.Printf("Update complete: %s\n", target)
	return nil
}

func downloadReleaseAssets(destination, asset string) error {
	cmd := exec.Command(
		"gh", "release", "download",
		"--repo", releaseRepository,
		"--pattern", asset,
		"--pattern", checksumAsset,
		"--dir", destination,
		"--clobber",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("downloading latest release with gh: %w", err)
	}
	return nil
}

func releaseAssetName() (string, error) {
	return releaseAssetNameFor(runtime.GOOS, runtime.GOARCH)
}

func releaseAssetNameFor(goos, goarch string) (string, error) {
	switch {
	case goos == "linux" && goarch == "amd64":
		return "agent-init-linux-amd64", nil
	case goos == "linux" && goarch == "arm64":
		return "agent-init-linux-arm64", nil
	case goos == "windows" && goarch == "amd64":
		return "agent-init-windows-amd64.exe", nil
	default:
		return "", fmt.Errorf("unsupported platform for self-update: %s/%s", goos, goarch)
	}
}

func verifyChecksum(assetPath, checksumPath, assetName string) error {
	checksums, err := os.Open(checksumPath)
	if err != nil {
		return fmt.Errorf("opening checksums: %w", err)
	}
	defer checksums.Close()

	var expected string
	scanner := bufio.NewScanner(checksums)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && strings.TrimPrefix(fields[1], "*") == assetName {
			expected = fields[0]
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading checksums: %w", err)
	}
	if expected == "" {
		return fmt.Errorf("checksum for %q not found", assetName)
	}

	asset, err := os.Open(assetPath)
	if err != nil {
		return fmt.Errorf("opening downloaded asset: %w", err)
	}
	defer asset.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, asset); err != nil {
		return fmt.Errorf("calculating checksum for %q: %w", assetName, err)
	}
	actual := fmt.Sprintf("%x", hash.Sum(nil))
	if !strings.EqualFold(expected, actual) {
		return fmt.Errorf("checksum mismatch for %q: expected %s, got %s", assetName, expected, actual)
	}
	return nil
}

func currentExecutablePath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("finding current executable: %w", err)
	}

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolving current executable: %w", err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolving executable path: %w", err)
	}
	return path, nil
}

func replaceExecutable(source, target string) error {
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening update asset: %w", err)
	}
	defer input.Close()

	temporary, err := os.CreateTemp(filepath.Dir(target), ".agent-init-update-")
	if err != nil {
		return fmt.Errorf("creating executable temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := io.Copy(temporary, input); err != nil {
		temporary.Close()
		return fmt.Errorf("copying update asset: %w", err)
	}
	if err := temporary.Chmod(0o755); err != nil {
		temporary.Close()
		return fmt.Errorf("setting executable permissions: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("syncing updated executable: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("closing updated executable: %w", err)
	}

	if err := os.Rename(temporaryPath, target); err != nil {
		return fmt.Errorf("replacing current executable: %w", err)
	}
	return nil
}

func relaunchUpgradeCommand() error {
	executable, err := currentExecutablePath()
	if err != nil {
		return err
	}

	args := append([]string(nil), os.Args[1:]...)
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = environmentWithValue(os.Environ(), skipSelfUpdateEnv, "1")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running updated agent-init: %w", err)
	}
	return nil
}

func environmentWithValue(environment []string, key, value string) []string {
	prefix := key + "="
	updated := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		if strings.HasPrefix(entry, prefix) {
			continue
		}
		updated = append(updated, entry)
	}
	return append(updated, prefix+value)
}
