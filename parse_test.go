package goconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		input          Conf
		expectedOutput Conf
		expectedErr    string
	}{
		{
			name:  "Successfully parsed the configs",
			input: Conf{},
			expectedOutput: Conf{
				Name: "coderx",
				Age:  99,
				Team: "backend",
			},
			expectedErr: "",
		},
		{
			name:           "config parse failure scenario",
			input:          Conf{},
			expectedOutput: Conf{},
			expectedErr:    "env: expected a pointer to a Struct",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedErr == "" {
				err := ParseEnv(&test.input)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOutput, test.input)
			} else {
				err := ParseEnv(test.input)
				require.ErrorContains(t, err, test.expectedErr)
			}
		})
	}
}

func updateEnv() {
	_ = os.Setenv("MY_NAME", "coderx")
	_ = os.Setenv("MY_AGE", "99")
	_ = os.Setenv("MY_TEAM", "backend")
}
