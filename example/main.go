package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/tryfix/log"
	"github.com/wgarunap/config"
	"os"
)

type Conf struct {
	Name string `env:"MY_NAME"`
}

var Config Conf

func (Conf) Register() {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal("error loading stream goconf, ", err)
	}
}

func (Conf) Validate() {
	if Config.Name == "" {
		log.Fatal(`MY_NAME environmental variable cannot be empty`)
	}
}

func (Conf) Print() interface{} {
	return Config
}

func main() {
	_ = os.Setenv("MY_NAME", "My First Configuration")

	goconf.Load(
		new(Conf),
	)

	if Config.Name != `My First Configuration` {
		log.Fatal(`error while comparing config`)
	}

	log.Info(`goconf successfully loaded`)
}
