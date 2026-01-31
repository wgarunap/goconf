# GoConf

A lightweight, flexible Go library for loading and validating environment-based configuration with beautiful output formatting.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Overview

GoConf simplifies environment variable management in Go applications by providing:
- **Type-safe configuration** - Parse environment variables directly into Go structs
- **Validation** - Built-in validation using struct tags
- **Default values** - Support for default values when environment variables are not set
- **Sensitive data masking** - Automatically mask sensitive fields in output
- **Multiple output formats** - Table format for development, JSON format for production
- **Zero configuration** - Works out of the box with sensible defaults

Perfect for containerized applications, microservices, and cloud-native deployments where configuration is managed through environment variables following the [12-factor app methodology](https://12factor.net/config).

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Features](#features)
- [Usage](#usage)
  - [Basic Configuration](#basic-configuration)
  - [Struct Tags](#struct-tags)
  - [Validation](#validation)
  - [Output Formats](#output-formats)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/wgarunap/goconf
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/wgarunap/goconf"
)

type Config struct {
    DatabaseURL string `env:"DATABASE_URL" validate:"required,uri"`
    Port        int    `env:"PORT" envDefault:"8080" validate:"gte=1024,lte=65535"`
    APIKey      string `env:"API_KEY" validate:"required" secret:"true"`
}

var AppConfig Config

func (Config) Register() error {
    return goconf.ParseEnv(&AppConfig)
}

func (Config) Validate() error {
    return goconf.StructValidator(AppConfig)
}

func (Config) Print() interface{} {
    return AppConfig
}

func main() {
    if err := goconf.Load(new(Config)); err != nil {
        log.Fatal(err)
    }

    log.Println("Configuration loaded successfully!")
}
```

## Features

### 🔧 Environment Variable Parsing
- Automatically parse environment variables into typed struct fields
- Support for common Go types: `string`, `int`, `bool`, `float64`, etc.
- Custom type support through struct tag configuration

### ✅ Built-in Validation
- Powered by [go-playground/validator](https://github.com/go-playground/validator)
- Extensive validation rules: `required`, `uri`, `email`, `min`, `max`, `gte`, `lte`, etc.
- Custom validation rules support

### 🎨 Multiple Output Formats

#### Table Format (Default)
Beautiful Unicode table output for development environments:
```
┌─────────────┬───────────────────────────┐
│   CONFIG    │          VALUE            │
├─────────────┼───────────────────────────┤
│ DatabaseURL │ postgres://localhost:5432 │
│ Port        │ 8080                      │
│ APIKey      │ ***************           │
└─────────────┴───────────────────────────┘
```

#### JSON Format
Timestamped JSON output for production and centralized logging:
```
2026/01/31 10:30:15 {
  "DatabaseURL": "postgres://localhost:5432",
  "Port": 8080,
  "APIKey": "***************"
}
```

### 🔒 Sensitive Data Protection
Automatically mask sensitive fields marked with `secret:"true"` tag in all output formats.

### 📦 Default Values
Set fallback values using the `envDefault` tag when environment variables are not provided.

## Usage

### Basic Configuration

1. **Define your configuration struct** with environment variable mappings:

```go
type AppConfig struct {
    // Basic string field
    AppName string `env:"APP_NAME" envDefault:"MyApp"`

    // Integer with validation
    Port int `env:"PORT" envDefault:"8080" validate:"gte=1024,lte=65535"`

    // Boolean flag
    Debug bool `env:"DEBUG" envDefault:"false"`

    // Required field with validation
    DatabaseURL string `env:"DATABASE_URL" validate:"required,uri"`

    // Sensitive data (will be masked in output)
    APISecret string `env:"API_SECRET" validate:"required" secret:"true"`
}
```

2. **Implement required interfaces**:

```go
var Config AppConfig

func (AppConfig) Register() error {
    return goconf.ParseEnv(&Config)
}

func (AppConfig) Validate() error {
    return goconf.StructValidator(Config)
}

func (AppConfig) Print() interface{} {
    return Config
}
```

3. **Load configuration**:

```go
func main() {
    if err := goconf.Load(new(AppConfig)); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Use your configuration
    log.Printf("Starting %s on port %d", Config.AppName, Config.Port)
}
```

### Struct Tags

GoConf uses struct tags to configure field behavior:

| Tag | Description | Example |
|-----|-------------|---------|
| `env` | Environment variable name | `env:"PORT"` |
| `envDefault` | Default value if env var not set | `envDefault:"8080"` |
| `validate` | Validation rules (comma-separated) | `validate:"required,uri"` |
| `secret` | Mark field as sensitive (masks in output) | `secret:"true"` |

### Validation

GoConf uses [go-playground/validator](https://github.com/go-playground/validator) for validation. Common validation rules:

| Rule | Description | Example |
|------|-------------|---------|
| `required` | Field must be set | `validate:"required"` |
| `uri` | Must be a valid URI | `validate:"uri"` |
| `email` | Must be a valid email | `validate:"email"` |
| `url` | Must be a valid URL | `validate:"url"` |
| `min=N` | Minimum length (string) or value (number) | `validate:"min=3"` |
| `max=N` | Maximum length (string) or value (number) | `validate:"max=100"` |
| `gte=N` | Greater than or equal to | `validate:"gte=0"` |
| `lte=N` | Less than or equal to | `validate:"lte=100"` |
| `oneof=A B C` | Value must be one of the options | `validate:"oneof=dev staging prod"` |

**Multiple rules:**
```go
Port int `env:"PORT" validate:"required,gte=1024,lte=65535"`
```

### Output Formats

#### Table Format (Default)

Best for local development and debugging:

```go
func main() {
    // Table format is used by default
    if err := goconf.Load(new(Config)); err != nil {
        log.Fatal(err)
    }
}
```

#### JSON Format

Ideal for production environments, containerized deployments, and centralized logging systems:

```go
func main() {
    // Switch to JSON format
    goconf.SetOutputFormat(goconf.OutputFormatJSON)

    if err := goconf.Load(new(Config)); err != nil {
        log.Fatal(err)
    }
}
```

**Environment-based format selection:**
```go
func main() {
    // Use JSON in production, table in development
    if os.Getenv("ENV") == "production" {
        goconf.SetOutputFormat(goconf.OutputFormatJSON)
    }

    if err := goconf.Load(new(Config)); err != nil {
        log.Fatal(err)
    }
}
```


### Interfaces

#### `Configer`
Must be implemented by configuration structs.

```go
type Configer interface {
    Register() error
}
```

#### `Validater` (Optional)
Implement to enable validation.

```go
type Validater interface {
    Validate() error
}
```

#### `Printer` (Optional)
Implement to enable configuration output.

```go
type Printer interface {
    Print() interface{}
}
```

### Constants

```go
const (
    OutputFormatTable OutputFormat = "table" // Default: Unicode table
    OutputFormatJSON  OutputFormat = "json"  // JSON with timestamps
)
```


## Best Practices

### 1. Use Validation Rules
Always validate critical configuration fields:
```go
type Config struct {
    APIKey string `env:"API_KEY" validate:"required,min=32"`
    Port   int    `env:"PORT" validate:"required,gte=1024,lte=65535"`
}
```

### 2. Provide Sensible Defaults
Use `envDefault` for non-critical configuration:
```go
type Config struct {
    Port     int    `env:"PORT" envDefault:"8080"`
    LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
```

### 3. Mark Sensitive Data
Always mark secrets and passwords:
```go
type Config struct {
    APIKey   string `env:"API_KEY" secret:"true"`
    Password string `env:"DB_PASSWORD" secret:"true"`
}
```

### 4. Use Environment-Specific Formats
```go
if os.Getenv("ENVIRONMENT") == "production" {
    goconf.SetOutputFormat(goconf.OutputFormatJSON)
}
```

### 5. Group Related Configuration
Organize configuration into logical groups:
```go
type Config struct {
    Server   ServerConfig
    Database DBConfig
    Cache    CacheConfig
}
```

### 6. Document Environment Variables
Create a `.env.example` file:
```bash
# Application
APP_NAME=myapp
APP_VERSION=1.0.0

# Server
PORT=8080
HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=postgres
DB_PASSWORD=secret
```

### 7. Fail Fast on Configuration Errors
```go
func main() {
    if err := goconf.Load(new(Config)); err != nil {
        log.Fatalf("Configuration error: %v", err)
    }
}
```

## Contributing

We welcome contributions! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit your changes**: `git commit -m 'Add amazing feature'`
4. **Push to the branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Setup

```bash
# Clone the repository
git clone https://github.com/wgarunap/goconf.git
cd goconf

# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -v -cover ./...
```

### Code Quality

- Write clear, idiomatic Go code
- Add tests for new features
- Ensure all tests pass
- Update documentation
- Follow existing code style

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [caarlos0/env](https://github.com/caarlos0/env) - Environment variable parsing
- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation
- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) - ASCII table formatting

## Support

- 📫 Issues: [GitHub Issues](https://github.com/wgarunap/goconf/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/wgarunap/goconf/discussions)
- 📖 Documentation: [pkg.go.dev](https://pkg.go.dev/github.com/wgarunap/goconf)

---

Made with ❤️ for the Go community
