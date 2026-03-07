package zippy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateCopy(t *testing.T) {
	tests := []struct {
		testName string
		filePath string
		written  int64
		expected int64
		wantErr  bool
	}{
		{"Valid Copy", "/valid/path", 100, 100, false},
		{"Invalid Path", "/invalid/path\0001", 100, 100, true},
		{"Mismatched Bytes", "/valid/path", 100, 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			pathToCheck := tt.filePath

			if tt.filePath == "/valid/path" {
				pathToCheck = filepath.Join(t.TempDir(), "valid")
				if err := os.WriteFile(pathToCheck, []byte("Test File"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			err := validateCopy(pathToCheck, tt.written, tt.expected)
			if (err != nil) != tt.wantErr {
				t.Errorf("test '%s': validateCopy() error = %v, wantErr %v", tt.testName, err, tt.wantErr)
			}
		})
	}
}
