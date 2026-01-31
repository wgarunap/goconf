// Package main demonstrates the usage of goconf library for loading environment-based configuration
package main

import (
	"log"
	"os"

	"github.com/wgarunap/goconf"
)

// Conf holds the application configuration loaded from environment variables
type Conf struct {
	Name        string `env:"MY_NAME" validate:"required"`
	ExampleHost string `env:"EXAMPLE_HOST" validate:"required,uri"`
	Port        int    `env:"EXAMPLE_PORT" validate:"gte=8080,lte=9000"`
	Password    string `env:"MY_PASSWORD" secret:"true"`
}

// Config is the global configuration instance
var Config Conf

// Register loads environment variables into the Config struct
func (Conf) Register() error {
	return goconf.ParseEnv(&Config)
}

// Validate ensures the loaded configuration meets validation requirements
func (Conf) Validate() error {
	return goconf.StructValidator(Config)
}

// Print returns the configuration for display purposes
func (Conf) Print() interface{} {
	return Config
}

func main() {
	_ = os.Setenv("MY_NAME", "GoConf")
	_ = os.Setenv("EXAMPLE_HOST", "https://github.com/wgarunap/goconf")
	_ = os.Setenv("EXAMPLE_PORT", "8090")
	_ = os.Setenv("MY_PASSWORD", "testUserPassword")

	// Uncomment the line below to use JSON output format
	// Useful for centralized logging systems in production
	goconf.SetOutputFormat(goconf.OutputFormatJSON)

	if err := goconf.Load(new(Conf)); err != nil {
		log.Fatal(err)
	}

	log.Println(`configuration successfully loaded`)

	log.Printf("%+v", Config)
}
