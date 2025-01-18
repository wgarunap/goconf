package main

import (
	"log"
	"os"

	"github.com/wgarunap/goconf"
)

type Conf struct {
	Name        string `env:"MY_NAME" validate:"required"`
	ExampleHost string `env:"EXAMPLE_HOST" validate:"required,uri"`
	Port        int    `env:"EXAMPLE_PORT" validate:"gte=8080,lte=9000"`
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
	_ = os.Setenv("MY_NAME", "GoConf")
	_ = os.Setenv("EXAMPLE_HOST", "https://github.com/wgarunap/goconf")
	_ = os.Setenv("EXAMPLE_PORT", "8090")
	_ = os.Setenv("MY_PASSWORD", "testUserPassword")

	if err := goconf.Load(new(Conf)); err != nil {
		log.Fatal(err)
	}

	log.Println(`configuration successfully loaded`)
	log.Printf("%+v", Config)
}
