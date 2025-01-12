# Go Config
Library to load env configuration in Golang

### Features
- Load Env Configuration
- Print configuration
- Validate Env Configuration

### How to use it

```go
package main

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/wgarunap/goconf"
	"log"
)

type Conf struct {
	Name     string `env:"MY_NAME"`
	Username string `env:"MY_USERNAME" secret:"true"`
	Password string `env:"MY_PASSWORD" secret:"true"`
}

var Config Conf

func (Conf) Register() error {
	return env.Parse(&Config)
}

func (Conf) Validate() error {
	if Config.Name == "" {
		return errors.New(`MY_NAME environmental variable cannot be empty`)
	}
	if Config.Username == "" {
		return errors.New(`MY_USERNAME environmental variable cannot be empty`)
	}
	if Config.Password == "" {
		return errors.New(`MY_PASSWORD environmental variable cannot be empty`)
	}
	return nil
}

func (Conf) Print() interface{} {
	return Config
}

func main() {
	err := goconf.Load(
		new(Conf),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(`configuration successfully loaded`)
}
```