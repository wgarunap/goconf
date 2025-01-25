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
	"log"

	"github.com/wgarunap/goconf"	
)

type Conf struct {
	Name        string `env:"MY_NAME" validate:"required"`
	ExampleHost string `env:"EXAMPLE_HOST" envDefault:"localhost" validate:"required,uri"`
	Port        int    `env:"EXAMPLE_PORT" envDefault:"8081" validate:"gte=8080,lte=9000"`
	Password    string `env:"MY_PASSWORD" secret:"true"`
}

var Config Conf

func (Conf) Register() error {
	return goconf.ParseEnv(&Config)
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
