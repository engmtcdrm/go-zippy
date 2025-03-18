package testutils

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = CreateTempFile(tempDir, "testfile-*.txt")
	if err != nil {
		t.Fatalf("CreateTempFile() error = %v", err)
	}

	// Check if the file was created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
}

func TestCreateTestFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if _, err = CreateTestFiles(tempDir, 3, 2); err != nil {
		t.Fatalf("CreateTestFiles() error = %v", err)
	}

	// Check if the files and subdirectories were created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	expectedFiles := 3
	expectedDirs := 2
	actualFiles := 0
	actualDirs := 0

	for _, file := range files {
		if file.IsDir() {
			actualDirs++
			subFiles, err := os.ReadDir(filepath.Join(tempDir, file.Name()))
			if err != nil {
				t.Fatalf("Failed to read subdir: %v", err)
			}
			actualFiles += len(subFiles)
		} else {
			actualFiles++
		}
	}

	if actualFiles != expectedFiles+expectedFiles*expectedDirs {
		t.Errorf("Expected %d files, got %d", expectedFiles+expectedFiles*expectedDirs, actualFiles)
	}

	if actualDirs != expectedDirs {
		t.Errorf("Expected %d directories, got %d", expectedDirs, actualDirs)
	}
}

func TestCreateZipFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = CreateZipFile(zipFilePath, 3, 2)
	if err != nil {
		t.Fatalf("CreateZipFile() error = %v", err)
	}

	// Check if the zip file was created
	_, err = os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Expected zip file %s does not exist", zipFilePath)
	}

	// Check the contents of the zip file
	zFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		t.Fatalf("Failed to open zip file: %v", err)
	}
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

	if actualFiles != expectedFiles+expectedFiles*expectedDirs {
		t.Errorf("Expected %d files, got %d", expectedFiles+expectedFiles*expectedDirs, actualFiles)
	}

	if actualDirs != expectedDirs {
		t.Errorf("Expected %d directories, got %d", expectedDirs, actualDirs)
	}
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
