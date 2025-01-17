package goconf

import "github.com/caarlos0/env/v11"

// ParseEnv parse the env values to given struct fields
// env variables are defined using struct tags. It utilizes the "github.com/caarlos0/env/v11"
// package for parse.
//
// Parameters:
//   - config (interface{}): The struct to be populated by env variables. The struct should have
//     env variables defined using tags such as `env:"USERNAME"`.
//
// Returns:
//   - error: Returns nil if the struct created using env variables. if fails, it returns an error
//     which provides details information about the error
//
// Usage Example:
//
//	type Config struct {
//	    Name string `env:"MY_NAME"`
//	    Age  int    `env:"MY_AGE"`
//      Team string `env:"MY_TEAM" envDefault:"backend"`
//	}
//
//	var conf Config
//	if err := Register(&conf); err != nil {
//	    // Handle env pass error, e.g., log or return
//	}
//
// Note:
//   - The function will panic if the `config` parameter is not a pointer to struct
//
// More env package information https://github.com/caarlos0/env/v11
func ParseEnv(config interface{}) error {
	return env.Parse(config)
}
