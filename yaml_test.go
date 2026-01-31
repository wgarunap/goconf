package goconf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseYaml(t *testing.T) {
	type Config struct {
		Name     string `yaml:"name" validate:"required"`
		Port     int    `yaml:"port" validate:"required,gte=1024,lte=65535"`
		Debug    bool   `yaml:"debug"`
		Database struct {
			Host     string `yaml:"host" validate:"required"`
			Port     int    `yaml:"port" validate:"required,gte=1024,lte=65535"`
			Username string `yaml:"username" validate:"required"`
			Password string `yaml:"password" validate:"required"`
		} `yaml:"database"`
	}

	tests := []struct {
		name        string
		yamlContent string
		expected    Config
		expectedErr bool
	}{
		{
			name: "valid YAML configuration",
			yamlContent: `name: TestApp
port: 8080
debug: true
database:
  host: localhost
  port: 5432
  username: testuser
  password: testpass
`,
			expected: Config{
				Name:  "TestApp",
				Port:  8080,
				Debug: true,
				Database: struct {
					Host     string `yaml:"host" validate:"required"`
					Port     int    `yaml:"port" validate:"required,gte=1024,lte=65535"`
					Username string `yaml:"username" validate:"required"`
					Password string `yaml:"password" validate:"required"`
				}{
					Host:     "localhost",
					Port:     5432,
					Username: "testuser",
					Password: "testpass",
				},
			},
			expectedErr: false,
		},
		{
			name: "partial YAML configuration with defaults",
			yamlContent: `name: PartialApp
port: 3000
`,
			expected: Config{
				Name: "PartialApp",
				Port: 3000,
			},
			expectedErr: false,
		},
		{
			name:        "empty YAML file",
			yamlContent: ``,
			expected:    Config{},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary YAML file
			tmpDir := t.TempDir()
			yamlFile := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(yamlFile, []byte(tt.yamlContent), 0644)
			require.NoError(t, err)

			// Parse YAML
			var cfg Config
			err = ParseYaml(&cfg, yamlFile)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Name, cfg.Name)
				assert.Equal(t, tt.expected.Port, cfg.Port)
				assert.Equal(t, tt.expected.Debug, cfg.Debug)
				assert.Equal(t, tt.expected.Database.Host, cfg.Database.Host)
				assert.Equal(t, tt.expected.Database.Port, cfg.Database.Port)
			}
		})
	}
}

func TestParseYaml_FileNotFound(t *testing.T) {
	type Config struct {
		Name string `yaml:"name"`
	}

	var cfg Config
	err := ParseYaml(&cfg, "/nonexistent/path/config.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read YAML file")
}

func TestParseYaml_InvalidYAML(t *testing.T) {
	type Config struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	}

	// Create temporary invalid YAML file
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "invalid.yaml")
	invalidYAML := `name: TestApp
port: [invalid
  structure
`
	err := os.WriteFile(yamlFile, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	var cfg Config
	err = ParseYaml(&cfg, yamlFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal YAML data")
}
