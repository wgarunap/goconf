package goconf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseYaml reads a YAML configuration file and unmarshals it into the provided struct.
// The struct fields should use yaml tags to map to YAML keys.
//
// Parameters:
//   - config (interface{}): Pointer to the struct to be populated from the YAML file.
//     The struct should have yaml struct tags for field mapping.
//   - filePath (string): Path to the YAML configuration file.
//
// Returns:
//   - error: Returns error if file reading or YAML parsing fails.
//
// Example:
//
//	type Config struct {
//	    Name string `yaml:"name"`
//	    Port int    `yaml:"port"`
//	}
//
//	var cfg Config
//	if err := goconf.ParseYaml(&cfg, "config.yaml"); err != nil {
//	    log.Fatal(err)
//	}
func ParseYaml(config interface{}, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file %s: %w", filePath, err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	return nil
}
