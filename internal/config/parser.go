package config

import (
	"os"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Step struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}

type Pipeline struct {
	Steps []Step `yaml:"steps"`
}

func LoadConfig(repoPath string) ([]string, error) {
	path := filepath.Join(repoPath, ".cicd.yaml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("⚠️  No .cicd.yaml found at: %s", path)
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("❌ Failed to read .cicd.yaml: %v", err)
		return nil, err
	}

	var pipeline Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		log.Printf("❌ Failed to parse YAML: %v", err)
		return nil, err
	}

	var cmds []string
	for _, step := range pipeline.Steps {
		cmds = append(cmds, step.Run)
	}

  log.Printf("✅ Loaded %d command(s) from .cicd.yaml", len(cmds))
	return cmds, nil
}

