package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const DefaultOutputDir = "output"

type Config struct {
	OutputDir string   `json:"output_dir"`
	Domains   []string `json:"domains"`
}

func Load(path string) (Config, error) {
	if strings.TrimSpace(path) == "" {
		return Config{}, fmt.Errorf("config path must not be empty")
	}
	if !strings.HasSuffix(path, ".json") {
		return Config{}, fmt.Errorf("config file %q must have .json extension", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("load config %q: %w", path, err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	if err := decoder.Decode(&struct{}{}); err == nil {
		return Config{}, fmt.Errorf("parse config %q: unexpected extra JSON content", path)
	}

	if err := cfg.normalize(); err != nil {
		return Config{}, fmt.Errorf("validate config %q: %w", path, err)
	}

	return cfg, nil
}

func Resolve(flagOutputDir string, outputDirFlagSet bool, cfg Config) Config {
	resolved := cfg
	if resolved.OutputDir == "" {
		resolved.OutputDir = DefaultOutputDir
	}

	if outputDirFlagSet {
		resolved.OutputDir = strings.TrimSpace(flagOutputDir)
		if resolved.OutputDir == "" {
			resolved.OutputDir = DefaultOutputDir
		}
	}

	return resolved
}

func (cfg *Config) normalize() error {
	if cfg.OutputDir = strings.TrimSpace(cfg.OutputDir); cfg.OutputDir == "" {
		cfg.OutputDir = DefaultOutputDir
	}

	if len(cfg.Domains) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(cfg.Domains))
	for index, domain := range cfg.Domains {
		trimmed := strings.TrimSpace(domain)
		if trimmed == "" {
			return fmt.Errorf("domains[%d] must not be empty", index)
		}
		normalized = append(normalized, trimmed)
	}

	cfg.Domains = normalized
	return nil
}
