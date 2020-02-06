package main

import (
	"github.com/caarlos0/env"
	"github.com/pickme-go/log/v2"
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
		log.Fatal("error loading stream config, ", err)
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

	config.Load(
		new(Conf),
	)

	log.Info(`config successfully loaded`)
}
