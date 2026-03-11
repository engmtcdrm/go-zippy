package testutils

import (
	"errors"
	"os"
	"testing"
)

// Tests for [PermissionTest] function.
func TestPermissionTest(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name      string
		funcToRun interface{}
		args      []interface{}
		wantErr   bool
	}{
		{
			name:      "0-arg function",
			funcToRun: func() error { return nil },
			args:      []interface{}{},
			wantErr:   false,
		},
		{
			name:      "1-arg function",
			funcToRun: func(arg1 string) error { return nil },
			args:      []interface{}{tempFile.Name()},
			wantErr:   false,
		},
		{
			name:      "2-arg function",
			funcToRun: func(arg1, arg2 string) error { return nil },
			args:      []interface{}{tempFile.Name(), "arg2"},
			wantErr:   false,
		},
		{
			name: "variadic function",
			funcToRun: func(args ...string) error {
				if len(args) != 2 {
					return errors.New("expected 2 arguments")
				}
				return nil
			},
			args:    []interface{}{tempFile.Name(), "arg2"},
			wantErr: false,
		},
		{
			name: "function with error",
			funcToRun: func(args ...string) error {
				return errors.New("test error")
			},
			args:    []interface{}{tempFile.Name(), "arg2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PermissionTest(tempFile.Name(), tt.funcToRun, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PermissionTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
