package testutils

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTempFile(t *testing.T) {
	tempDir := t.TempDir()

	_, err := CreateTempFile(tempDir, "testfile-*.txt")
	assert.NoError(t, err)

	// Check if the file was created
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)

	assert.Len(t, files, 1)
}

func TestCreateTestFiles(t *testing.T) {
	tempDir := t.TempDir()

	_, err := CreateTestFiles(tempDir, 3, 2)
	assert.NoError(t, err)

	// Check if the files and subdirectories were created
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)

	expectedFiles := 3
	expectedDirs := 2
	actualFiles := 0
	actualDirs := 0

	for _, file := range files {
		if file.IsDir() {
			actualDirs++
			subFiles, err := os.ReadDir(filepath.Join(tempDir, file.Name()))
			assert.NoError(t, err)
			actualFiles += len(subFiles)
		} else {
			actualFiles++
		}
	}

	assert.Equal(t, expectedFiles+expectedFiles*expectedDirs, actualFiles)
	assert.Equal(t, expectedDirs, actualDirs)
}

func TestCreateZipFile(t *testing.T) {
	tempDir := t.TempDir()
	zipFilePath := filepath.Join(tempDir, "test.zip")

	err := CreateZipFile(zipFilePath, 3, 2)
	assert.NoError(t, err)

	// Check if the zip file was created
	_, err = os.Stat(zipFilePath)
	assert.NoError(t, err)

	// Check the contents of the zip file
	zFile, err := zip.OpenReader(zipFilePath)
	assert.NoError(t, err)
	defer zFile.Close()

	expectedFiles := 3
	expectedDirs := 2
	actualFiles := 0
	actualDirs := 0

	for _, file := range zFile.File {
		if file.FileInfo().IsDir() {
			actualDirs++
		} else {
			actualFiles++
		}
	}

	assert.Equal(t, expectedFiles+expectedFiles*expectedDirs, actualFiles)
	assert.Equal(t, expectedDirs, actualDirs)
}

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
