package main

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/tryfix/log"
	"github.com/wgarunap/goconf"
	"os"
)

type Conf struct {
	Name     string `env:"MY_NAME"`
	Username string `env:"MY_USERNAME" secret:""`
	Password string `env:"MY_PASSWORD" secret:""`
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
	_ = os.Setenv("MY_NAME", "My First Configuration")
	_ = os.Setenv("MY_USERNAME", "testUserName")
	_ = os.Setenv("MY_PASSWORD", "testUserPassword")

	err := goconf.Load(
		new(Conf),
	)
	if err != nil {
		log.Fatal(err)
	}
	if Config.Name != `My First Configuration` {
		log.Fatal(`error while comparing config`)
	}

	log.Info(`goconf successfully loaded`)
}
