package zippy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestContents(t *testing.T) {
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
	_, err = zipWriter.Create("testfile.txt")
	if err != nil {
		t.Fatalf("Failed to add file to test zip: %v", err)
	}
	zipWriter.Close()
	zipFile.Close()

	tests := []struct {
		testName        string
		zipFile         string
		expectedFileCnt int
		wantErr         bool
	}{
		{"Zip Exists", filepath.Join(tempDir, "test.zip"), 1, false},
		{"Zip Does Not Exist", "nonexistent.zip", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			zipFiles, err := Contents(tt.zipFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contents() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(zipFiles) != tt.expectedFileCnt {
				t.Errorf("Contents() got %d files, want %d", len(zipFiles), tt.expectedFileCnt)
			}
		})
	}
}
