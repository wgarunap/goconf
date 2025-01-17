package goconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Conf struct {
	Name string `env:"MY_NAME"`
	Age  int    `env:"MY_AGE"`
	Team string `env:"MY_TEAM"`
}

func TestParseEnv(t *testing.T) {
	updateEnv()
	tests := []struct {
		name           string
		expectedOutput Conf
		errorExpected  bool
	}{
		{
			name: "Successfully parsed the configs",
			expectedOutput: Conf{
				Name: "coderx",
				Age:  99,
				Team: "backend",
			},
			errorExpected: false,
		},
		{
			name:           "config parse failure scenario",
			expectedOutput: Conf{},
			errorExpected:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.errorExpected {
				conf := newConf()
				err := ParseEnv(&conf)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOutput, conf)
			} else {
				confPointer := newPointerToConf()
				err := ParseEnv(&confPointer)
				assert.Error(t, err)
			}
		})
	}
}

func updateEnv() {
	_ = os.Setenv("MY_NAME", "coderx")
	_ = os.Setenv("MY_AGE", "99")
	_ = os.Setenv("MY_TEAM", "backend")
}

func newConf() Conf {
	return Conf{}
}

func newPointerToConf() *Conf {
	return &Conf{}
}
