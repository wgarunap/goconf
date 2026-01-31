package goconf

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/wgarunap/goconf/mocks"
)

func TestLoad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		mockConfig     Configer
		expectedErr    error
		expectedOutput string
	}{
		{
			name:           "successful config registration and printing",
			mockConfig:     mock(ctrl, nil, nil),
			expectedErr:    nil,
			expectedOutput: "┌──────────────┬─────────────────┐\n│    CONFIG    │      VALUE      │\n├──────────────┼─────────────────┤\n│ DatabaseName │ test_db         │\n│ Username     │ *************** │\n│ Password     │ *************** │\n└──────────────┴─────────────────┘\n",
		},
		{
			name:           "config validation failure scenario",
			mockConfig:     mock(ctrl, nil, errors.New("validation failed")),
			expectedErr:    errors.New("validation failed"),
			expectedOutput: "",
		},
		{
			name:           "config registration failure scenario",
			mockConfig:     mock(ctrl, errors.New("registration failed"), nil),
			expectedErr:    errors.New("registration failed"),
			expectedOutput: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			oldStdOut := os.Stdout
			os.Stdout = w

			err := Load(test.mockConfig)
			if err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErr, err)
			}

			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)
			os.Stdout = oldStdOut

			if err == nil {
				output := buf.String()
				assert.Contains(t, output, test.expectedOutput)
			}
		})
	}
}

