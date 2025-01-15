package goconf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type Conf struct {
	Name string `env:"MY_NAME"`
	Age  int    `env:"MY_AGE"`
	Team string `env:"MY_TEAM"`
}

func TestParseEnv_Success(t *testing.T) {
	updateEnv()
	tests := []struct {
		name           string
		input          Conf
		expectedOutput Conf
		errorExpected  bool
	}{
		{
			name:  "Successfully parsed the configs",
			input: Conf{},
			expectedOutput: Conf{
				Name: "coderx",
				Age:  99,
				Team: "backend",
			},
			errorExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ParseEnv(&test.input)
			if !test.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOutput, test.input)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestParseEnv_Failure(t *testing.T) {
	updateEnv()
	tests := []struct {
		name           string
		input          *Conf
		expectedOutput *Conf
		errorExpected  bool
	}{
		{
			name:           "config parse failure scenario",
			input:          &Conf{},
			expectedOutput: &Conf{},
			errorExpected:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ParseEnv(&test.input)
			if !test.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOutput, test.input)
			} else {
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
