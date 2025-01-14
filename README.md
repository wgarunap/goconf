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
	Name        string `env:"MY_NAME" validate:"required"`
	ExampleHost string `env:"EXAMPLE_HOST" validate:"required,uri"`
	Port        int    `env:"EXAMPLE_PORT" validate:"gte=8080,lte=9000"`
	Password    string `env:"MY_PASSWORD" secret:"true"`
}

var Config Conf

func (Conf) Register() error {
	return env.Parse(&Config)
}

func (Conf) Validate() error {
	return goconf.StructValidator(Config)
}

func (Conf) Print() interface{} {
	return Config
}

func main() {
    if err := goconf.Load(new(Conf));err!=nil{
        log.Fatal(err)
    }
    
    log.Println(`configuration successfully loaded`)
}
```