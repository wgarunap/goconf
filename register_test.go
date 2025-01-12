package goconf

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wgarunap/goconf/mocks"
	"os"
	"testing"
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
			expectedOutput: "+--------------+-----------------+\n|    CONFIG    |      VALUE      |\n+--------------+-----------------+\n| DatabaseName | test_db         |\n| Username     | *************** |\n| Password     | *************** |\n+--------------+-----------------+\n",
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
