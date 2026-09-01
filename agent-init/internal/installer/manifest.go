package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Wather17/Agents/agent-init/internal/templates"
)

const (
	manifestFilename = ".agent-init.json"
	manifestVersion  = 1
)

type manifest struct {
	Version int               `json:"version"`
	Agents  []templates.Agent `json:"agents"`
}

func readManifest(targetPath string) (manifest, bool, error) {
	path := filepath.Join(targetPath, manifestFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return manifest{Version: manifestVersion}, false, nil
		}
		return manifest{}, false, fmt.Errorf("reading agent manifest: %w", err)
	}

	var value manifest
	if err := json.Unmarshal(data, &value); err != nil {
		return manifest{}, true, fmt.Errorf("invalid %s: %w", manifestFilename, err)
	}
	if err := validateManifest(value); err != nil {
		return manifest{}, true, err
	}
	normalizeManifest(&value)
	return value, true, nil
}

func validateManifest(value manifest) error {
	if value.Version != manifestVersion {
		return fmt.Errorf("unsupported %s version %d (expected %d)", manifestFilename, value.Version, manifestVersion)
	}
	for _, agent := range value.Agents {
		parsed, err := templates.ParseAgent(string(agent))
		if err != nil {
			return fmt.Errorf("invalid %s agent %q", manifestFilename, agent)
		}
		if parsed != agent {
			return fmt.Errorf("non-canonical %s agent %q", manifestFilename, agent)
		}
	}
	return nil
}

func normalizeManifest(value *manifest) {
	unique := make(map[templates.Agent]bool, len(value.Agents))
	agents := make([]templates.Agent, 0, len(value.Agents))
	for _, agent := range value.Agents {
		if unique[agent] {
			continue
		}
		unique[agent] = true
		agents = append(agents, agent)
	}
	sort.Slice(agents, func(i, j int) bool { return agents[i] < agents[j] })
	value.Agents = agents
}

func addManifestAgent(value *manifest, agent templates.Agent) bool {
	for _, existing := range value.Agents {
		if existing == agent {
			return false
		}
	}
	value.Agents = append(value.Agents, agent)
	normalizeManifest(value)
	return true
}

func writeManifest(targetPath string, value manifest) error {
	normalizeManifest(&value)
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding agent manifest: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(filepath.Join(targetPath, manifestFilename), data, 0o644); err != nil {
		return fmt.Errorf("writing agent manifest: %w", err)
	}
	return nil
}

func inferLegacyManifest(targetPath string) (manifest, error) {
	value := manifest{Version: manifestVersion}
	legacyPrompts := []struct {
		path  string
		agent templates.Agent
	}{
		{path: "GEMINI.md", agent: templates.Antigravity},
		{path: "AGENTS.md", agent: templates.OpenCode},
	}

	for _, candidate := range legacyPrompts {
		exists, err := fileExists(filepath.Join(targetPath, candidate.path))
		if err != nil {
			return manifest{}, fmt.Errorf("checking legacy agent prompt %q: %w", candidate.path, err)
		}
		if exists {
			addManifestAgent(&value, candidate.agent)
		}
	}
	return value, nil
}
