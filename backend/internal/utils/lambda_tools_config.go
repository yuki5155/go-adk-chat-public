package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ToolParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type ToolSchema struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  []ToolParameter `json:"parameters"`
}

type LambdaToolEntry struct {
	Name      string     `yaml:"name"`
	Route     string     `yaml:"route"`
	LocalPort int        `yaml:"local_port"`
	Schema    *ToolSchema
}

type lambdaToolsManifest struct {
	Tools []LambdaToolEntry `yaml:"tools"`
}

// LoadLambdaTools reads tools.yaml and each tool's schema.json.
func LoadLambdaTools(configPath string) ([]LambdaToolEntry, error) {
	cleanConfig := filepath.Clean(configPath)
	data, err := os.ReadFile(cleanConfig) // #nosec G304 -- path comes from trusted config env var, not user input
	if err != nil {
		return nil, fmt.Errorf("lambda tools config: read %s: %w", cleanConfig, err)
	}
	var manifest lambdaToolsManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("lambda tools config: parse: %w", err)
	}

	dir := filepath.Dir(cleanConfig)
	for i, entry := range manifest.Tools {
		schemaPath := filepath.Join(dir, filepath.Clean(entry.Name), "schema.json")
		schemaData, err := os.ReadFile(schemaPath) // #nosec G304 -- tool name comes from trusted yaml, not user input
		if err != nil {
			return nil, fmt.Errorf("lambda tools config: read schema for %s: %w", entry.Name, err)
		}
		var schema ToolSchema
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			return nil, fmt.Errorf("lambda tools config: parse schema for %s: %w", entry.Name, err)
		}
		manifest.Tools[i].Schema = &schema
	}

	return manifest.Tools, nil
}

// BuildLambdaToolURL returns the endpoint URL for a tool.
// In development it uses the local RIE port; in production it uses the API Gateway base URL + route.
func BuildLambdaToolURL(entry LambdaToolEntry, baseURL string, isDev bool) string {
	if isDev {
		return fmt.Sprintf(
			"http://host.docker.internal:%d/2015-03-31/functions/function/invocations",
			entry.LocalPort,
		)
	}
	return baseURL + entry.Route
}
