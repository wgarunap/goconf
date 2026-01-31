// Package main demonstrates YAML configuration loading with goconf
package main

import (
	"log"

	"github.com/wgarunap/goconf"
)

// Config holds configuration loaded from a YAML file
type Config struct {
	AppName  string `yaml:"app_name" validate:"required"`
	Port     int    `yaml:"port" validate:"required,gte=1024,lte=65535"`
	Debug    bool   `yaml:"debug"`
	Database struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required,gte=1024,lte=65535"`
		Database string `yaml:"database" validate:"required"`
		Username string `yaml:"username" validate:"required"`
		Password string `yaml:"password" validate:"required" secret:"true"`
	} `yaml:"database" validate:"required"`
	Redis struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required,gte=1024,lte=65535"`
		Password string `yaml:"password" secret:"true"`
	} `yaml:"redis" validate:"required"`
}

// AppConfig is the global configuration instance
var AppConfig Config

// Register loads the YAML configuration file
func (Config) Register() error {
	return goconf.ParseYaml(&AppConfig, "config.yaml")
}

// Validate validates the loaded configuration
func (Config) Validate() error {
	return goconf.StructValidator(AppConfig)
}

// Print returns the configuration for display
func (Config) Print() interface{} {
	return AppConfig
}

func main() {
	// Load YAML configuration
	// Uncomment to use JSON format: goconf.SetOutputFormat(goconf.OutputFormatJSON)
	if err := goconf.Load(new(Config)); err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded configuration for: %s", AppConfig.AppName)
	log.Printf("Server will run on port: %d", AppConfig.Port)
}
