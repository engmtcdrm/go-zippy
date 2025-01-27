package testutils

import (
	"archive/zip"
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
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "testfile.txt")
	_, err = CreateTempFile(tempDir, "testfile-*.txt")
	if err != nil {
		t.Fatalf("CreateTempFile() error = %v", err)
	}

	// Define a callback function that requires specific permissions
	callback := func(arg1 string, arg2 string) error {
		_, err := os.OpenFile(arg1, os.O_WRONLY, 0666)
		return err
	}

	// Test with granted permissions
	err = PermissionTest(tempDir, callback, tempDir, testFilePath)
	if err == nil {
		t.Errorf("PermissionTest() did not produce an error as expected = %v", err)
	}
}
