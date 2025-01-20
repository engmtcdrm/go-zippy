package zippy

import (
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
		{"Invalid Path", "/invalid/path\0001", 100, 100, true},
		{"Mismatched Bytes", "/valid/path", 100, 200, true},
		{"Valid Copy", "/valid/path", 100, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := validateCopy(tt.filePath, tt.written, tt.expected)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
