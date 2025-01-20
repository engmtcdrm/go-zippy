package zippy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestUnzipFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test zip file
	zipFile, err := os.Create(filepath.Join(tempDir, "test.zip"))
	if err != nil {
		t.Fatalf("Failed to create test zip file: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)
	fileWriter, err := zipWriter.Create("testfile.txt")
	if err != nil {
		t.Fatalf("Failed to add file to test zip: %v", err)
	}
	_, err = fileWriter.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}
	zipWriter.Close()
	zipFile.Close()

	// Test unzipFile function
	zipReader, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		t.Fatalf("Failed to open test zip file: %v", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		err := unzipFile(file, filepath.Join(tempDir, "test_output"))
		if err != nil {
			t.Errorf("unzipFile() error = %v", err)
		}
	}
}

func TestUnzip(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		testName string
		zipFile  string
		dest     string
		wantErr  bool
	}{
		{"Zip Does Not Exist", "nonexistent.zip", filepath.Join(tempDir, "output"), true},
		{"Zip Exists", filepath.Join(tempDir, "test.zip"), filepath.Join(tempDir, "output"), false},
	}

	// Create a test zip file
	zipFile, err := os.Create(filepath.Join(tempDir, "test.zip"))
	if err != nil {
		t.Fatalf("Failed to create test zip file: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)
	fileWriter, err := zipWriter.Create("testfile.txt")
	if err != nil {
		t.Fatalf("Failed to add file to test zip: %v", err)
	}
	_, err = fileWriter.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}
	zipWriter.Close()
	zipFile.Close()

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := Unzip(tt.zipFile, tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unzip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
