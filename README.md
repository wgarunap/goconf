# Go Config
Library to load env configuration

### How to use it

```go
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
	config.Load(
		new(Conf),
	)
    log.Info(`configuration loaded, `,Config.Name)
}

```