func TestLoadWithJSONOutput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Save original format and restore after test
	originalFormat := currentOutputFormat
	defer func() { currentOutputFormat = originalFormat }()

	tests := []struct {
		name           string
		outputFormat   OutputFormat
		mockConfig     Configer
		expectedErr    error
		expectedOutput string
	}{
		{
			name:           "successful config registration with JSON output",
			outputFormat:   OutputFormatJSON,
			mockConfig:     mock(ctrl, nil, nil),
			expectedErr:    nil,
			expectedOutput: `"DatabaseName": "test_db"`,
		},
		{
			name:           "JSON output masks sensitive data",
			outputFormat:   OutputFormatJSON,
			mockConfig:     mock(ctrl, nil, nil),
			expectedErr:    nil,
			expectedOutput: `"Username": "***************"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetOutputFormat(test.outputFormat)

			r, w, _ := os.Pipe()
			oldStdOut := os.Stdout
			os.Stdout = w

			err := Load(test.mockConfig)
			if err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErr, err)
			}

			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)
			os.Stdout = oldStdOut

			if err == nil {
				output := buf.String()
				assert.Contains(t, output, test.expectedOutput)
			}
		})
	}
}

func TestNestedStructSecretMasking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type NestedConfig struct {
		AppName  string `yaml:"app_name"`
		Port     int    `yaml:"port"`
		Database struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password" secret:"true"`
		} `yaml:"database"`
		Redis struct {
			Host     string `yaml:"host"`
			Password string `yaml:"password" secret:"true"`
		} `yaml:"redis"`
	}

	tests := []struct {
		name          string
		outputFormat  OutputFormat
		config        NestedConfig
		expectedTable []string // Expected strings in table output
		expectedJSON  []string // Expected strings in JSON output
	}{
		{
			name:         "table format masks nested secrets",
			outputFormat: OutputFormatTable,
			config: NestedConfig{
				AppName: "TestApp",
				Port:    8080,
				Database: struct {
					Host     string `yaml:"host"`
					Port     int    `yaml:"port"`
					Username string `yaml:"username"`
					Password string `yaml:"password" secret:"true"`
				}{
					Host:     "localhost",
					Port:     5432,
					Username: "dbuser",
					Password: "secret123",
				},
				Redis: struct {
					Host     string `yaml:"host"`
					Password string `yaml:"password" secret:"true"`
				}{
					Host:     "redis-host",
					Password: "redis-secret",
				},
			},
			expectedTable: []string{
				"Database.Password",
				SensitiveDataMaskString,
				"Redis.Password",
				"Database.Username │ dbuser",
			},
		},
		{
			name:         "json format masks nested secrets",
			outputFormat: OutputFormatJSON,
			config: NestedConfig{
				AppName: "TestApp",
				Port:    8080,
				Database: struct {
					Host     string `yaml:"host"`
					Port     int    `yaml:"port"`
					Username string `yaml:"username"`
					Password string `yaml:"password" secret:"true"`
				}{
					Host:     "localhost",
					Port:     5432,
					Username: "dbuser",
					Password: "secret123",
				},
				Redis: struct {
					Host     string `yaml:"host"`
					Password string `yaml:"password" secret:"true"`
				}{
					Host:     "redis-host",
					Password: "redis-secret",
				},
			},
			expectedJSON: []string{
				`"Password": "***************"`,
				`"Username": "dbuser"`,
				`"Host": "localhost"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original format and restore after test
			originalFormat := currentOutputFormat
			defer func() { currentOutputFormat = originalFormat }()

			SetOutputFormat(tt.outputFormat)

			mockConfiger := mocks.NewMockConfiger(ctrl)
			mockValidater := mocks.NewMockValidater(ctrl)
			mockPrinter := mocks.NewMockPrinter(ctrl)

			mockConfiger.EXPECT().Register().Return(nil).AnyTimes()
			mockValidater.EXPECT().Validate().Return(nil).AnyTimes()
			mockPrinter.EXPECT().Print().Return(tt.config).AnyTimes()

			configInterface := &struct {
				*mocks.MockConfiger
				*mocks.MockValidater
				*mocks.MockPrinter
			}{
				MockConfiger:  mockConfiger,
				MockValidater: mockValidater,
				MockPrinter:   mockPrinter,
			}

			r, w, _ := os.Pipe()
			oldStdOut := os.Stdout
			os.Stdout = w

			err := Load(configInterface)
			assert.NoError(t, err)

			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)
			os.Stdout = oldStdOut

			output := buf.String()

			if tt.outputFormat == OutputFormatTable {
				for _, expected := range tt.expectedTable {
					assert.Contains(t, output, expected, "Table output should contain: %s", expected)
				}
				// Ensure actual passwords are NOT in output
				assert.NotContains(t, output, "secret123")
				assert.NotContains(t, output, "redis-secret")
			} else {
				for _, expected := range tt.expectedJSON {
					assert.Contains(t, output, expected, "JSON output should contain: %s", expected)
				}
				// Ensure actual passwords are NOT in output
				assert.NotContains(t, output, "secret123")
				assert.NotContains(t, output, "redis-secret")
			}
		})
	}
}

func TestExtractFields(t *testing.T) {
	type NestedStruct struct {
		PublicField  string
		SecretField  string `secret:"true"`
		NestedLevel2 struct {
			Field1 string
			Field2 int
			Secret string `secret:"true"`
		}
	}

	testData := NestedStruct{
		PublicField: "public",
		SecretField: "secret",
		NestedLevel2: struct {
			Field1 string
			Field2 int
			Secret string `secret:"true"`
		}{
			Field1: "value1",
			Field2: 42,
			Secret: "nested-secret",
		},
	}

	values := reflect.ValueOf(testData)
	result := extractFields("", values)

	// Verify field names with dot notation
	fieldNames := make(map[string]bool)
	for _, row := range result {
		fieldNames[row[0]] = true
	}

	assert.True(t, fieldNames["PublicField"], "Should have PublicField")
	assert.True(t, fieldNames["SecretField"], "Should have SecretField")
	assert.True(t, fieldNames["NestedLevel2.Field1"], "Should have NestedLevel2.Field1")
	assert.True(t, fieldNames["NestedLevel2.Field2"], "Should have NestedLevel2.Field2")
	assert.True(t, fieldNames["NestedLevel2.Secret"], "Should have NestedLevel2.Secret")

	// Verify secret masking
	for _, row := range result {
		if row[0] == "SecretField" || row[0] == "NestedLevel2.Secret" {
			assert.Equal(t, SensitiveDataMaskString, row[1], "Secret field should be masked")
		}
	}
}

func TestExtractJSONFields(t *testing.T) {
	type Config struct {
		Name     string
		Password string `secret:"true"`
		Database struct {
			Host     string
			Password string `secret:"true"`
		}
	}

	testData := Config{
		Name:     "TestApp",
		Password: "secret123",
		Database: struct {
			Host     string
			Password string `secret:"true"`
		}{
			Host:     "localhost",
			Password: "db-secret",
		},
	}

	values := reflect.ValueOf(testData)
	result := extractJSONFields(values)

	assert.Equal(t, "TestApp", result["Name"])
	assert.Equal(t, SensitiveDataMaskString, result["Password"])

	dbMap, ok := result["Database"].(map[string]interface{})
	assert.True(t, ok, "Database should be a map")
	assert.Equal(t, "localhost", dbMap["Host"])
	assert.Equal(t, SensitiveDataMaskString, dbMap["Password"])
}

func mock(ctrl *gomock.Controller, registerErr, validateErr error) Configer {
	mockConfiger := mocks.NewMockConfiger(ctrl)
	mockValidater := mocks.NewMockValidater(ctrl)
	mockPrinter := mocks.NewMockPrinter(ctrl)

	mockConfiger.EXPECT().Register().Return(registerErr).AnyTimes()
	mockValidater.EXPECT().Validate().Return(validateErr).AnyTimes()
	mockPrinter.EXPECT().Print().Return(struct {
		DatabaseName string `secret:"false"`
		Username     string `secret:"true"`
		Password     string `secret:"true"`
	}{
		DatabaseName: "test_db",
		Username:     "test_user",
		Password:     "test_password",
	}).AnyTimes()

	return &struct {
		*mocks.MockConfiger
		*mocks.MockValidater
		*mocks.MockPrinter
	}{
		MockConfiger:  mockConfiger,
		MockValidater: mockValidater,
		MockPrinter:   mockPrinter,
	}
}
