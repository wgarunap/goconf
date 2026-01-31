package goconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructValidator(t *testing.T) {
	type config struct {
		Name string `validate:"required"`
		Age  int    `validate:"gte=0,lte=130"`
	}
	tests := map[string]struct {
		config      config
		expectedErr string
	}{
		"required field available": {
			config: config{
				Name: "XX",
			},
			expectedErr: "",
		},
		"required field not available": {
			config:      config{},
			expectedErr: "Error:Field validation",
		},
		"validation condition successful": {
			config: config{
				Name: "some name",
				Age:  10,
			},
			expectedErr: "",
		},
		"validation condition failed": {
			config: config{
				Name: "some name",
				Age:  -5,
			},
			expectedErr: "Error:Field validation",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := StructValidator(test.config)
			if test.expectedErr != "" {
				require.ErrorContains(t, err, test.expectedErr)
			}
		})
	}
}
